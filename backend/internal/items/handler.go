package items

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"techverin-backend/internal/httpx"
	"techverin-backend/internal/types"
)

type Handler struct {
	DB     *pgxpool.Pool
	Origin string
}

// ---------- helpers ----------

func idFromPath(prefix string, path string) (int64, error) {
	// expects /api/items/{id}
	s := strings.TrimPrefix(path, prefix)
	id, err := strconv.ParseInt(strings.Trim(s, "/"), 10, 64)
	if err != nil || id <= 0 {
		return 0, err
	}
	return id, nil
}

// ---------- health ----------

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, "ok", map[string]string{"status": "ok"})
}

// ---------- batch: preview & insert ----------

// POST /api/items/preview  (no save, just batch stats)
func (h *Handler) Preview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req types.BulkItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if len(req.Items) == 0 {
		httpx.Error(w, http.StatusBadRequest, "items cannot be empty")
		return
	}

	rows := make([]Row, 0, len(req.Items))
	for _, it := range req.Items {
		name := strings.TrimSpace(it.Name)
		if name == "" {
			httpx.Error(w, http.StatusBadRequest, "item name cannot be empty")
			return
		}
		if it.Price < 0 || isNaN(it.Price) || isInf(it.Price) {
			httpx.Error(w, http.StatusBadRequest, "price must be valid and >= 0")
			return
		}
		if it.Quantity < 1 {
			httpx.Error(w, http.StatusBadRequest, "quantity must be >= 1")
			return
		}
		rows = append(rows, Row{Name: name, Price: round2(it.Price), Quantity: it.Quantity})
	}

	line, qty, total, avgUnit, avgLine := ComputeBatchStats(rows)
	httpx.OK(w, "preview", &types.StatsResponse{
		LineItemCount: line,
		TotalQuantity: qty,
		TotalCost:     total,
		AvgUnitPrice:  avgUnit,
		AvgLineCost:   avgLine,
	})
}

// POST /api/items/bulk  (save with upsert, return batch + global)
func (h *Handler) BulkInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req types.BulkItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if len(req.Items) == 0 {
		httpx.Error(w, http.StatusBadRequest, "items cannot be empty")
		return
	}

	rows := make([]Row, 0, len(req.Items))
	for _, it := range req.Items {
		name := strings.TrimSpace(it.Name)
		if name == "" {
			httpx.Error(w, http.StatusBadRequest, "item name cannot be empty")
			return
		}
		if it.Price < 0 || isNaN(it.Price) || isInf(it.Price) {
			httpx.Error(w, http.StatusBadRequest, "price must be valid and >= 0")
			return
		}
		if it.Quantity < 1 {
			httpx.Error(w, http.StatusBadRequest, "quantity must be >= 1")
			return
		}
		rows = append(rows, Row{Name: name, Price: round2(it.Price), Quantity: it.Quantity})
	}

	// Save in a tx using upsert-on-conflict
	tx, err := h.DB.Begin(r.Context())
	if err != nil {
		log.Printf("BEGIN error: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "db transaction error")
		return
	}
	defer tx.Rollback(r.Context())

	repo := Repo{DB: h.DB}
	if err := repo.InsertBatch(r.Context(), tx, rows); err != nil {
		log.Printf("BATCH error: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "db insert error")
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		log.Printf("COMMIT error: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "db commit error")
		return
	}

	// Compute batch + global
	bLine, bQty, bTotal, bAvgUnit, bAvgLine := ComputeBatchStats(rows)
	gLine, gQty, gTotal, gAvgUnit, gAvgLine, err := repo.GlobalStats(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to compute global stats")
		return
	}

	httpx.OK(w, "saved", map[string]any{
		"batch": &types.StatsResponse{
			LineItemCount: bLine, TotalQuantity: bQty, TotalCost: bTotal,
			AvgUnitPrice: bAvgUnit, AvgLineCost: bAvgLine,
		},
		"global": &types.StatsResponse{
			LineItemCount: gLine, TotalQuantity: gQty, TotalCost: gTotal,
			AvgUnitPrice: gAvgUnit, AvgLineCost: gAvgLine,
		},
	})
}

// ---------- stats ----------

// GET /api/stats  (global stats)
func (h *Handler) GlobalStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	repo := Repo{DB: h.DB}
	line, qty, total, avgUnit, avgLine, err := repo.GlobalStats(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get stats")
		return
	}
	httpx.OK(w, "ok", &types.StatsResponse{
		LineItemCount: line,
		TotalQuantity: qty,
		TotalCost:     total,
		AvgUnitPrice:  avgUnit,
		AvgLineCost:   avgLine,
	})
}

// ---------- items CRUD ----------

// GET /api/items
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	limit := 200
	if q := r.URL.Query().Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 && n <= 1000 {
			limit = n
		}
	}
	repo := Repo{DB: h.DB}
	items, err := repo.List(r.Context(), limit)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to list")
		return
	}
	httpx.OK(w, "ok", items)
}

type updateBody struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// PUT /api/items/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id, err := idFromPath("/api/items/", r.URL.Path)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var b updateBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	name := strings.TrimSpace(b.Name)
	if name == "" {
		httpx.Error(w, http.StatusBadRequest, "name required")
		return
	}
	if b.Price < 0 || isNaN(b.Price) || isInf(b.Price) {
		httpx.Error(w, http.StatusBadRequest, "price must be valid and >= 0")
		return
	}
	if b.Quantity < 1 {
		httpx.Error(w, http.StatusBadRequest, "quantity must be >= 1")
		return
	}

	repo := Repo{DB: h.DB}
	row, merged, err := repo.UpdateWithMerge(r.Context(), id, name, round2(b.Price), b.Quantity)
	if err != nil {
		// repo returns pgx.ErrNoRows for not found
		httpx.Error(w, http.StatusNotFound, "not found")
		return
	}
	msg := "updated"
	if merged {
		msg = "merged"
	}
	httpx.OK(w, msg, row)
}

// DELETE /api/items/{id}
func (h *Handler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id, err := idFromPath("/api/items/", r.URL.Path)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	repo := Repo{DB: h.DB}
	ok, err := repo.DeleteOne(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "delete failed")
		return
	}
	if !ok {
		httpx.Error(w, http.StatusNotFound, "not found")
		return
	}
	httpx.OK(w, "deleted", nil)
}

// DELETE /api/items  (clear all)
func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	repo := Repo{DB: h.DB}
	if err := repo.DeleteAll(r.Context()); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "clear failed")
		return
	}
	httpx.OK(w, "cleared", nil)
}

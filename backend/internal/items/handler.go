package items

import (
	"context"
	"log"
	"net/http"
	"strings"
	"encoding/json"

	"techverin-backend/internal/httpx"
	"techverin-backend/internal/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB     *pgxpool.Pool
	Origin string
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /api/items/bulk
// - validate
// - insert
// - return STATS FOR THIS SUBMISSION ONLY
func (h *Handler) BulkInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.BulkItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if len(req.Items) == 0 {
		http.Error(w, "items cannot be empty", http.StatusBadRequest)
		return
	}

	// validate + normalize
	rows := make([]Row, 0, len(req.Items))
	for _, it := range req.Items {
		name := strings.TrimSpace(it.Name)
		if name == "" {
			http.Error(w, "item name cannot be empty", http.StatusBadRequest)
			return
		}
		if it.Price < 0 || isNaN(it.Price) || isInf(it.Price) {
			http.Error(w, "price must be valid and >= 0", http.StatusBadRequest)
			return
		}
		if it.Quantity < 1 {
			http.Error(w, "quantity must be >= 1", http.StatusBadRequest)
			return
		}
		rows = append(rows, Row{Name: name, Price: round2(it.Price), Quantity: it.Quantity})
	}

	ctx := r.Context()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		log.Printf("BEGIN error: %v", err)
		http.Error(w, "db transaction error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	repo := Repo{DB: h.DB}
	if err := repo.InsertBatch(ctx, tx, rows); err != nil {
		log.Printf("BATCH error: %v", err)
		http.Error(w, "db insert error", http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(ctx); err != nil {
		log.Printf("COMMIT error: %v", err)
		http.Error(w, "db commit error", http.StatusInternalServerError)
		return
	}

	// batch-only stats
	lineCount, totalQty, totalCost, avgUnit, avgLine := ComputeBatchStats(rows)
	httpx.WriteJSON(w, http.StatusOK, &types.StatsResponse{
		LineItemCount: lineCount,
		TotalQuantity: totalQty,
		TotalCost:     totalCost,
		AvgUnitPrice:  avgUnit,
		AvgLineCost:   avgLine,
	})
}

// GET /api/stats (optional global verification)
func (h *Handler) GlobalStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	repo := Repo{DB: h.DB}
	line, qty, total, avgUnit, avgLine, err := repo.GlobalStats(context.Background())
	if err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, &types.StatsResponse{
		LineItemCount: line, TotalQuantity: qty, TotalCost: total,
		AvgUnitPrice: avgUnit, AvgLineCost: avgLine,
	})
}

package items

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"techverin-backend/internal/types"
	"techverin-backend/internal/httpx"
)

type Handler struct {
	Repo *Repo
}

func NewHandler(r *Repo) *Handler { return &Handler{Repo: r} }

// POST /api/items/bulk  -> save batch (upsert) + return { batch, global }
func (h *Handler) Bulk(w http.ResponseWriter, r *http.Request) {
	var req types.BulkItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, types.Envelope{"message": "invalid JSON"})
		return
	}
	if len(req.Items) == 0 {
		httpx.WriteJSON(w, http.StatusBadRequest, types.Envelope{"message": "items cannot be empty"})
		return
	}

	// Normalize minimally
	rows := make([]types.Item, 0, len(req.Items))
	for _, it := range req.Items {
		name := strings.TrimSpace(it.Name)
		if name == "" || it.Price < 0 || it.Quantity < 1 {
			httpx.WriteJSON(w, http.StatusBadRequest, types.Envelope{"message": "invalid item(s): name/price/quantity"})
			return
		}
		rows = append(rows, types.Item{Name: name, Price: it.Price, Quantity: it.Quantity})
	}

	// Batch (pre-save) stats
	batch := ComputeBatchStats(rows)

	// Save/merge
	if err := h.Repo.InsertBatchUpsert(r.Context(), rows); err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, types.Envelope{"message": "db insert error"})
		return
	}

	// Global stats (after save)
	global, err := h.Repo.GlobalStats(r.Context())
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, types.Envelope{"message": "failed to compute global stats"})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, types.Envelope{
		"message": "saved",
		"batch":   batch,
		"global":  global,
	})
}

// GET /api/stats -> global stats only
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	s, err := h.Repo.GlobalStats(r.Context())
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, types.Envelope{"message": "failed to get stats"})
		return
	}
	httpx.WriteJSON(w, http.StatusOK, types.Envelope{"data": s})
}

// GET /api/items -> list items
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Repo.List(r.Context(), 100)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, types.Envelope{"message": "failed to list"})
		return
	}
	httpx.WriteJSON(w, http.StatusOK, types.Envelope{"data": items})
}

// PUT /api/items/{id} -> edit (with merge on conflict)
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/items/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpx.WriteJSON(w, http.StatusBadRequest, types.Envelope{"message": "invalid id"})
		return
	}
	var body struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, types.Envelope{"message": "invalid JSON"})
		return
	}
	body.Name = strings.TrimSpace(body.Name)
	if body.Name == "" || body.Price < 0 || body.Quantity < 1 {
		httpx.WriteJSON(w, http.StatusBadRequest, types.Envelope{"message": "invalid fields"})
		return
	}

	item, msg, err := h.Repo.Update(r.Context(), id, body.Name, body.Price, body.Quantity)
	if err != nil {
		if err.Error() == "not found" {
			httpx.WriteJSON(w, http.StatusNotFound, types.Envelope{"message": "not found"})
			return
		}
		httpx.WriteJSON(w, http.StatusInternalServerError, types.Envelope{"message": "update error"})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, types.Envelope{"message": msg, "data": item})
}

// DELETE /api/items/{id}
func (h *Handler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/items/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpx.WriteJSON(w, http.StatusBadRequest, types.Envelope{"message": "invalid id"})
		return
	}
	ok, err := h.Repo.DeleteOne(r.Context(), id)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, types.Envelope{"message": "delete error"})
		return
	}
	if !ok {
		httpx.WriteJSON(w, http.StatusNotFound, types.Envelope{"message": "not found"})
		return
	}
	httpx.WriteJSON(w, http.StatusOK, types.Envelope{"message": "deleted"})
}

// DELETE /api/items -> clear all
func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if err := h.Repo.ClearAll(r.Context()); err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, types.Envelope{"message": "clear error"})
		return
	}
	httpx.WriteJSON(w, http.StatusOK, types.Envelope{"message": "cleared"})
}

// (Optional) POST /api/items/preview -> batch stats without saving
func (h *Handler) Preview(w http.ResponseWriter, r *http.Request) {
	var req types.BulkItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, types.Envelope{"message": "invalid JSON"})
		return
	}
	stats := ComputeBatchStats(req.Items)
	httpx.WriteJSON(w, http.StatusOK, types.Envelope{
		"message": "ok",
		"data":    stats,
	})
}

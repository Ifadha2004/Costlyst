package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"techverin-backend/internal/db"
	"techverin-backend/internal/httpx"
	"techverin-backend/internal/items"
)

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	// Load env (exported in shell; simple is fine for the take-home)
	dsn := getenv("DATABASE_URL", "postgres://postgres:postgres@127.0.0.1:5432/techverin?sslmode=disable")
	port := getenv("PORT", "8080")
	origin := getenv("FRONTEND_ORIGIN", "http://localhost:5173")

	// DB pool + ping
	pool, err := db.NewPool(context.Background(), dsn)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	defer pool.Close()

	// Handlers
	h := &items.Handler{
		DB:     pool,
		Origin: origin,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Health)

	// Stats
	mux.HandleFunc("GET /api/stats",         httpx.WithCORS(origin, h.GlobalStats))

	// Items
	mux.HandleFunc("GET /api/items",         httpx.WithCORS(origin, h.List))
	mux.HandleFunc("PUT /api/items/",        httpx.WithCORS(origin, h.Update))     // expects /api/items/{id}
	mux.HandleFunc("DELETE /api/items/",     httpx.WithCORS(origin, h.DeleteOne))  // expects /api/items/{id}

	// Batch save + (optional) preview
	mux.HandleFunc("POST /api/items/bulk",   httpx.WithCORS(origin, h.BulkInsert))
	mux.HandleFunc("POST /api/items/preview",httpx.WithCORS(origin, h.Preview))

	// CORS preflight catch-all
	mux.HandleFunc("OPTIONS /", func(w http.ResponseWriter, r *http.Request) {
		httpx.ApplyCORSHeaders(origin, w, r); w.WriteHeader(http.StatusNoContent)
	})

	log.Printf("✅ Backend listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

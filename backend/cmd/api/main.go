package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"techverin-backend/internal/items"
)

func main() {
	ctx := context.Background()

	dbURL := getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/techverin?sslmode=disable")
	port := getenv("PORT", "8080")
	origin := getenv("FRONTEND_ORIGIN", "http://localhost:5173")

	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil { log.Fatalf("Invalid DATABASE_URL: %v", err) }
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil { log.Fatalf("Database connection error: %v", err) }
	defer db.Close()

	repo := items.NewRepo(db)
	h := items.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealth)

	// Items APIs
	mux.HandleFunc("POST /api/items/bulk", withCORS(origin, h.Bulk))
	mux.HandleFunc("POST /api/items/preview", withCORS(origin, h.Preview))  // optional
	mux.HandleFunc("GET /api/items", withCORS(origin, h.List))
	mux.HandleFunc("PUT /api/items/", withCORS(origin, h.Update))           // expects /api/items/{id}
	mux.HandleFunc("DELETE /api/items/", withCORS(origin, h.DeleteOne))     // expects /api/items/{id}
	mux.HandleFunc("DELETE /api/items", withCORS(origin, h.DeleteAll))

	// Stats
	mux.HandleFunc("GET /api/stats", withCORS(origin, h.Stats))

	// Preflight for all
	mux.HandleFunc("OPTIONS /", func(w http.ResponseWriter, r *http.Request) {
		applyCORSHeaders(origin, w, r)
		w.WriteHeader(http.StatusNoContent)
	})

	log.Printf("✅ Backend listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

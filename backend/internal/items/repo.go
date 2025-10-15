package items

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"techverin-backend/internal/types"
)

type Repo struct{ DB *pgxpool.Pool }

func NewRepo(db *pgxpool.Pool) *Repo { return &Repo{DB: db} }

// ----- Bulk upsert (merges duplicates on (name_lc, price)) -----
func (r *Repo) InsertBatchUpsert(ctx context.Context, rows []types.Item) error {
	b := &pgxpool.Batch{}
	for _, v := range rows {
		b.Queue(`
			INSERT INTO items(name, price, quantity)
			VALUES ($1,$2,$3)
			ON CONFLICT (name_lc, price)
			DO UPDATE SET quantity = items.quantity + EXCLUDED.quantity
		`, v.Name, v.Price, v.Quantity)
	}
	br := r.DB.SendBatch(ctx, b)
	return br.Close()
}

// ----- Global stats (DB) -----
func (r *Repo) GlobalStats(ctx context.Context) (types.StatsResponse, error) {
	var s types.StatsResponse
	err := r.DB.QueryRow(ctx, `
		SELECT
			COALESCE(COUNT(*),0),
			COALESCE(SUM(quantity),0),
			COALESCE(SUM(line_total),0),
			CASE WHEN COUNT(*)=0 THEN 0 ELSE SUM(price)/COUNT(*) END,
			CASE WHEN COUNT(*)=0 THEN 0 ELSE SUM(line_total)/COUNT(*) END
		FROM items;
	`).Scan(&s.LineItemCount, &s.TotalQuantity, &s.TotalCost, &s.AvgUnitPrice, &s.AvgLineCost)
	if err != nil {
		return types.StatsResponse{}, err
	}
	// round to 2dp where needed (optional; your SQL can also format)
	s.TotalCost = float64(int64(s.TotalCost*100+0.5)) / 100
	s.AvgUnitPrice = float64(int64(s.AvgUnitPrice*100+0.5)) / 100
	s.AvgLineCost = float64(int64(s.AvgLineCost*100+0.5)) / 100
	return s, nil
}

// ----- List items -----
func (r *Repo) List(ctx context.Context, limit int) ([]types.DBItem, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	rows, err := r.DB.Query(ctx, `
		SELECT id, name, price, quantity, created_at
		FROM items
		ORDER BY id DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []types.DBItem
	for rows.Next() {
		var it types.DBItem
		if err := rows.Scan(&it.ID, &it.Name, &it.Price, &it.Quantity, &it.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// ----- Update one (with merge on unique conflict) -----
func (r *Repo) Update(ctx context.Context, id int64, name string, price float64, quantity int) (types.DBItem, string, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return types.DBItem{}, "", err
	}
	defer tx.Rollback(ctx)

	// Is there another row with same (name_lc, price)?
	var dstID int64
	err = tx.QueryRow(ctx, `
		SELECT id FROM items
		WHERE name_lc = LOWER(TRIM($1)) AND price = $2 AND id <> $3
		LIMIT 1
	`, name, price, id).Scan(&dstID)

	if err == nil && dstID > 0 {
		// Merge quantities into dstID, delete source
		if _, err := tx.Exec(ctx, `UPDATE items SET quantity = quantity + $1 WHERE id = $2`, quantity, dstID); err != nil {
			return types.DBItem{}, "", err
		}
		if _, err := tx.Exec(ctx, `DELETE FROM items WHERE id = $1`, id); err != nil {
			return types.DBItem{}, "", err
		}
		var it types.DBItem
		err = tx.QueryRow(ctx, `SELECT id, name, price, quantity, created_at FROM items WHERE id=$1`, dstID).
			Scan(&it.ID, &it.Name, &it.Price, &it.Quantity, &it.CreatedAt)
		if err != nil {
			return types.DBItem{}, "", err
		}
		if err := tx.Commit(ctx); err != nil {
			return types.DBItem{}, "", err
		}
		return it, "merged", nil
	}

	// Otherwise, simple UPDATE on this id
	res, err := tx.Exec(ctx, `
		UPDATE items SET name=$1, price=$2, quantity=$3 WHERE id=$4
	`, name, price, quantity, id)
	if err != nil {
		return types.DBItem{}, "", err
	}
	if res.RowsAffected() == 0 {
		return types.DBItem{}, "", errors.New("not found")
	}

	var it types.DBItem
	err = tx.QueryRow(ctx, `SELECT id, name, price, quantity, created_at FROM items WHERE id=$1`, id).
		Scan(&it.ID, &it.Name, &it.Price, &it.Quantity, &it.CreatedAt)
	if err != nil {
		return types.DBItem{}, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return types.DBItem{}, "", err
	}
	return it, "updated", nil
}

// ----- Delete one -----
func (r *Repo) DeleteOne(ctx context.Context, id int64) (bool, error) {
	res, err := r.DB.Exec(ctx, `DELETE FROM items WHERE id=$1`, id)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

// ----- Clear all -----
func (r *Repo) ClearAll(ctx context.Context) error {
	_, err := r.DB.Exec(ctx, `TRUNCATE TABLE items RESTART IDENTITY`)
	return err
}

package items

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct{ DB *pgxpool.Pool }

func (r *Repo) InsertBatch(ctx context.Context, tx pgx.Tx, rows []Row) error {
	b := &pgx.Batch{}
	for _, v := range rows {
		b.Queue(`
		  INSERT INTO items(name, price, quantity)
		  VALUES ($1,$2,$3)
		  ON CONFLICT (name_lc, price)
		  DO UPDATE SET quantity = items.quantity + EXCLUDED.quantity
		`, v.Name, v.Price, v.Quantity)
	}
	return tx.SendBatch(ctx, b).Close()
}

type ItemRow struct {
	ID        int64
	Name      string
	Price     float64
	Quantity  int
	CreatedAt string
}

func (r *Repo) List(ctx context.Context, limit int) ([]ItemRow, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, name, price, quantity, TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SSOF')
		FROM items
		ORDER BY id DESC
		LIMIT $1
	`, limit)
	if err != nil { return nil, err }
	defer rows.Close()

	var out []ItemRow
	for rows.Next() {
		var it ItemRow
		if err := rows.Scan(&it.ID, &it.Name, &it.Price, &it.Quantity, &it.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *Repo) DeleteOne(ctx context.Context, id int64) (bool, error) {
	ct, err := r.DB.Exec(ctx, `DELETE FROM items WHERE id=$1`, id)
	return ct.RowsAffected() == 1, err
}

func (r *Repo) DeleteAll(ctx context.Context) error {
	_, err := r.DB.Exec(ctx, `TRUNCATE TABLE items`)
	return err
}

// UpdateWithMerge updates item {id}. If (name,price) collides with another row,
// it merges quantities (add) and deletes the source row — all in a single tx.
// Returns the final (possibly merged) row and whether a merge happened.
func (r *Repo) UpdateWithMerge(ctx context.Context, id int64, name string, price float64, qty int) (ItemRow, bool, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil { return ItemRow{}, false, err }
	defer tx.Rollback(ctx)

	// 1) Ensure the source row exists (and lock it)
	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM items WHERE id=$1 FOR UPDATE)`, id).Scan(&exists); err != nil {
		return ItemRow{}, false, err
	}
	if !exists { return ItemRow{}, false, pgx.ErrNoRows }

	// 2) Does a different target row already exist for (name_lc, price)?
	var targetID *int64
	if err := tx.QueryRow(ctx, `
	  SELECT id FROM items
	  WHERE name_lc = LOWER(TRIM($1)) AND price = $2 AND id <> $3
	  LIMIT 1
	`, name, price, id).Scan(&targetID); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return ItemRow{}, false, err
	}
	merged := targetID != nil

	// 3) Delete the source row
	if _, err := tx.Exec(ctx, `DELETE FROM items WHERE id=$1`, id); err != nil {
		return ItemRow{}, false, err
	}

	// 4) Re-insert with upsert to merge quantities if needed
	var out ItemRow
	if err := tx.QueryRow(ctx, `
	  INSERT INTO items(name, price, quantity)
	  VALUES ($1,$2,$3)
	  ON CONFLICT (name_lc, price)
	  DO UPDATE SET quantity = items.quantity + EXCLUDED.quantity
	  RETURNING id, name, price, quantity, TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SSOF')
	`, name, price, qty).Scan(&out.ID, &out.Name, &out.Price, &out.Quantity, &out.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return ItemRow{}, false, fmt.Errorf("pg: %s (%s)", pgErr.Message, pgErr.Code)
		}
		return ItemRow{}, false, err
	}

	if err := tx.Commit(ctx); err != nil { return ItemRow{}, false, err }
	return out, merged, nil
}

func (r *Repo) GlobalStats(ctx context.Context) (lineCount, totalQty int64, totalCost, avgUnit, avgLine float64, err error) {
	var totalStr, avgUnitStr, avgLineStr string
	err = r.DB.QueryRow(ctx, `
		SELECT
			COALESCE(COUNT(*),0),
			COALESCE(SUM(quantity),0),
			COALESCE(TO_CHAR(COALESCE(SUM(line_total),0), 'FM9999999990.00'), '0.00'),
			CASE WHEN COALESCE(SUM(quantity),0)=0 THEN '0.00'
				 ELSE TO_CHAR(SUM(line_total)/SUM(quantity), 'FM9999999990.00') END,
			CASE WHEN COUNT(*)=0 THEN '0.00'
				 ELSE TO_CHAR(SUM(line_total)/COUNT(*), 'FM9999999990.00') END
		FROM items
	`).Scan(&lineCount, &totalQty, &totalStr, &avgUnitStr, &avgLineStr)
	if err != nil { return }

	parse := func(s string) float64 { var f float64; _, _ = fmt.Sscanf(s, "%f", &f); return f }
	totalCost = parse(totalStr); avgUnit = parse(avgUnitStr); avgLine = parse(avgLineStr)
	return
}

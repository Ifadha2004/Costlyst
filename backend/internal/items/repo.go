package items

import (
	"context"

	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct{ DB *pgxpool.Pool }

func (r *Repo) InsertBatch(ctx context.Context, tx pgx.Tx, rows []Row) error {
	b := &pgx.Batch{}
	for _, v := range rows {
		b.Queue(`
			INSERT INTO items (name, price, quantity)
			VALUES ($1, $2, $3)
			ON CONFLICT (name_lc, price)
			DO UPDATE SET quantity = items.quantity + EXCLUDED.quantity
			`, v.Name, v.Price, v.Quantity)
	}
	return tx.SendBatch(ctx, b).Close()
}

func (r *Repo) GlobalStats(ctx context.Context) (lineCount, totalQty int64, totalCost, avgUnit, avgLine float64, err error) {
	var totalStr, avgUnitStr, avgLineStr string
	err = r.DB.QueryRow(ctx, `
		SELECT
			COALESCE(COUNT(*),0),
			COALESCE(SUM(quantity),0),
			COALESCE(TO_CHAR(COALESCE(SUM(line_total),0), 'FM9999999990.00'), '0.00'),
			CASE WHEN COALESCE(SUM(quantity),0)=0 THEN '0.00'
				ELSE TO_CHAR(SUM(line_total)/SUM(quantity), 'FM9999999990.00')
			END,
			CASE WHEN COUNT(*)=0 THEN '0.00'
				ELSE TO_CHAR(SUM(line_total)/COUNT(*), 'FM9999999990.00')
			END
		FROM items;
	`).Scan(&lineCount, &totalQty, &totalStr, &avgUnitStr, &avgLineStr)
	if err != nil {
		return
	}
	parse := func(s string) float64 {
		// safe, no commas because of FM format, but keep simple:
		var f float64
		_, _ = fmt.Sscanf(s, "%f", &f)
		return f
	}
	totalCost = parse(totalStr)
	avgUnit = parse(avgUnitStr)
	avgLine = parse(avgLineStr)
	return
}

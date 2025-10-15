package items

func round2(f float64) float64 { return float64(int64(f*100+0.5)) / 100 }
func isNaN(f float64) bool     { return f != f }
func isInf(f float64) bool     { return f > 1e308 || f < -1e308 }

type Row struct {
	Name     string
	Price    float64
	Quantity int
}

// Compute stats for the CURRENT submission (batch-only).
func ComputeBatchStats(rows []Row) (lineCount int64, totalQty int64, totalCost, avgUnit, avgLine float64) {
	lineCount = int64(len(rows))
	for _, r := range rows {
		totalQty += int64(r.Quantity)
		totalCost += float64(r.Quantity) * r.Price
	}
	totalCost = round2(totalCost)
	if totalQty > 0 {
		avgUnit = round2(totalCost / float64(totalQty))
	}
	if lineCount > 0 {
		avgLine = round2(totalCost / float64(lineCount))
	}
	return
}

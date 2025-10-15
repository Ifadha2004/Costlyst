package items

import (
	"math"

	"techverin-backend/internal/types"
)

func round2(f float64) float64 { return math.Round(f*100) / 100 }

func ComputeBatchStats(items []types.Item) types.StatsResponse {
	var (
		count    int64
		totalQty int64
		sumCost  float64
		sumUnit  float64
	)

	for _, it := range items {
		if it.Name == "" || it.Price < 0 || it.Quantity < 1 {
			continue
		}
		count++
		totalQty += int64(it.Quantity)
		sumCost += float64(it.Quantity) * it.Price
		sumUnit += it.Price
	}

	var avgUnit, avgLine float64
	if count > 0 {
		avgUnit = sumUnit / float64(count)
		avgLine = sumCost / float64(count)
	}

	return types.StatsResponse{
		LineItemCount: count,
		TotalQuantity: totalQty,
		TotalCost:     round2(sumCost),
		AvgUnitPrice:  round2(avgUnit),
		AvgLineCost:   round2(avgLine),
	}
}

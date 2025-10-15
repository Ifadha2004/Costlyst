package types

import "time"

type Envolpe map[string]any

type DBItem struct {
	ID		int64     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

type BulkResult struct {
	Batch StatsResponse `json:"batch"`
	Global StatsResponse `json:"global"`
}
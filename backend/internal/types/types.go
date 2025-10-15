package types

type Item struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type BulkItemsRequest struct {
	Items []Item `json:"items"`
}

type StatsResponse struct {
	LineItemCount int64   `json:"lineItemCount"`
	TotalQuantity int64   `json:"totalQuantity"`
	TotalCost     float64 `json:"totalCost"`
	AvgUnitPrice  float64 `json:"avgUnitPrice"`
	AvgLineCost   float64 `json:"avgLineCost"`
}

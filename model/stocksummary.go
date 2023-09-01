package model

type StockSummary struct {
	StockCode     string  `json:"stock_code"`
	PreviousPrice float64 `json:"previous_price"`
	OpenPrice     float64 `json:"open_price"`
	HighestPrice  float64 `json:"highest_price"`
	LowestPrice   float64 `json:"lowest_price"`
	ClosePrice    float64 `json:"close_price"`
	Transaction   float64 `json:"transaction"`
	AveragePrice  float64 `json:"average_price"`
	Volume        float64 `json:"volume"`
	Value         float64 `json:"value"`
}

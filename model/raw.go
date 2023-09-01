package model

type Raw struct {
	Type             string `json:"type"`
	OrderBook        string `json:"order_book"`
	Quantity         string `json:"quantity"`
	Price            string `json:"price"`
	StockCode        string `json:"stock_code"`
	ExecutedQuantity string `json:"executed_quantity"`
	ExecutionPrice   string `json:"execution_price"`
	OrderNumber      string `json:"order_number"`
	OrderVerb        string `json:"order_verb"`
}

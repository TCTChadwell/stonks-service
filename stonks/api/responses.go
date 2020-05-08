package api

type PriceApiResponse struct {
	QuoteResponse QuoteResponse `json:"quoteResponse"`
}

type QuoteResponse struct {
	Result []Stonk     `json:"result"`
	Error  interface{} `json:"error"`
}

type Stonk struct {
	Symbol string  `json:"symbol"`
	Source string  `json:"quoteType"`
	Price  float64 `json:"regularMarketPrice"`
	High   float64 `json:"fiftyTwoWeekHigh"`
	Low    float64 `json:"fiftyTwoWeekLow"`
}

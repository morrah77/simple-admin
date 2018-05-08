package common

type Record struct {
	Id           string  `json:"id,omitempty"`
	CurrencyFrom string  `json:"currency_from"`
	CurrencyTo   string  `json:"currency_to"`
	Rate         float64 `json:"rate"`
	Time         float64 `json:"time"`
}

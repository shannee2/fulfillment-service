package model

type Order struct {
	OrderID     int64       `json:"orderId"`
	TotalAmount TotalAmount `json:"grandTotal"`
	//UserID       string      `json:"userId"`
	CurrencyType string     `json:"currencyType"`
	Restaurant   Restaurant `json:"restaurant"`
}

type Restaurant struct {
	Name      string  `json:"name"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type TotalAmount struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

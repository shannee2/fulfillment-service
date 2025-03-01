package model

type DeliveryStatus string

const (
	Available       DeliveryStatus = "AVAILABLE"
	DeliveringOrder DeliveryStatus = "DELIVERING_ORDER"
	Unavailable     DeliveryStatus = "UNAVAILABLE"
)

type DeliveryPartner struct {
	ID     int            `json:"id"`
	Name   string         `json:"name"`
	Phone  string         `json:"phone"`
	City   string         `json:"city"`
	Status DeliveryStatus `json:"status"`
}

package model

import "time"

type DeliveryAssignment struct {
	ID                int       `json:"id"`
	DeliveryPartnerID int       `json:"delivery_partner_id"`
	OrderID           int64     `json:"order_id"`
	AssignedAt        time.Time `json:"assigned_at"`
}

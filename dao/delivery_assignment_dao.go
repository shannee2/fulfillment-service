package dao

import (
	"database/sql"
	"fulfillment/model"
	"log"
)

type DeliveryAssignmentDAOInterface interface {
	InsertDeliveryAssignment(assignment model.DeliveryAssignment) (int, error)
}

type DeliveryAssignmentDAO struct {
	db *sql.DB
}

func NewDeliveryAssignmentDAO(db *sql.DB) *DeliveryAssignmentDAO {
	return &DeliveryAssignmentDAO{db: db}
}

// Insert a new delivery assignment
func (dao *DeliveryAssignmentDAO) InsertDeliveryAssignment(assignment model.DeliveryAssignment) (int, error) {
	query := `INSERT INTO delivery_assignment (delivery_partner_id, order_id, assigned_at) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := dao.db.QueryRow(query, assignment.DeliveryPartnerID, assignment.OrderID, assignment.AssignedAt).Scan(&id)
	if err != nil {
		log.Println("Error inserting delivery assignment:", err)
		return 0, err
	}
	return id, nil
}

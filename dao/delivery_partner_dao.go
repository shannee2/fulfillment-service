package dao

import (
	"database/sql"
	"fulfillment/model"
	"log"
)

type DeliveryPartnerDAOInterface interface {
	InsertDeliveryPartner(partner model.DeliveryPartner) error
	UpdateDeliveryPartnerStatus(partnerID int, status string) error
	GetAvailableDeliveryPartnersByCity(city string) ([]model.DeliveryPartner, error)
}

type DeliveryPartnerDAO struct {
	db *sql.DB
}

func NewDeliveryPartnerDAO(db *sql.DB) *DeliveryPartnerDAO {
	return &DeliveryPartnerDAO{db: db}
}

func (dao *DeliveryPartnerDAO) InsertDeliveryPartner(partner model.DeliveryPartner) error {
	query := `INSERT INTO delivery_partner (name, phone, city, status) VALUES ($1, $2, $3, $4)`
	_, err := dao.db.Exec(query, partner.Name, partner.Phone, partner.City, partner.Status)
	if err != nil {
		log.Printf("Error inserting delivery partner: %v", err)
	}
	return err
}

// Update delivery partner status
func (dao *DeliveryPartnerDAO) UpdateDeliveryPartnerStatus(partnerID int, status string) error {
	query := `UPDATE delivery_partner SET status = $1 WHERE id = $2`
	_, err := dao.db.Exec(query, status, partnerID)
	if err != nil {
		log.Println("Error updating delivery partner status:", err)
		return err
	}
	return nil
}

func (dao *DeliveryPartnerDAO) GetAvailableDeliveryPartnersByCity(city string) ([]model.DeliveryPartner, error) {
	query := `SELECT id, name, phone, city, status FROM delivery_partner WHERE city = $1 AND status = 'AVAILABLE'`
	rows, err := dao.db.Query(query, city)
	if err != nil {
		log.Printf("Error fetching delivery partners by city: %v", err)
		return nil, err
	}
	defer rows.Close()

	var partners []model.DeliveryPartner
	for rows.Next() {
		var partner model.DeliveryPartner
		err := rows.Scan(&partner.ID, &partner.Name, &partner.Phone, &partner.City, &partner.Status)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		partners = append(partners, partner)
	}
	return partners, nil
}

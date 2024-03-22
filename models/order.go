package models

import "time"

type Order struct {
	ID           int       `gorm:"column:order_id; primaryKey" json:"id"`
	CustomerName string    `gorm:"column:customer_name" json:"customer_name"`
	OrderedAt    time.Time `gorm:"column:ordered_at" json:"ordered_at"`
}

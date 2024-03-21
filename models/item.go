package models

type Item struct {
	ID          int    `gorm:"column:item_id;primaryKey" json:"id"`
	Code        string `gorm:"column:item_code" json:"code"`
	Description string `gorm:"type:text" json:"description"`
	Quantity    int    `json:"quantity"`
	OrderID     int    `gorm:"column:order_id" json:"orderId"`
}

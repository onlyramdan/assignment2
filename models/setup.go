package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	database, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/go_orders_by"))

	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&Item{}, &Order{})

	DB = database
}

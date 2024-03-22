package main

import (
	"assg2/controllers/order_controller"
	"assg2/models"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	models.ConnectDB()

	r.POST("/order", order_controller.Create)
	r.GET("/orders", order_controller.Index)
	r.PUT("/order/:orderId", order_controller.Update)
	r.DELETE("/order/:orderId", order_controller.Delete)
	r.GET("/order/:orderId")

	r.Run("localhost:5000")
}

package order_controller

import (
	"assg2/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Index(c *gin.Context) {

	var orders []models.Order

	models.DB.Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"order": orders,
	})
}

func Create(c *gin.Context) {
	var reqBody struct {
		OrderedAt    time.Time `json:"orderedAt"`
		CustomerName string    `json:"customerName"`
		Items        []struct {
			ItemCode    string `json:"itemCode"`
			Description string `json:"description"`
			Quantity    int    `json:"quantity"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	order := models.Order{
		CustomerName: reqBody.CustomerName,
		OrderedAt:    reqBody.OrderedAt,
	}

	if err := models.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan order",
		})
		return
	}

	for _, itemData := range reqBody.Items {
		item := models.Item{
			Code:        itemData.ItemCode,
			Description: itemData.Description,
			Quantity:    itemData.Quantity,
			OrderID:     order.ID,
		}

		if err := models.DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal menyimpan item",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order berhasil dibuat",
		"order":   order,
	})
}

func Update(c *gin.Context) {

}

func Delete(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, err := strconv.Atoi(orderIDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID pesanan tidak valid",
		})
		return
	}

	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Pesanan tidak ditemukan",
		})
		return
	}

	if err := models.DB.Where("order_id = ?", orderID).Delete(&models.Item{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus item terkait",
		})
		return
	}

	if err := models.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus pesanan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pesanan berhasil dihapus",
	})

}

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
	orderIDStr := c.Param("orderId")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	var reqBody struct {
		LineItemID   int    `json:"lineItemId"`
		OrderedAt    string `json:"orderedAt"` // Mengambil sebagai string untuk validasi
		CustomerName string `json:"customerName"`
		Items        []struct {
			ItemCode    string `json:"itemCode"`
			Description string `json:"description"`
			Quantity    int    `json:"quantity"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Parse orderedAt
	orderedAt, err := time.Parse(time.RFC3339, reqBody.OrderedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid orderedAt format",
		})
		return
	}

	// Retrieve the order from the database based on orderID
	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found",
		})
		return
	}

	// Begin a database transaction
	tx := models.DB.Begin()

	// Update order details
	order.OrderedAt = orderedAt
	order.CustomerName = reqBody.CustomerName

	// Delete existing items related to the order
	if err := tx.Where("order_id = ?", order.ID).Delete(&models.Item{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete existing items",
		})
		return
	}

	// Create new items related to the order
	var items []models.Item
	for _, itemData := range reqBody.Items {
		item := models.Item{
			Code:        itemData.ItemCode,
			Description: itemData.Description,
			Quantity:    itemData.Quantity,
			OrderID:     order.ID,
		}
		items = append(items, item)
	}

	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create new items",
		})
		return
	}

	// Commit the transaction
	tx.Commit()

	// Save the updated order to the database
	if err := models.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update the order",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order successfully updated",
		"order":   order,
	})
}

func Delete(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, err := strconv.Atoi(orderIDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID Pesanan Tidak Valid",
		})
		return
	}

	tx := models.DB.Begin()

	if err := tx.Where("order_id = ?", orderID).Delete(&models.Item{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus item terkait",
		})
		return
	}

	if err := tx.Where("order_id = ?", orderID).Delete(&models.Order{}).Error; err != nil {
		tx.Rollback() // Batalkan transaksi jika terjadi kesalahan
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus pesanan",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Pesanan berhasil dihapus",
	})

}

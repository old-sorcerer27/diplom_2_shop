package handlers

import (
	"net/http"
	"store/database"
	"store/database/source"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllUsersWithCarts(c *gin.Context) {
	db := database.GetDB()

	var users []source.User
	// Загружаем пользователей с их корзинами и товарами в корзине
	if err := db.Preload("Cart.Items.Product").Where("role != ?", source.RoleOwner).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Скрываем пароли
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

// func GetUserCart(c *gin.Context) {
// 	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
// 		return
// 	}

// 	db := database.GetDB()

// 	var cart source.Cart
// 	if err := db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, cart)
// }

func GetAllOrders(c *gin.Context) {
	db := database.GetDB()

	var orders []source.Order
	if err := db.Preload("User").Preload("Items.Product").Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func UpdateOrderStatus(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending paid shipped delivered cancelled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()

	result := db.Model(&source.Order{}).Where("id = ?", orderID).Update("status", req.Status)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

func UpdateProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var updates source.Product
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()

	var product source.Product
	if err := db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Обновляем только переданные поля
	if updates.Name != "" {
		product.Name = updates.Name
	}
	if updates.Description != "" {
		product.Description = updates.Description
	}
	if updates.Price > 0 {
		product.Price = updates.Price
	}
	if updates.Stock >= 0 {
		product.Stock = updates.Stock
	}
	if updates.ImageURL != "" {
		product.ImageURL = updates.ImageURL
	}
	if updates.Category != "" {
		product.Category = updates.Category
	}

	result := db.Save(&product)
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	db := database.GetDB()
	result := db.Delete(&source.Product{}, productID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

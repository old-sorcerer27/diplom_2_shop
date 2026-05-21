package handlers

import (
	"log"
	"net/http"
	"store/database"
	"store/database/generated"
	"store/database/source"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderRequest struct {
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	Items         []struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	} `json:"items"`
}

func CreateOrder(c *gin.Context) {
	var req OrderRequest
	log.Printf("Order data %v\n", req)

	if err := c.BindJSON(&req); err != nil {
		log.Printf("Error parsing order request: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received order request: %+v\n", req)

	db := database.GetDB()

	Error := db.Transaction(func(tx *gorm.DB) error {
		var total float64
		for _, item := range req.Items {
			var product source.Product
			if err := db.First(&product, item.ProductID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
					log.Printf("Error parsing order request: %v\n", err)
					return err
				}
			}
			if product.Stock < item.Quantity {
				c.JSON(http.StatusBadRequest, gin.H{"error": "not enough stock for product " + product.Name})
				log.Printf("Error parsing order request: not enough stock for product %s\n", product.Name)
				return gorm.ErrInvalidData
			}
			total += product.Price * float64(item.Quantity)
		}
		orderID, err := generated.OrderQuery[source.Order](db).CreateOrder(c.Request.Context(), req.CustomerName, req.CustomerEmail, total)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error creating order": err.Error()})
			log.Printf("Error creating order: %v\n", err)
			return err
		}
		for _, item := range req.Items {
			if err := generated.OrderQuery[source.Order](db).AddToOrder(c.Request.Context(), orderID, item.ProductID, item.Quantity); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error adding item to order": err.Error()})
				log.Printf("Error adding item to order: %v\n", err)
				return err
			}
		}
		c.JSON(http.StatusOK, gin.H{"order_id": orderID, "total": total})
		log.Printf("Order created successfully: ID=%d, Total=%.2f\n", orderID, total)
		return nil
	})
	if Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error global": Error.Error()})
		log.Printf("Transaction failed: %v\n", Error)
		return
	}

}

func GetOrders(c *gin.Context) {
	db := database.GetDB()

	orders, err := gorm.G[source.Order](db).
		Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

package handlers

import (
	"log"
	"net/http"
	"store/database"
	"store/database/generated"
	"store/database/source"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetReviews(c *gin.Context) {
	log.Println("Received request for reviews")
	db := database.GetDB()

	reviews, err := gorm.G[source.Review](db).
		Where(generated.Product.Stock.Gt(0)).
		Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func GetProductReviews(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	db := database.GetDB()

	reviews, err := generated.ReviewQuery[source.Review](db).GetByProductID(c.Request.Context(), int(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

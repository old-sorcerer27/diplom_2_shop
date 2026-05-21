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

func GetProducts(c *gin.Context) {
	log.Println("Received request for product list")
	db := database.GetDB()

	products, err := gorm.G[source.Product](db).
		Where(generated.Product.Stock.Gt(0)).
		Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductsByName(c *gin.Context) {
	log.Println("Received request to search products by name")

	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product name required"})
		log.Println("Product name query parameter is missing")
		return
	}

	db := database.GetDB()
	products, err := generated.ProductQuery[source.Product](db).SearchByName(c.Request.Context(), "%"+name+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Error searching products by name '%s': %v", name, err)
		return
	}
	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	db := database.GetDB()

	var product source.Product
	result := db.First(&product, id)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query required"})
		return
	}

	db := database.GetDB()

	products, err := generated.ProductQuery[source.Product](db).SearchByName(c.Request.Context(), "%"+query+"%")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func CreateProduct(c *gin.Context) {
	var product source.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Create(&product)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Проверяем, что запись действительно создана
	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.Header("Cache-Control", "no-cache")
	c.JSON(http.StatusCreated, product)
}

// func AddComment(c *gin.Context) {
// 	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
// 		return
// 	}

// 	userID := middleware.GetCurrentUserID(c)
// 	if userID == 0 {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login to leave a comment"})
// 		return
// 	}

// 	var req struct {
// 		Rating int    `json:"rating" binding:"required,min=1,max=5"`
// 		Text   string `json:"text" binding:"required,min=3,max=1000"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	db := database.GetDB()

// 	// Проверяем существование товара
// 	var product source.Product
// 	if err := db.First(&product, productID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
// 		return
// 	}

// 	// Проверяем, покупал ли пользователь этот товар
// 	var orderCount int64
// 	db.Table("order_items").
// 		Joins("JOIN orders ON orders.id = order_items.order_id").
// 		Where("orders.user_id = ? AND order_items.product_id = ? AND orders.status != ?",
// 			userID, productID, "cancelled").
// 		Count(&orderCount)

// 	isVerified := orderCount > 0

// 	if !isVerified {
// 		c.JSON(http.StatusForbidden, gin.H{
// 			"error": "You can only comment on products you have purchased",
// 		})
// 		return
// 	}

// 	// Проверяем, не оставлял ли пользователь уже комментарий
// 	var existingComment source.Review
// 	result := db.Where("product_id = ? AND user_id = ?", productID, userID).First(&existingComment)
// 	if result.Error == nil {
// 		c.JSON(http.StatusConflict, gin.H{
// 			"error":   "You have already reviewed this product",
// 			"comment": existingComment,
// 		})
// 		return
// 	}

// 	// Создаем комментарий
// 	comment := source.Comment{
// 		ProductID:  uint(productID),
// 		UserID:     userID,
// 		Rating:     req.Rating,
// 		Text:       req.Text,
// 		IsVerified: isVerified,
// 	}

// 	if err := db.Create(&comment).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
// 		return
// 	}

// 	// Загружаем информацию о пользователе
// 	db.Preload("User").First(&comment, comment.ID)
// 	comment.User.Password = ""

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "Comment added successfully",
// 		"comment": comment,
// 	})
// }

// UpdateComment - обновление комментария
// func UpdateComment(c *gin.Context) {
// 	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
// 		return
// 	}

// 	userID := middleware.GetCurrentUserID(c)
// 	userRole := middleware.GetUserRole(c)

// 	var req struct {
// 		Rating int    `json:"rating" binding:"required,min=1,max=5"`
// 		Text   string `json:"text" binding:"required,min=3,max=1000"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	db := database.GetDB()

// 	var comment models.Comment
// 	if err := db.First(&comment, commentID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
// 		return
// 	}

// 	// Проверяем права: только автор, админ или владелец могут редактировать
// 	if comment.UserID != userID && userRole != "admin" && userRole != "owner" {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to edit this comment"})
// 		return
// 	}

// 	// Обновляем
// 	comment.Rating = req.Rating
// 	comment.Text = req.Text

// 	if err := db.Save(&comment).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Comment updated successfully",
// 		"comment": comment,
// 	})
// }

// func DeleteComment(c *gin.Context) {
// 	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
// 		return
// 	}

// 	userID := middleware.GetCurrentUserID(c)
// 	userRole := middleware.GetUserRole(c)

// 	db := database.GetDB()

// 	var comment models.Comment
// 	if err := db.First(&comment, commentID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
// 		return
// 	}

// 	// Проверяем права
// 	if comment.UserID != userID && userRole != "admin" && userRole != "owner" {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this comment"})
// 		return
// 	}

// 	// Удаляем комментарий (мягкое удаление)
// 	if err := db.Delete(&comment).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
// }

// func LikeComment(c *gin.Context) {
// 	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
// 		return
// 	}

// 	userID := middleware.GetCurrentUserID(c)
// 	if userID == 0 {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login to like comments"})
// 		return
// 	}

// 	db := database.GetDB()

// 	// Проверяем существование комментария
// 	var comment models.Comment
// 	if err := db.First(&comment, commentID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
// 		return
// 	}

// 	// Проверяем, не лайкал ли уже пользователь
// 	var like models.CommentLike
// 	result := db.Where("comment_id = ? AND user_id = ?", commentID, userID).First(&like)

// 	if result.Error == nil {
// 		// Удаляем лайк
// 		db.Delete(&like)
// 		db.Model(&comment).Update("likes", gorm.Expr("likes - ?", 1))
// 		c.JSON(http.StatusOK, gin.H{"message": "Like removed", "likes": comment.Likes - 1})
// 		return
// 	}

// 	// Добавляем лайк
// 	like = models.CommentLike{
// 		CommentID: uint(commentID),
// 		UserID:    userID,
// 	}

// 	if err := db.Create(&like).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like comment"})
// 		return
// 	}

// 	db.Model(&comment).Update("likes", gorm.Expr("likes + ?", 1))

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Comment liked",
// 		"likes":   comment.Likes + 1,
// 	})
// }

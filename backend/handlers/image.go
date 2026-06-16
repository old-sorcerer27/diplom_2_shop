package handlers

import (
	"net/http"
	"store/database/source"
	"store/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ImageHandler struct {
	db     *gorm.DB
	imgSvc *services.ImageService
}

func NewImageHandler(db *gorm.DB) *ImageHandler {
	return &ImageHandler{
		db:     db,
		imgSvc: services.NewImageService(),
	}
}

// UploadProductImage - загрузка основного изображения товара
func (h *ImageHandler) UploadProductImage(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Получаем файл из формы
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// Находим товар
	var product source.Product
	if err := h.db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Удаляем старое изображение
	if product.ImageURL != "" {
		h.imgSvc.DeleteImage(product.ImageURL)
	}

	// Обрабатываем и сохраняем изображение
	processed, err := h.imgSvc.ProcessAndSave(src, file.Filename, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновляем товар
	updates := map[string]interface{}{
		"image_url":     processed.Original,
		"thumbnail_url": processed.Thumbnail,
		"medium_url":    processed.Medium,
	}

	if err := h.db.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"image_url":     processed.Original,
		"thumbnail_url": processed.Thumbnail,
		"medium_url":    processed.Medium,
		"message":       "Image uploaded successfully",
	})
}

// UploadProductGallery - загрузка галереи изображений
func (h *ImageHandler) UploadProductGallery(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Получаем форму с несколькими файлами
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one image is required"})
		return
	}

	if len(files) > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 5 images per upload"})
		return
	}

	// Находим товар
	var product source.Product
	if err := h.db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var imageURLs []string

	// Обрабатываем каждый файл
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			continue
		}
		defer src.Close()

		processed, err := h.imgSvc.ProcessAndSave(src, file.Filename, false)
		if err != nil {
			continue
		}

		imageURLs = append(imageURLs, processed.Original)
	}

	if len(imageURLs) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process images"})
		return
	}

	// Обновляем галерею
	gallery := product.Gallery
	gallery = append(gallery, imageURLs...)

	if err := h.db.Model(&product).Update("gallery", gallery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gallery"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"images":     imageURLs,
		"all_images": gallery,
		"message":    "Gallery updated successfully",
	})
}

// DeleteProductImage - удаление основного изображения
func (h *ImageHandler) DeleteProductImage(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product source.Product
	if err := h.db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.ImageURL != "" {
		h.imgSvc.DeleteImage(product.ImageURL)
	}

	updates := map[string]interface{}{
		"image_url":     nil,
		"thumbnail_url": nil,
		"medium_url":    nil,
	}

	if err := h.db.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image deleted successfully",
	})
}

// DeleteGalleryImage - удаление изображения из галереи
func (h *ImageHandler) DeleteGalleryImage(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req struct {
		ImageURL string `json:"image_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product source.Product
	if err := h.db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Удаляем физический файл
	h.imgSvc.DeleteImage(req.ImageURL)

	// Удаляем из галереи
	gallery := product.Gallery
	newGallery := make(source.JSONArray, 0)
	for _, url := range gallery {
		if url != req.ImageURL {
			newGallery = append(newGallery, url)
		}
	}

	if err := h.db.Model(&product).Update("gallery", newGallery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gallery"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image removed from gallery",
	})
}

// GetProductImages - получение всех изображений товара
func (h *ImageHandler) GetProductImages(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product source.Product
	if err := h.db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	images := []source.ProductImage{}

	// Добавляем основное изображение
	if product.ImageURL != "" {
		images = append(images, source.ProductImage{
			ID:        0,
			URL:       product.ImageURL,
			Thumbnail: product.ThumbnailURL,
			Medium:    product.MediumURL,
			IsMain:    true,
		})
	}

	// Добавляем изображения из галереи
	for i, url := range product.Gallery {
		images = append(images, source.ProductImage{
			ID:        uint(i + 1),
			URL:       url,
			Thumbnail: "",
			Medium:    "",
			IsMain:    false,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"images": images,
		"count":  len(images),
	})
}

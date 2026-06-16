package main

import (
	"log"
	"store/database"
	"store/handlers"
	"store/middleware"
	"store/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	if err := services.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	database.InitDB()

	services.SetupUploadDirectories()
}
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	if err := services.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	database.InitDB()
	services.SetupUploadDirectories()
}

func main() {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	db := database.GetDB()
	imageHandler := handlers.NewImageHandler(db)

	// === СТАТИЧЕСКИЕ ФАЙЛЫ (ТОЛЬКО ОДИН СПОСОБ) ===

	// ВАРИАНТ 1: Только Static (РЕКОМЕНДУЕТСЯ)
	r.Static("/uploads", "./uploads")
	r.Static("/api/v1/uploads", "./uploads")

	// ВАРИАНТ 2: Static + обработчик (если нужен fallback)
	// Раскомментируйте этот блок и закомментируйте r.Static выше
	/*
		r.StaticFS("/uploads", http.Dir("./uploads"))
		r.Any("/api/v1/uploads/*filepath", func(c *gin.Context) {
			filepath := c.Param("filepath")
			fullPath := "./uploads" + filepath

			// Проверяем существование файла
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				// Возвращаем 404 или плейсхолдер
				c.File("./uploads/products/original/placeholder.jpg")
				return
			}
			c.File(fullPath)
		})
	*/

	// === ПУБЛИЧНЫЕ МАРШРУТЫ ===
	public := r.Group("/api/v1")
	{
		public.POST("/auth/register", handlers.Register)
		public.POST("/auth/login", handlers.Login)

		public.GET("/checkconn", CheckConnection)
		public.GET("/products", handlers.GetProducts)
		public.GET("/products/search", handlers.GetProductsByName)
		public.GET("/products/:id", handlers.GetProductByID)
		public.GET("/orders", handlers.GetOrders)
		public.POST("/orders", handlers.CreateOrder)

		public.GET("/products/:id/images", imageHandler.GetProductImages)
	}

	// === ЗАЩИЩЕННЫЕ МАРШРУТЫ ===
	auth := r.Group("/api/v1")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/auth/me", handlers.GetMe)
	}

	// === АДМИН МАРШРУТЫ ===
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RequireAdminOrOwner())
	{
		admin.GET("/users", handlers.GetAllUsersWithCarts)

		admin.POST("/products", handlers.CreateProduct)
		admin.PUT("/products/:id", handlers.UpdateProduct)
		admin.DELETE("/products/:id", handlers.DeleteProduct)

		admin.POST("/products/:id/image", imageHandler.UploadProductImage)
		admin.POST("/products/:id/gallery", imageHandler.UploadProductGallery)
		admin.DELETE("/products/:id/image", imageHandler.DeleteProductImage)
		admin.DELETE("/products/:id/gallery", imageHandler.DeleteGalleryImage)

		admin.GET("/orders", handlers.GetAllOrders)
		admin.PUT("/orders/:id/status", handlers.UpdateOrderStatus)
	}

	// === ВЛАДЕЛЕЦ МАРШРУТЫ ===
	owner := r.Group("/api/v1/owner")
	owner.Use(middleware.AuthMiddleware(), middleware.RequireOwner())
	{
		owner.GET("/users", handlers.GetAllUsersWithCarts)
	}

	log.Println("Server starting on :8080")
	r.Run(":8080")
}

func CheckConnection(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Connected to backend"})
}

// package main

// import (
// 	"log"
// 	"store/database"
// 	"store/handlers"
// 	"store/middleware"
// 	"store/services"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// )

// func init() {
// 	// loads values from .env into the system
// 	if err := godotenv.Load(); err != nil {
// 		log.Print("No .env file found")
// 	}
// 	if err := services.LoadConfig(); err != nil {
// 		log.Fatal("Failed to load config:", err)
// 	}

// 	database.InitDB()

// 	services.SetupUploadDirectories()
// }

// func main() {

// 	r := gin.Default()

// 	r.Use(middleware.CORSMiddleware())

// 	db := database.GetDB()
// 	imageHandler := handlers.NewImageHandler(db)

// 	r.Use(func(c *gin.Context) {
// 		c.Header("Access-Control-Allow-Origin", "*")
// 		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		c.Header("Access-Control-Allow-Headers", "Content-Type")
// 		c.Next()
// 	})

// 	// === СТАТИЧЕСКИЕ ФАЙЛЫ ===
// 	// Основная папка для загрузок
// 	r.Static("/uploads", "./uploads")

// 	// Для совместимости с API путями
// 	r.Static("/api/v1/uploads", "./uploads")

// 	// Прямой доступ к изображениям товаров
// 	r.Static("/images/products", "./uploads/products")
// 	r.Static("/products/images", "./uploads/products")

// 	public := r.Group("/api/v1")
// 	{
// 		public.POST("/auth/register", handlers.Register)
// 		public.POST("/auth/login", handlers.Login)

// 		public.GET("/checkconn", CheckConnection)
// 		public.GET("/products", handlers.GetProducts)
// 		public.GET("/products/search", handlers.GetProductsByName)
// 		public.GET("/products/:id", handlers.GetProductByID)
// 		public.GET("/orders", handlers.GetOrders)
// 		public.POST("/orders", handlers.CreateOrder)

// 		public.GET("/products/:id/images", imageHandler.GetProductImages)
// 	}

// 	auth := r.Group("/api/v1")
// 	auth.Use(middleware.AuthMiddleware())
// 	{
// 		auth.GET("/auth/me", handlers.GetMe)
// 	}

// 	admin := r.Group("/api/v1/admin")
// 	admin.Use(middleware.AuthMiddleware(), middleware.RequireAdminOrOwner())
// 	{
// 		admin.GET("/users", handlers.GetAllUsersWithCarts)

// 		admin.POST("/admin/products", handlers.CreateProduct)
// 		admin.PUT("/products/:id", handlers.UpdateProduct)
// 			admin.PUT("/orders/:id/status", handlers.UpdateOrderStatus)
// 	}

// 	owner := r.Group("/api/v1/owner")
// 	owner.Use(middleware.AuthMiddleware(), middleware.RequireOwner())
// 	{
// 		owner.GET("/users", handlers.GetAllUsersWithCarts)
// 	}

// 	// r.Run(":" + services.AppConfig.Port)
// 	r.Run(":8080")
// }

// func CheckConnection(c *gin.Context) {
// 	c.JSON(200, gin.H{"message": "Connected to backend"})
// }
// admin.DELETE("/products/:id", handlers.DeleteProduct)

// 		admin.POST("/products/:id/image", imageHandler.UploadProductImage)
// 		admin.POST("/products/:id/gallery", imageHandler.UploadProductGallery)
// 		admin.DELETE("/products/:id/image", imageHandler.DeleteProductImage)
// 		admin.DELETE("/products/:id/gallery", imageHandler.DeleteGalleryImage)

// 		admin.GET("/orders", handlers.GetAllOrders)

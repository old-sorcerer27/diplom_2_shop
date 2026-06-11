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
	services.LoadConfig()

	database.InitDB()
}

func main() {

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	})

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
	}

	auth := r.Group("/api/v1")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/auth/me", handlers.GetMe)
	}

	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RequireAdminOrOwner())
	{
		admin.GET("/users", handlers.GetAllUsersWithCarts)
		admin.GET("/orders", handlers.GetAllOrders)
		admin.PUT("/orders/:id/status", handlers.UpdateOrderStatus)
		admin.POST("/products", handlers.CreateProduct)
		admin.PUT("/products/:id", handlers.UpdateProduct)
		admin.DELETE("/products/:id", handlers.DeleteProduct)
	}

	// r.Run(":" + services.AppConfig.Port)
	r.Run(":8080")
}

func CheckConnection(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Connected to backend"})
}

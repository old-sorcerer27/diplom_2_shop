package main

import (
	"store/database"
	"store/handlers"
	"store/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

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

	r.Run(":8080")
}

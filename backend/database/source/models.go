package source

import (
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

const (
	RoleClient = "client"
	RoleAdmin  = "admin"
	RoleOwner  = "owner"
)

type User struct {
	gorm.Model
	Email    string   `json:"email" gorm:"unique"`
	Username string   `json:"username" gorm:"unique"`
	Password string   `json:"password"`
	Phone    string   `json:"phone"`
	Role     string   `json:"role" gorm:"default:'client'"`
	IsActive bool     `json:"is_active" gorm:"default:true"`
	Reviews  []Review `gorm:"foreignKey:UserID"`
	Orders   []Order  `gorm:"foreignKey:UserID"`
}

type Product struct {
	gorm.Model
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Category      string   `gorm:"size:100;index" json:"category"`
	Price         float64  `json:"price"`
	Stock         int      `json:"stock"`
	ImageURL      string   `json:"image_url"`
	Orders        []Order  `gorm:"many2many:order_items;"`
	Reviews       []Review `gorm:"foreignKey:ProductID"`
	AverageRating float64  `json:"average_rating"`
	ReviewsCount  int      `json:"reviews_count"`
}

type Order struct {
	gorm.Model
	UserID        uint      `json:"user_id"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	CustomerPhone string    `json:"customer_phone"`
	Total         float64   `json:"total"`
	Status        string    `json:"status"`
	Products      []Product `gorm:"many2many:order_items;"`
}

type OrderItems struct {
	gorm.Model
	OrderID   uint `json:"order_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type Review struct {
	gorm.Model
	ProductID    uint         `json:"product_id"`
	UserID       uint         `json:"user_id"`
	Rating       int          `json:"rating"`
	Comment      string       `json:"comment"`
	CommentLikes []ReviewLike `gorm:"foreignKey:CommentID"`
}

type ReviewLike struct {
	gorm.Model
	CommentID uint `gorm:"not null;index" json:"comment_id"`
	UserID    uint `gorm:"not null;index" json:"user_id"`
}

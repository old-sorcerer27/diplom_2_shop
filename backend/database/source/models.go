package source

import (
	"database/sql/driver"
	"encoding/json"

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
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `gorm:"size:100;index" json:"category"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`

	ImageURL     string    `json:"image_url" gorm:"size:500"`
	ThumbnailURL string    `json:"thumbnail_url" gorm:"size:500"`
	MediumURL    string    `json:"medium_url" gorm:"size:500"`
	Gallery      JSONArray `json:"gallery" gorm:"type:json"`

	Orders        []Order  `gorm:"many2many:order_items;"`
	Reviews       []Review `gorm:"foreignKey:ProductID"`
	AverageRating float64  `json:"average_rating"`
	ReviewsCount  int      `json:"reviews_count"`
}

type JSONArray []string

func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return json.Marshal(j)
}

type ProductImage struct {
	ID        uint   `json:"id"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	Medium    string `json:"medium"`
	IsMain    bool   `json:"is_main"`
}

type ProductWithImages struct {
	Product
	Images        []ProductImage `json:"images"`
	MainImage     string         `json:"main_image"`
	AverageRating float64        `json:"average_rating"`
	ReviewsCount  int            `json:"reviews_count"`
}

type Order struct {
	gorm.Model
	UserID        uint      `json:"user_id"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	CustomerPhone string    `json:"customer_phone"`
	Comment       string    `json:"comment"`
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

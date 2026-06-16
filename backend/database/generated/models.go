package generated

import (
	"store/database/source"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/cli/gorm/field"
	"gorm.io/gorm"
)

var User = struct {
	ID        field.Number[uint]
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field[gorm.DeletedAt]
	Email     field.String
	Username  field.String
	Password  field.String
	Phone     field.String
	Role      field.String
	IsActive  field.Bool
	Reviews   field.Slice[source.Review]
	Orders    field.Slice[source.Order]
}{
	ID:        field.Number[uint]{}.WithColumn("id"),
	CreatedAt: field.Time{}.WithColumn("created_at"),
	UpdatedAt: field.Time{}.WithColumn("updated_at"),
	DeletedAt: field.Field[gorm.DeletedAt]{}.WithColumn("deleted_at"),
	Email:     field.String{}.WithColumn("email"),
	Username:  field.String{}.WithColumn("username"),
	Password:  field.String{}.WithColumn("password"),
	Phone:     field.String{}.WithColumn("phone"),
	Role:      field.String{}.WithColumn("role"),
	IsActive:  field.Bool{}.WithColumn("is_active"),
	Reviews:   field.Slice[source.Review]{}.WithName("Reviews"),
	Orders:    field.Slice[source.Order]{}.WithName("Orders"),
}

var Product = struct {
	ID          field.Number[uint]
	CreatedAt   field.Time
	UpdatedAt   field.Time
	DeletedAt   field.Field[gorm.DeletedAt]
	Name        field.String
	Description field.String
	Category    field.String
	Price       field.Number[float64]
	Stock       field.Number[int]
	ImageURL    field.String
	Orders      field.Slice[source.Order]
	Reviews     field.Slice[source.Review]
}{
	ID:          field.Number[uint]{}.WithColumn("id"),
	CreatedAt:   field.Time{}.WithColumn("created_at"),
	UpdatedAt:   field.Time{}.WithColumn("updated_at"),
	DeletedAt:   field.Field[gorm.DeletedAt]{}.WithColumn("deleted_at"),
	Name:        field.String{}.WithColumn("name"),
	Description: field.String{}.WithColumn("description"),
	Category:    field.String{}.WithColumn("category"),
	Price:       field.Number[float64]{}.WithColumn("price"),
	Stock:       field.Number[int]{}.WithColumn("stock"),
	ImageURL:    field.String{}.WithColumn("image_url"),
	Orders:      field.Slice[source.Order]{}.WithName("Orders"),
	Reviews:     field.Slice[source.Review]{}.WithName("Reviews"),
}

var Order = struct {
	ID            field.Number[uint]
	CreatedAt     field.Time
	UpdatedAt     field.Time
	DeletedAt     field.Field[gorm.DeletedAt]
	UserID        field.Number[uint]
	CustomerName  field.String
	CustomerEmail field.String
	CustomerPhone field.String
	Comment       field.String
	Total         field.Number[float64]
	Status        field.String
	Products      field.Slice[source.Product]
}{
	ID:            field.Number[uint]{}.WithColumn("id"),
	CreatedAt:     field.Time{}.WithColumn("created_at"),
	UpdatedAt:     field.Time{}.WithColumn("updated_at"),
	DeletedAt:     field.Field[gorm.DeletedAt]{}.WithColumn("deleted_at"),
	UserID:        field.Number[uint]{}.WithColumn("user_id"),
	CustomerName:  field.String{}.WithColumn("customer_name"),
	CustomerEmail: field.String{}.WithColumn("customer_email"),
	CustomerPhone: field.String{}.WithColumn("customer_phone"),
	Comment:       field.String{}.WithColumn("comment"),
	Total:         field.Number[float64]{}.WithColumn("total"),
	Status:        field.String{}.WithColumn("status"),
	Products:      field.Slice[source.Product]{}.WithName("Products"),
}

var OrderItems = struct {
	ID        field.Number[uint]
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field[gorm.DeletedAt]
	OrderID   field.Number[uint]
	ProductID field.Number[uint]
	Quantity  field.Number[int]
}{
	ID:        field.Number[uint]{}.WithColumn("id"),
	CreatedAt: field.Time{}.WithColumn("created_at"),
	UpdatedAt: field.Time{}.WithColumn("updated_at"),
	DeletedAt: field.Field[gorm.DeletedAt]{}.WithColumn("deleted_at"),
	OrderID:   field.Number[uint]{}.WithColumn("order_id"),
	ProductID: field.Number[uint]{}.WithColumn("product_id"),
	Quantity:  field.Number[int]{}.WithColumn("quantity"),
}

var Review = struct {
	ID        field.Number[uint]
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field[gorm.DeletedAt]
	ProductID field.Number[uint]
	UserID    field.Number[uint]
	Rating    field.Number[int]
	Comment   field.String
}{
	ID:        field.Number[uint]{}.WithColumn("id"),
	CreatedAt: field.Time{}.WithColumn("created_at"),
	UpdatedAt: field.Time{}.WithColumn("updated_at"),
	DeletedAt: field.Field[gorm.DeletedAt]{}.WithColumn("deleted_at"),
	ProductID: field.Number[uint]{}.WithColumn("product_id"),
	UserID:    field.Number[uint]{}.WithColumn("user_id"),
	Rating:    field.Number[int]{}.WithColumn("rating"),
	Comment:   field.String{}.WithColumn("comment"),
}

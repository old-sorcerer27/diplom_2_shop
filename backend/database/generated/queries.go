package generated

import (
	"context"
	"store/database/source"
	"strings"

	"gorm.io/cli/gorm/typed"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TypedQuery[T any](db *gorm.DB, opts ...clause.Expression) _TypedQueryInterface[T] {
	return _TypedQueryImpl[T]{
		Interface: typed.G[T](db, opts...),
	}
}

type _TypedQueryInterface[T any] interface {
	typed.Interface[T]
	GetByID(ctx context.Context, id int) (T, error)
	GetAll(ctx context.Context) ([]T, error)
	DeleteByID(ctx context.Context, id int) error
}

type _TypedQueryImpl[T any] struct {
	typed.Interface[T]
}

func (e _TypedQueryImpl[T]) GetByID(ctx context.Context, id int) (T, error) {
	var sb strings.Builder
	_params := make([]any, 0, 2)

	sb.WriteString("SELECT * FROM ? WHERE id = ?")
	_params = append(_params, clause.Table{Name: clause.CurrentTable}, id)

	var result T
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func (e _TypedQueryImpl[T]) GetAll(ctx context.Context) ([]T, error) {
	var sb strings.Builder
	_params := make([]any, 0, 1)

	sb.WriteString("SELECT * FROM ?")
	_params = append(_params, clause.Table{Name: clause.CurrentTable})

	var result []T
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func (e _TypedQueryImpl[T]) DeleteByID(ctx context.Context, id int) error {
	var sb strings.Builder
	_params := make([]any, 0, 2)

	sb.WriteString("DELETE FROM ? WHERE id = ?")
	_params = append(_params, clause.Table{Name: clause.CurrentTable}, id)

	return e.Exec(ctx, sb.String(), _params...)
}

func ProductQuery[T any](db *gorm.DB, opts ...clause.Expression) _ProductQueryInterface[T] {
	return _ProductQueryImpl[T]{
		Interface: typed.G[T](db, opts...),
	}
}

type _ProductQueryInterface[T any] interface {
	typed.Interface[T]
	GetByPriceRange(ctx context.Context, maxPrice float64, minStock int) ([]source.Product, error)
	SearchByName(ctx context.Context, namePattern string) ([]source.Product, error)
}

type _ProductQueryImpl[T any] struct {
	typed.Interface[T]
}

func (e _ProductQueryImpl[T]) GetByPriceRange(ctx context.Context, maxPrice float64, minStock int) ([]source.Product, error) {
	var sb strings.Builder
	_params := make([]any, 0, 2)

	sb.WriteString("SELECT * FROM products WHERE price < ? AND stock > ?")
	_params = append(_params, maxPrice, minStock)

	var result []source.Product
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func (e _ProductQueryImpl[T]) SearchByName(ctx context.Context, namePattern string) ([]source.Product, error) {
	var sb strings.Builder
	_params := make([]any, 0, 1)

	sb.WriteString("SELECT * FROM products WHERE name LIKE ?")
	_params = append(_params, namePattern)

	var result []source.Product
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func ReviewQuery[T any](db *gorm.DB, opts ...clause.Expression) _ReviewQueryInterface[T] {
	return _ReviewQueryImpl[T]{
		Interface: typed.G[T](db, opts...),
	}
}

type _ReviewQueryInterface[T any] interface {
	typed.Interface[T]
	AddReview(ctx context.Context, productID int, userID int, rating int) error
	EditReview(ctx context.Context, reviewID int, rating int, comment string) error
	GetByProductID(ctx context.Context, productID int) ([]source.Review, error)
	GetByUserID(ctx context.Context, userID int) ([]source.Review, error)
}

type _ReviewQueryImpl[T any] struct {
	typed.Interface[T]
}

func (e _ReviewQueryImpl[T]) AddReview(ctx context.Context, productID int, userID int, rating int) error {
	var sb strings.Builder
	_params := make([]any, 0, 3)

	sb.WriteString("INSERT INTO reviews (product_id, user_id, rating) VALUES (?, ?, ?)")
	_params = append(_params, productID, userID, rating)

	return e.Exec(ctx, sb.String(), _params...)
}

func (e _ReviewQueryImpl[T]) EditReview(ctx context.Context, reviewID int, rating int, comment string) error {
	var sb strings.Builder
	_params := make([]any, 0, 3)

	sb.WriteString("UPDATE reviews SET rating = ?, comment = ? WHERE id = ?")
	_params = append(_params, rating, comment, reviewID)

	return e.Exec(ctx, sb.String(), _params...)
}

func (e _ReviewQueryImpl[T]) GetByProductID(ctx context.Context, productID int) ([]source.Review, error) {
	var sb strings.Builder
	_params := make([]any, 0, 1)

	sb.WriteString("SELECT * FROM reviews WHERE product_id = ?")
	_params = append(_params, productID)

	var result []source.Review
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func (e _ReviewQueryImpl[T]) GetByUserID(ctx context.Context, userID int) ([]source.Review, error) {
	var sb strings.Builder
	_params := make([]any, 0, 1)

	sb.WriteString("SELECT * FROM reviews WHERE user_id = ?")
	_params = append(_params, userID)

	var result []source.Review
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func OrderQuery[T any](db *gorm.DB, opts ...clause.Expression) _OrderQueryInterface[T] {
	return _OrderQueryImpl[T]{
		Interface: typed.G[T](db, opts...),
	}
}

type _OrderQueryInterface[T any] interface {
	typed.Interface[T]
	GetByCustomerEmail(ctx context.Context, customerEmail string) ([]source.Order, error)
	GetByStatus(ctx context.Context, status string) ([]source.Order, error)
	GetByDateRange(ctx context.Context, startDate string, endDate string) ([]source.Order, error)
	GetOrdersProducts(ctx context.Context, orderID int) ([]source.Product, error)
	RemoveFromOrder(ctx context.Context, userID int, productID int) error
	ClearOrder(ctx context.Context, userID int) error
	CreateOrder(ctx context.Context, customerName string, customerEmail string, customerPhone string, comment string, total float64) (int, error)
	AddToOrder(ctx context.Context, orderID int, productID int, quantity int) error
	UpdateOrderItem(ctx context.Context, userID int, productID int, quantity int) error
	CalculateTotal(ctx context.Context, userID int) (float64, error)
	SendOrder(ctx context.Context, userID int) error
	GetQuantity(ctx context.Context, userID int, productID int) (int, error)
}

type _OrderQueryImpl[T any] struct {
	typed.Interface[T]
}

func (e _OrderQueryImpl[T]) GetByCustomerEmail(ctx context.Context, customerEmail string) ([]source.Order, error) {
	var sb strings.Builder
	_params := make([]any, 0, 1)

	sb.WriteString("SELECT * FROM orders WHERE customer_email = ?")
	_params = append(_params, customerEmail)

	var result []source.Order
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func (e _OrderQueryImpl[T]) GetByStatus(ctx context.Context, status string) ([]source.Order, error) {
	var sb strings.Builder
	_params := make([]any, 0, 1)

	sb.WriteString("SELECT * FROM orders WHERE status = ?")
	_params = append(_params, status)

	var result []source.Order
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func (e _OrderQueryImpl[T]) GetByDateRange(ctx context.Context, startDate string, endDate string) ([]source.Order, error) {
	var sb strings.Builder
	_params := make([]any, 0, 2)

	sb.WriteString("SELECT * FROM orders WHERE created_at >= ? AND created_at <= ?")
	_params = append(_params, startDate, endDate)

	var result []source.Order
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func (e _OrderQueryImpl[T]) GetOrdersProducts(ctx context.Context, orderID int) ([]source.Product, error) {
	var sb strings.Builder
	_params := make([]any, 0, 1)

	sb.WriteString("SELECT p.* FROM products p JOIN order_items oi ON p.id = oi.product_id WHERE oi.order_id = ?")
	_params = append(_params, orderID)

	var result []source.Product
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func (e _OrderQueryImpl[T]) RemoveFromOrder(ctx context.Context, userID int, productID int) error {
	var sb strings.Builder
	_params := make([]any, 0, 2)

	sb.WriteString("DELETE FROM order_items WHERE user_id = ? AND product_id = ?")
	_params = append(_params, userID, productID)

	return e.Exec(ctx, sb.String(), _params...)
}

func (e _OrderQueryImpl[T]) ClearOrder(ctx context.Context, userID int) error {
	var sb strings.Builder
	_params := make([]any, 0, 1)

	sb.WriteString("DELETE FROM order_items WHERE user_id = ?")
	_params = append(_params, userID)

	return e.Exec(ctx, sb.String(), _params...)
}

func (e _OrderQueryImpl[T]) CreateOrder(ctx context.Context, customerName string, customerEmail string, customerPhone string, comment string, total float64) (int, error) {
	var sb strings.Builder
	_params := make([]any, 0, 5)

	sb.WriteString("INSERT INTO orders (customer_name, customer_email, customer_phone, comment, total) VALUES (?, ?, ?, ?, ?)")
	_params = append(_params, customerName, customerEmail, customerPhone, comment, total)

	var result int
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func (e _OrderQueryImpl[T]) AddToOrder(ctx context.Context, orderID int, productID int, quantity int) error {
	var sb strings.Builder
	_params := make([]any, 0, 3)

	sb.WriteString("INSERT INTO order_items  (order_id, product_id, quantity, created_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)")
	_params = append(_params, orderID, productID, quantity)

	return e.Exec(ctx, sb.String(), _params...)
}

func (e _OrderQueryImpl[T]) UpdateOrderItem(ctx context.Context, userID int, productID int, quantity int) error {
	var sb strings.Builder
	_params := make([]any, 0, 3)

	sb.WriteString("UPDATE order_items SET quantity = ? WHERE user_id = ? AND product_id = ?")
	_params = append(_params, quantity, userID, productID)

	return e.Exec(ctx, sb.String(), _params...)
}

func (e _OrderQueryImpl[T]) CalculateTotal(ctx context.Context, userID int) (float64, error) {
	var sb strings.Builder
	_params := make([]any, 0, 1)

	sb.WriteString("SELECT SUM(p.price * oi.quantity) FROM products p JOIN order_items oi ON p.id = oi.product_id WHERE oi.user_id = ?")
	_params = append(_params, userID)

	var result float64
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

func (e _OrderQueryImpl[T]) SendOrder(ctx context.Context, userID int) error {
	var sb strings.Builder
	_params := make([]any, 0, 1)

	sb.WriteString("UPDATE orders SET status = 'sent' WHERE user_id = ?")
	_params = append(_params, userID)

	return e.Exec(ctx, sb.String(), _params...)
}

func (e _OrderQueryImpl[T]) GetQuantity(ctx context.Context, userID int, productID int) (int, error) {
	var sb strings.Builder
	_params := make([]any, 0, 2)

	sb.WriteString("SELECT quantity FROM order_items WHERE user_id = ? AND product_id = ?")
	_params = append(_params, userID, productID)

	var result int
	err := e.Raw(sb.String(), _params...).Scan(ctx, &result)
	return result, err
}

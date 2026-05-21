package source

import (
	"context"
)

type TypedQuery[T any] interface {
	// SELECT * FROM @@table WHERE id = @id
	GetByID(ctx context.Context, id int) (T, error)

	// SELECT * FROM @@table
	GetAll(ctx context.Context) ([]T, error)

	// DELETE FROM @@table WHERE id = @id
	DeleteByID(ctx context.Context, id int) error
}

type ProductQuery interface {
	// SELECT * FROM products WHERE price < @maxPrice AND stock > @minStock
	GetByPriceRange(ctx context.Context, maxPrice float64, minStock int) ([]Product, error)

	// SELECT * FROM products WHERE name LIKE @namePattern
	SearchByName(ctx context.Context, namePattern string) ([]Product, error)
}

type ReviewQuery interface {
	// INSERT INTO reviews (product_id, user_id, rating) VALUES (@productID, @userID, @rating)
	AddReview(ctx context.Context, productID int, userID int, rating int) error

	// UPDATE reviews SET rating = @rating, comment = @comment WHERE id = @reviewID
	EditReview(ctx context.Context, reviewID int, rating int, comment string) error

	// SELECT * FROM reviews WHERE product_id = @productID
	GetByProductID(ctx context.Context, productID int) ([]Review, error)

	// SELECT * FROM reviews WHERE user_id = @userID
	GetByUserID(ctx context.Context, userID int) ([]Review, error)
}

type OrderQuery interface {
	// SELECT * FROM orders WHERE customer_email = @customerEmail
	GetByCustomerEmail(ctx context.Context, customerEmail string) ([]Order, error)

	// SELECT * FROM orders WHERE status = @status
	GetByStatus(ctx context.Context, status string) ([]Order, error)

	// SELECT * FROM orders WHERE created_at >= @startDate AND created_at <= @endDate
	GetByDateRange(ctx context.Context, startDate, endDate string) ([]Order, error)

	// SELECT p.* FROM products p JOIN order_items oi ON p.id = oi.product_id WHERE oi.order_id = @orderID
	GetOrdersProducts(ctx context.Context, orderID int) ([]Product, error)

	// DELETE FROM order_items WHERE user_id = @userID AND product_id = @productID
	RemoveFromOrder(ctx context.Context, userID int, productID int) error

	// DELETE FROM order_items WHERE user_id = @userID
	ClearOrder(ctx context.Context, userID int) error

	// INSERT INTO orders (customer_name, customer_email, total) VALUES (@customerName, @customerEmail, @total)
	CreateOrder(ctx context.Context, customerName string, customerEmail string, total float64) (int, error)

	// INSERT INTO order_items  (order_id, product_id, quantity, created_at) VALUES (@orderID, @productID, @quantity, CURRENT_TIMESTAMP)
	AddToOrder(ctx context.Context, orderID int, productID int, quantity int) error

	// UPDATE order_items SET quantity = @quantity WHERE user_id = @userID AND product_id = @productID
	UpdateOrderItem(ctx context.Context, userID int, productID int, quantity int) error

	// SELECT SUM(p.price * oi.quantity) FROM products p JOIN order_items oi ON p.id = oi.product_id WHERE oi.user_id = @userID
	CalculateTotal(ctx context.Context, userID int) (float64, error)

	// UPDATE orders SET status = 'sent' WHERE user_id = @userID
	SendOrder(ctx context.Context, userID int) error

	// SELECT quantity FROM order_items WHERE user_id = @userID AND product_id = @productID
	GetQuantity(ctx context.Context, userID int, productID int) (int, error)
}

package repositories

import (
	"backend/models"
	"database/sql"
	"time"
)

type OrderRepository interface {
	Create(order models.Order) (*models.Order, error)
	GetAll() ([]models.Order, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order models.Order) (*models.Order, error) {
	// Start database transaction for atomic operation
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Will be no-op if tx.Commit() succeeds

	// Insert order into orders table
	var orderID int
	var orderedAt time.Time
	orderQuery := `
		INSERT INTO orders (user_id, total_amount, shipping_name, shipping_address, shipping_phone_number, shipping_email)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, ordered_at
	`

	// Use placeholder email if not provided
	shippingEmail := "customer@example.com"

	err = tx.QueryRow(orderQuery,
		order.UserID,
		order.TotalAmount,
		order.ShippingInfo.Name,
		order.ShippingInfo.Address,
		order.ShippingInfo.PhoneNumber,
		shippingEmail,
	).Scan(&orderID, &orderedAt)

	if err != nil {
		return nil, err
	}

	// Insert order items
	itemQuery := `
		INSERT INTO order_items (order_id, product_id, product_name, quantity, price)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, item := range order.Items {
		_, err = tx.Exec(itemQuery, orderID, item.ProductID, item.Name, item.Quantity, item.Price)
		if err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Create response order object
	resultOrder := models.Order{
		ID:           orderID,
		UserID:       order.UserID,
		Timestamp:    orderedAt.Format("2006-01-02 15:04:05"),
		OrderedAt:    orderedAt,
		Items:        order.Items,
		TotalAmount:  order.TotalAmount,
		ShippingInfo: order.ShippingInfo,
	}

	// Set order ID for each item in response
	for i := range resultOrder.Items {
		resultOrder.Items[i].OrderID = orderID
		resultOrder.Items[i].ID = i + 1 // Simple indexing for response
	}

	return &resultOrder, nil
}

func (r *orderRepository) GetAll() ([]models.Order, error) {
	// Query to get all orders with their items
	orderQuery := `
		SELECT o.id, o.user_id, o.total_amount, o.ordered_at, 
		       o.shipping_name, o.shipping_address, o.shipping_phone_number, o.shipping_email
		FROM orders o
		ORDER BY o.ordered_at DESC
	`

	rows, err := r.db.Query(orderQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	orderMap := make(map[int]*models.Order) // To efficiently group items by order

	for rows.Next() {
		var order models.Order
		var shippingEmail string

		err := rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.OrderedAt,
			&order.ShippingInfo.Name, &order.ShippingInfo.Address,
			&order.ShippingInfo.PhoneNumber, &shippingEmail)
		if err != nil {
			return nil, err
		}

		// Format timestamp for frontend compatibility
		order.Timestamp = order.OrderedAt.Format("2006-01-02 15:04:05")
		order.Items = []models.OrderItem{} // Initialize empty items slice

		orders = append(orders, order)
		orderMap[order.ID] = &orders[len(orders)-1] // Reference to the order in slice
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Now fetch all order items and group them by order
	if len(orders) > 0 {
		itemQuery := `
			SELECT oi.id, oi.order_id, oi.product_id, oi.product_name, oi.quantity, oi.price
			FROM order_items oi
			ORDER BY oi.order_id, oi.id
		`

		itemRows, err := r.db.Query(itemQuery)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item models.OrderItem

			err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID,
				&item.Name, &item.Quantity, &item.Price)
			if err != nil {
				return nil, err
			}

			// Calculate subtotal
			item.Subtotal = item.Price * float64(item.Quantity)

			// Add item to corresponding order
			if order, exists := orderMap[item.OrderID]; exists {
				order.Items = append(order.Items, item)
			}
		}

		if err = itemRows.Err(); err != nil {
			return nil, err
		}
	}

	return orders, nil
}
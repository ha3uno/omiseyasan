package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"backend/data"
	"backend/database"
	"backend/models"

	"github.com/gorilla/mux"
)

// In-memory user storage
var users = []models.User{}
var nextUserID = 1

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "hello world from Go!")
}

// GetHistoryHandler handles GET /api/history - returns all history entries sorted by timestamp descending
func GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, timestamp, description, effort_hours, claude_prompt 
		FROM update_history 
		ORDER BY timestamp DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Failed to fetch history data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []models.HistoryEntry

	for rows.Next() {
		var entry models.HistoryEntry
		var timestamp time.Time

		err := rows.Scan(&entry.ID, &timestamp, &entry.Description, &entry.EffortHours, &entry.ClaudePrompt)
		if err != nil {
			log.Printf("Database scan error: %v", err)
			http.Error(w, "Failed to parse history data", http.StatusInternalServerError)
			return
		}

		// Format timestamp for frontend compatibility
		entry.Timestamp = timestamp.Format("2006-01-02 15:04:05")
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Database rows error: %v", err)
		http.Error(w, "Failed to fetch history data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to encode history data", http.StatusInternalServerError)
		return
	}
}

// CreateHistoryHandler handles POST /api/history - creates a new history entry
func CreateHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var newEntry struct {
		Description  string  `json:"description"`
		EffortHours  float64 `json:"effortHours"`
		ClaudePrompt string  `json:"claudePrompt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newEntry); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if newEntry.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	// Insert new entry into database
	query := `
		INSERT INTO update_history (description, effort_hours, claude_prompt) 
		VALUES ($1, $2, $3) 
		RETURNING id, timestamp, description, effort_hours, claude_prompt
	`

	var entry models.HistoryEntry
	var timestamp time.Time

	err := database.DB.QueryRow(query, newEntry.Description, newEntry.EffortHours, newEntry.ClaudePrompt).Scan(
		&entry.ID, &timestamp, &entry.Description, &entry.EffortHours, &entry.ClaudePrompt,
	)

	if err != nil {
		log.Printf("Database insert error: %v", err)
		http.Error(w, "Failed to create history entry", http.StatusInternalServerError)
		return
	}

	// Format timestamp for frontend compatibility
	entry.Timestamp = timestamp.Format("2006-01-02 15:04:05")

	// Return the created entry
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetProductsHandler handles GET /api/products - returns all products with optional filtering
func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")

	// Build dynamic query
	query := "SELECT id, name, price, description, image_url, category FROM products WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if category != "" && category != "all" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR category ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	query += " ORDER BY id"

	// Execute query
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []data.Product
	for rows.Next() {
		var product data.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.ImageURL, &product.Category)
		if err != nil {
			log.Printf("Database scan error: %v", err)
			http.Error(w, "Failed to parse products", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Database rows error: %v", err)
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to encode products data", http.StatusInternalServerError)
		return
	}
}

// GetProductByIDHandler handles GET /api/products/:id - returns a specific product
func GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	// Convert string ID to int
	id := 0
	for i, char := range productID {
		if char < '0' || char > '9' {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		id = id*10 + int(char-'0')
		if i > 10 { // Prevent overflow
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
	}

	// Find product by ID from database
	query := "SELECT id, name, price, description, image_url, category FROM products WHERE id = $1"
	var product data.Product
	err := database.DB.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.ImageURL, &product.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		log.Printf("Database query error: %v", err)
		http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode product data", http.StatusInternalServerError)
		return
	}
}

// GetCategoriesHandler handles GET /api/categories - returns all unique categories
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	query := "SELECT DISTINCT category FROM products ORDER BY category"
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			log.Printf("Database scan error: %v", err)
			http.Error(w, "Failed to parse categories", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Database rows error: %v", err)
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to encode categories data", http.StatusInternalServerError)
		return
	}
}

// RegisterUserHandler handles POST /api/users/register - creates a new user account
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser struct {
		Name        string `json:"name"`
		Address     string `json:"address"`
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" {
		http.Error(w, "Name, email, and password are required", http.StatusBadRequest)
		return
	}

	// Create new user with unique ID
	user := models.User{
		ID:          nextUserID,
		Name:        newUser.Name,
		Address:     newUser.Address,
		PhoneNumber: newUser.PhoneNumber,
		Email:       newUser.Email,
		Password:    newUser.Password,
	}

	// Add to users slice and increment ID counter
	users = append(users, user)
	nextUserID++

	// Return the created user (excluding password for security)
	userResponse := struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Address     string `json:"address"`
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
	}{
		ID:          user.ID,
		Name:        user.Name,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(userResponse); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// LoginUserHandler handles POST /api/users/login - authenticates a user
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if loginData.Email == "" || loginData.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Search for user with matching email and password
	for _, user := range users {
		if user.Email == loginData.Email && user.Password == loginData.Password {
			// User found - return user data (excluding password for security)
			userResponse := struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Address     string `json:"address"`
				PhoneNumber string `json:"phoneNumber"`
				Email       string `json:"email"`
			}{
				ID:          user.ID,
				Name:        user.Name,
				Address:     user.Address,
				PhoneNumber: user.PhoneNumber,
				Email:       user.Email,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(userResponse); err != nil {
				log.Printf("JSON encoding error: %v", err)
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	// User not found or invalid credentials
	http.Error(w, "Invalid email or password", http.StatusUnauthorized)
}

// CreateOrderHandler handles POST /api/orders - creates a new order in database
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var newOrder struct {
		Items        []models.OrderItem   `json:"items"`
		TotalAmount  float64              `json:"totalAmount"`
		ShippingInfo models.ShippingInfo  `json:"shippingInfo"`
		UserID       *int                 `json:"userId,omitempty"` // Optional user ID for logged-in users
	}

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if len(newOrder.Items) == 0 {
		http.Error(w, "Order must contain at least one item", http.StatusBadRequest)
		return
	}

	if newOrder.ShippingInfo.Name == "" || newOrder.ShippingInfo.Address == "" {
		http.Error(w, "Shipping name and address are required", http.StatusBadRequest)
		return
	}

	// Calculate total amount from items to prevent manipulation
	calculatedTotal := 0.0
	for i := range newOrder.Items {
		item := &newOrder.Items[i]
		item.Subtotal = item.Price * float64(item.Quantity)
		calculatedTotal += item.Subtotal
	}

	// Start database transaction for atomic operation
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
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

	// Use empty string for shipping_email if not provided (frontend doesn't send email)
	shippingEmail := ""
	if newOrder.ShippingInfo.Address != "" {
		shippingEmail = "customer@example.com" // Placeholder - can be updated later
	}

	err = tx.QueryRow(orderQuery,
		newOrder.UserID,
		calculatedTotal,
		newOrder.ShippingInfo.Name,
		newOrder.ShippingInfo.Address,
		newOrder.ShippingInfo.PhoneNumber,
		shippingEmail,
	).Scan(&orderID, &orderedAt)

	if err != nil {
		log.Printf("Failed to insert order: %v", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Insert order items
	itemQuery := `
		INSERT INTO order_items (order_id, product_id, product_name, quantity, price)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, item := range newOrder.Items {
		_, err = tx.Exec(itemQuery, orderID, item.ProductID, item.Name, item.Quantity, item.Price)
		if err != nil {
			log.Printf("Failed to insert order item: %v", err)
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Create response order object
	order := models.Order{
		ID:           orderID,
		UserID:       newOrder.UserID,
		Timestamp:    orderedAt.Format("2006-01-02 15:04:05"),
		OrderedAt:    orderedAt,
		Items:        newOrder.Items,
		TotalAmount:  calculatedTotal,
		ShippingInfo: newOrder.ShippingInfo,
	}

	// Set order ID for each item in response
	for i := range order.Items {
		order.Items[i].OrderID = orderID
		order.Items[i].ID = i + 1 // Simple indexing for response
	}

	log.Printf("New order created in database: ID=%d, Total=%.2f, Items=%d", orderID, calculatedTotal, len(newOrder.Items))

	// Return the created order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetOrdersHandler handles GET /api/orders - returns all orders from database sorted by timestamp descending
func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	// Query to get all orders with their items
	orderQuery := `
		SELECT o.id, o.user_id, o.total_amount, o.ordered_at, 
		       o.shipping_name, o.shipping_address, o.shipping_phone_number, o.shipping_email
		FROM orders o
		ORDER BY o.ordered_at DESC
	`

	rows, err := database.DB.Query(orderQuery)
	if err != nil {
		log.Printf("Failed to query orders: %v", err)
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
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
			log.Printf("Failed to scan order: %v", err)
			http.Error(w, "Failed to parse orders", http.StatusInternalServerError)
			return
		}

		// Format timestamp for frontend compatibility
		order.Timestamp = order.OrderedAt.Format("2006-01-02 15:04:05")
		order.Items = []models.OrderItem{} // Initialize empty items slice

		orders = append(orders, order)
		orderMap[order.ID] = &orders[len(orders)-1] // Reference to the order in slice
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating orders: %v", err)
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	// Now fetch all order items and group them by order
	if len(orders) > 0 {
		itemQuery := `
			SELECT oi.id, oi.order_id, oi.product_id, oi.product_name, oi.quantity, oi.price
			FROM order_items oi
			ORDER BY oi.order_id, oi.id
		`

		itemRows, err := database.DB.Query(itemQuery)
		if err != nil {
			log.Printf("Failed to query order items: %v", err)
			http.Error(w, "Failed to fetch order items", http.StatusInternalServerError)
			return
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item models.OrderItem

			err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID,
				&item.Name, &item.Quantity, &item.Price)
			if err != nil {
				log.Printf("Failed to scan order item: %v", err)
				http.Error(w, "Failed to parse order items", http.StatusInternalServerError)
				return
			}

			// Calculate subtotal
			item.Subtotal = item.Price * float64(item.Quantity)

			// Add item to corresponding order
			if order, exists := orderMap[item.OrderID]; exists {
				order.Items = append(order.Items, item)
			}
		}

		if err = itemRows.Err(); err != nil {
			log.Printf("Error iterating order items: %v", err)
			http.Error(w, "Failed to fetch order items", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to encode orders data", http.StatusInternalServerError)
		return
	}
}
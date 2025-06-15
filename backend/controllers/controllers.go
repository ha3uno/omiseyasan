package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"backend/models"
	"backend/usecases"

	"github.com/gorilla/mux"
)

type Controllers struct {
	historyUsecase usecases.HistoryUsecase
	productUsecase usecases.ProductUsecase
	userUsecase    usecases.UserUsecase
	orderUsecase   usecases.OrderUsecase
}

func NewControllers(
	historyUsecase usecases.HistoryUsecase,
	productUsecase usecases.ProductUsecase,
	userUsecase usecases.UserUsecase,
	orderUsecase usecases.OrderUsecase,
) *Controllers {
	return &Controllers{
		historyUsecase: historyUsecase,
		productUsecase: productUsecase,
		userUsecase:    userUsecase,
		orderUsecase:   orderUsecase,
	}
}

func (c *Controllers) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "hello world from Go!")
}

// GetHistoryHandler handles GET /api/history - returns all history entries sorted by timestamp descending
func (c *Controllers) GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := c.historyUsecase.GetAllHistory()
	if err != nil {
		log.Printf("History usecase error: %v", err)
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
func (c *Controllers) CreateHistoryHandler(w http.ResponseWriter, r *http.Request) {
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

	entry, err := c.historyUsecase.CreateHistory(newEntry.Description, newEntry.EffortHours, newEntry.ClaudePrompt)
	if err != nil {
		log.Printf("History usecase error: %v", err)
		if err.Error() == "description is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to create history entry", http.StatusInternalServerError)
		}
		return
	}

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
func (c *Controllers) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")

	products, err := c.productUsecase.GetAllProducts(category, search)
	if err != nil {
		log.Printf("Product usecase error: %v", err)
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
func (c *Controllers) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	product, err := c.productUsecase.GetProductByID(productID)
	if err != nil {
		log.Printf("Product usecase error: %v", err)
		if err.Error() == "invalid product ID" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err.Error() == "product not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode product data", http.StatusInternalServerError)
		return
	}
}

// GetCategoriesHandler handles GET /api/categories - returns all unique categories
func (c *Controllers) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := c.productUsecase.GetCategories()
	if err != nil {
		log.Printf("Product usecase error: %v", err)
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
func (c *Controllers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := c.userUsecase.RegisterUser(
		newUser.Name,
		newUser.Address,
		newUser.PhoneNumber,
		newUser.Email,
		newUser.Password,
	)
	if err != nil {
		log.Printf("User usecase error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
func (c *Controllers) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	user, err := c.userUsecase.LoginUser(loginData.Email, loginData.Password)
	if err != nil {
		log.Printf("User usecase error: %v", err)
		if err.Error() == "email and password are required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}

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
}

// CreateOrderHandler handles POST /api/orders - creates a new order in database
func (c *Controllers) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var newOrder struct {
		Items        []models.OrderItem  `json:"items"`
		TotalAmount  float64             `json:"totalAmount"`
		ShippingInfo models.ShippingInfo `json:"shippingInfo"`
		UserID       *int                `json:"userId,omitempty"` // Optional user ID for logged-in users
	}

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	order, err := c.orderUsecase.CreateOrder(newOrder.Items, newOrder.ShippingInfo, newOrder.UserID)
	if err != nil {
		log.Printf("Order usecase error: %v", err)
		if err.Error() == "order must contain at least one item" || err.Error() == "shipping name and address are required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("New order created: ID=%d, Total=%.2f, Items=%d", order.ID, order.TotalAmount, len(order.Items))

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
func (c *Controllers) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	orders, err := c.orderUsecase.GetAllOrders()
	if err != nil {
		log.Printf("Order usecase error: %v", err)
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to encode orders data", http.StatusInternalServerError)
		return
	}
}
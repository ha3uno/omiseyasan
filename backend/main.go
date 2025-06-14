package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// HistoryEntry represents a single history record
type HistoryEntry struct {
	ID           int    `json:"id"`
	Timestamp    string `json:"timestamp"`
	Description  string `json:"description"`
	EffortHours  float64 `json:"effortHours"`
	ClaudePrompt string `json:"claudePrompt"`
}

// Product represents a product in the e-commerce site
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
}

// User represents a user account
type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        int     `json:"id" db:"id"`
	OrderID   int     `json:"orderId" db:"order_id"`
	ProductID int     `json:"productId" db:"product_id"`
	Name      string  `json:"name" db:"product_name"`
	Price     float64 `json:"price" db:"price"`
	Quantity  int     `json:"quantity" db:"quantity"`
	Subtotal  float64 `json:"subtotal"` // 計算値なのでdbタグなし
}

// ShippingInfo represents shipping information
type ShippingInfo struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
}

// Order represents a customer order
type Order struct {
	ID           int           `json:"id" db:"id"`
	UserID       *int          `json:"userId,omitempty" db:"user_id"` // Nullable for guest orders
	Timestamp    string        `json:"timestamp"` // フォーマット済みの文字列
	OrderedAt    time.Time     `json:"-" db:"ordered_at"` // データベース用のタイムスタンプ
	Items        []OrderItem   `json:"items"` // JOINで取得
	TotalAmount  float64       `json:"totalAmount" db:"total_amount"`
	ShippingInfo ShippingInfo  `json:"shippingInfo"` // 個別のフィールドからマッピング
}

// Database connection
var db *sql.DB

// In-memory user storage
var users = []User{}
var nextUserID = 1

// Note: Order storage is now handled by PostgreSQL database

// Sample product data - hardcoded in memory with Picsum Photos images
var products = []Product{
	{
		ID:          1,
		Name:        "モダンアートフレーム",
		Price:       3980,
		Description: "お部屋を彩るスタイリッシュなアートフレーム。モダンなデザインでどんなインテリアにも合います。",
		ImageURL:    "https://picsum.photos/300/300?random=1",
	},
	{
		ID:          2,
		Name:        "ナチュラルウッドテーブル",
		Price:       12800,
		Description: "天然木を使用したシンプルで美しいテーブル。温かみのある木目が特徴的です。",
		ImageURL:    "https://picsum.photos/300/300?random=2",
	},
	{
		ID:          3,
		Name:        "ヴィンテージレザーバッグ",
		Price:       8900,
		Description: "上質なレザーを使用したヴィンテージ風バッグ。使うほどに味が出てきます。",
		ImageURL:    "https://picsum.photos/300/300?random=3",
	},
	{
		ID:          4,
		Name:        "セラミック花瓶",
		Price:       2200,
		Description: "手作りの温かみを感じるセラミック花瓶。お花を生けて空間を華やかに演出できます。",
		ImageURL:    "https://picsum.photos/300/300?random=4",
	},
	{
		ID:          5,
		Name:        "オーガニックコットンクッション",
		Price:       1680,
		Description: "肌に優しいオーガニックコットン100%のクッション。ナチュラルな風合いが魅力です。",
		ImageURL:    "https://picsum.photos/300/300?random=5",
	},
	{
		ID:          6,
		Name:        "アロマキャンドルセット",
		Price:       2450,
		Description: "リラックス効果の高い天然アロマを使用したキャンドルセット。癒しの時間をお届けします。",
		ImageURL:    "https://picsum.photos/300/300?random=6",
	},
	{
		ID:          7,
		Name:        "ハンドメイドソープコレクション",
		Price:       1890,
		Description: "天然成分にこだわったハンドメイドソープのコレクション。お肌に優しく香りも豊かです。",
		ImageURL:    "https://picsum.photos/300/300?random=7",
	},
	{
		ID:          8,
		Name:        "竹製キッチンツールセット",
		Price:       3200,
		Description: "環境に優しい竹素材のキッチンツールセット。軽くて丈夫、長くご愛用いただけます。",
		ImageURL:    "https://picsum.photos/300/300?random=8",
	},
	{
		ID:          9,
		Name:        "ガラス製デキャンタ",
		Price:       4650,
		Description: "エレガントなフォルムのガラス製デキャンタ。ワインやウイスキーをより美味しく楽しめます。",
		ImageURL:    "https://picsum.photos/300/300?random=9",
	},
	{
		ID:          10,
		Name:        "フェルト製収納ボックス",
		Price:       1280,
		Description: "柔らかなフェルト素材の収納ボックス。シンプルなデザインで様々な用途に使えます。",
		ImageURL:    "https://picsum.photos/300/300?random=10",
	},
}

func main() {
	// Initialize database connection
	initDB()
	defer db.Close()

	// Run database migrations
	runMigrations()

	// Create tables if they don't exist (fallback for existing code)
	createTables()

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/hello", helloHandler).Methods("GET")
	api.HandleFunc("/history", getHistoryHandler).Methods("GET")
	api.HandleFunc("/history", createHistoryHandler).Methods("POST")
	api.HandleFunc("/products", getProductsHandler).Methods("GET")
	api.HandleFunc("/products/{id}", getProductByIDHandler).Methods("GET")
	api.HandleFunc("/users/register", registerUserHandler).Methods("POST")
	api.HandleFunc("/users/login", loginUserHandler).Methods("POST")
	api.HandleFunc("/orders", createOrderHandler).Methods("POST")
	api.HandleFunc("/orders", getOrdersHandler).Methods("GET")

	// Serve static files from frontend/build directory
	staticDir := "../frontend/build"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("Warning: Static directory %s does not exist. Run 'cd ../frontend && npm run build' first.", staticDir)
	}

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(staticDir, "static")))))

	// Serve other assets
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "favicon.ico"))
	})
	r.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "manifest.json"))
	})

	// Catch-all handler for SPA routing - serve index.html for all non-API routes
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't serve index.html for API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		indexPath := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			http.Error(w, "index.html not found. Please build the frontend first.", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, indexPath)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("API endpoints:\n")
	fmt.Printf("  - http://localhost:%s/api/hello\n", port)
	fmt.Printf("  - http://localhost:%s/api/products\n", port)
	fmt.Printf("  - http://localhost:%s/api/history\n", port)
	fmt.Printf("Frontend: http://localhost:%s/\n", port)
	
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// getProductByIDHandler handles GET /api/products/:id - returns a specific product
func getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
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
	
	// Find product by ID
	for _, product := range products {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(product); err != nil {
				http.Error(w, "Failed to encode product data", http.StatusInternalServerError)
				return
			}
			return
		}
	}
	
	// Product not found
	http.Error(w, "Product not found", http.StatusNotFound)
}

// registerUserHandler handles POST /api/users/register - creates a new user account
func registerUserHandler(w http.ResponseWriter, r *http.Request) {
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
	user := User{
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

// loginUserHandler handles POST /api/users/login - authenticates a user
func loginUserHandler(w http.ResponseWriter, r *http.Request) {
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

// initDB initializes the database connection
func initDB() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
}

// runMigrations runs database migrations using golang-migrate
func runMigrations() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required for migrations")
	}

	// main.go自身のファイルパスを取得（デプロイ環境で確実にパスを特定）
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Failed to get caller information for main.go")
	}
	
	// main.goのディレクトリ (backend/)
	currentFileDir := filepath.Dir(filename)
	// プロジェクトルート (backend/の親ディレクトリ)
	projectRoot := filepath.Dir(currentFileDir)
	// マイグレーションディレクトリの絶対パス
	migrationsPath := "file://" + filepath.Join(projectRoot, "migrations")
	
	log.Printf("main.go file path: %s", filename)
	log.Printf("backend directory: %s", currentFileDir)
	log.Printf("Project root: %s", projectRoot)
	log.Printf("Migrations path: %s", migrationsPath)

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create postgres driver for migrations: %v", err)
	}

	// Create migrate instance with file source using absolute path
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath, // Source: absolute path to migrations directory
		"postgres",     // Database name
		driver,         // Database driver instance
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance with path '%s': %v", migrationsPath, err)
	}

	// Run up migrations
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Database schema is up to date")
		} else {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	} else {
		log.Println("Database migrations applied successfully")
	}
}

// createTables creates the update_history table if it doesn't exist (fallback)
func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS update_history (
		id SERIAL PRIMARY KEY,
		timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		description TEXT NOT NULL,
		effort_hours NUMERIC(5, 2) NOT NULL,
		claude_prompt TEXT NOT NULL
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatal("Failed to create update_history table:", err)
	}

	log.Println("Database tables initialized successfully")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "hello world from Go!")
}

// getHistoryHandler handles GET /api/history - returns all history entries sorted by timestamp descending
func getHistoryHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, timestamp, description, effort_hours, claude_prompt 
		FROM update_history 
		ORDER BY timestamp DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Failed to fetch history data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []HistoryEntry

	for rows.Next() {
		var entry HistoryEntry
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

// createHistoryHandler handles POST /api/history - creates a new history entry
func createHistoryHandler(w http.ResponseWriter, r *http.Request) {
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

	var entry HistoryEntry
	var timestamp time.Time

	err := db.QueryRow(query, newEntry.Description, newEntry.EffortHours, newEntry.ClaudePrompt).Scan(
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

// getProductsHandler handles GET /api/products - returns all products
func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode products data", http.StatusInternalServerError)
		return
	}
}

// createOrderHandler handles POST /api/orders - creates a new order in database
func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var newOrder struct {
		Items        []OrderItem   `json:"items"`
		TotalAmount  float64       `json:"totalAmount"`
		ShippingInfo ShippingInfo  `json:"shippingInfo"`
		UserID       *int          `json:"userId,omitempty"` // Optional user ID for logged-in users
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
	tx, err := db.Begin()
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
	order := Order{
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

// getOrdersHandler handles GET /api/orders - returns all orders from database sorted by timestamp descending
func getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	// Query to get all orders with their items
	orderQuery := `
		SELECT o.id, o.user_id, o.total_amount, o.ordered_at, 
		       o.shipping_name, o.shipping_address, o.shipping_phone_number, o.shipping_email
		FROM orders o
		ORDER BY o.ordered_at DESC
	`
	
	rows, err := db.Query(orderQuery)
	if err != nil {
		log.Printf("Failed to query orders: %v", err)
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	var orders []Order
	orderMap := make(map[int]*Order) // To efficiently group items by order
	
	for rows.Next() {
		var order Order
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
		order.Items = []OrderItem{} // Initialize empty items slice
		
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
		
		itemRows, err := db.Query(itemQuery)
		if err != nil {
			log.Printf("Failed to query order items: %v", err)
			http.Error(w, "Failed to fetch order items", http.StatusInternalServerError)
			return
		}
		defer itemRows.Close()
		
		for itemRows.Next() {
			var item OrderItem
			
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
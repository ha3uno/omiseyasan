package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// TestMain runs setup and teardown for all tests
func TestMain(m *testing.M) {
	// Set up test database
	setupTestDB()
	
	// Run all tests
	code := m.Run()
	
	// Clean up test database
	teardownTestDB()
	
	os.Exit(code)
}

// setupTestDB starts PostgreSQL container and runs migrations
func setupTestDB() {
	log.Println("Setting up test database...")
	
	// Start PostgreSQL container using docker compose (space-separated for Render compatibility)
	cmd := exec.Command("docker", "compose", "-f", "../docker-compose.test.yml", "up", "-d")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to start test database: %v", err)
	}
	
	// Wait for PostgreSQL to be ready
	log.Println("Waiting for PostgreSQL to be ready...")
	for i := 0; i < 30; i++ {
		if checkDBConnection() {
			break
		}
		time.Sleep(time.Second)
		if i == 29 {
			log.Fatal("PostgreSQL did not start within 30 seconds")
		}
	}
	
	// Set environment variable for test database
	os.Setenv("DATABASE_URL", "postgres://test_user:test_password@localhost/omiseyasan_test?sslmode=disable")
	
	// Initialize database connection
	initDB()
	
	// Run migrations
	runTestMigrations()
	
	log.Println("Test database setup completed")
}

// teardownTestDB stops PostgreSQL container
func teardownTestDB() {
	log.Println("Tearing down test database...")
	
	// Close database connection
	if db != nil {
		db.Close()
	}
	
	// Stop PostgreSQL container using docker compose (space-separated for Render compatibility)
	cmd := exec.Command("docker", "compose", "-f", "../docker-compose.test.yml", "down")
	if err := cmd.Run(); err != nil {
		log.Printf("Warning: Failed to stop test database: %v", err)
	}
	
	log.Println("Test database teardown completed")
}

// checkDBConnection checks if PostgreSQL is ready to accept connections
func checkDBConnection() bool {
	testDB, err := sql.Open("postgres", "postgres://test_user:test_password@localhost/omiseyasan_test?sslmode=disable")
	if err != nil {
		return false
	}
	defer testDB.Close()
	
	return testDB.Ping() == nil
}

// runTestMigrations runs database migrations for test database
func runTestMigrations() {
	log.Println("Running test migrations...")
	
	// Get migrations path using runtime
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get caller information")
	}
	
	// Get project root directory
	currentFileDir := filepath.Dir(filename)
	projectRoot := filepath.Dir(currentFileDir)
	migrationsPath := "file://" + filepath.Join(projectRoot, "migrations")
	
	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create postgres driver for test migrations: %v", err)
	}
	
	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance for tests: %v", err)
	}
	
	// Run up migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run test migrations: %v", err)
	}
	
	log.Println("Test migrations completed")
}

// clearTestData clears all test data from database
func clearTestData() {
	_, err := db.Exec("DELETE FROM order_items")
	if err != nil {
		log.Printf("Warning: Failed to clear order_items: %v", err)
	}
	
	_, err = db.Exec("DELETE FROM orders")
	if err != nil {
		log.Printf("Warning: Failed to clear orders: %v", err)
	}
	
	_, err = db.Exec("DELETE FROM update_history")
	if err != nil {
		log.Printf("Warning: Failed to clear update_history: %v", err)
	}
	
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		log.Printf("Warning: Failed to clear users: %v", err)
	}
}

// TestAPIs runs integration tests for all API endpoints
func TestAPIs(t *testing.T) {
	// Clear test data before each test
	clearTestData()
	
	// Create router
	r := mux.NewRouter()
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
	
	// Test hello endpoint
	t.Run("GET /api/hello", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/hello", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		
		expected := "hello world from Go!"
		if w.Body.String() != expected {
			t.Errorf("Expected body '%s', got '%s'", expected, w.Body.String())
		}
	})
	
	// Test products endpoint
	t.Run("GET /api/products", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/products", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		
		var products []Product
		if err := json.NewDecoder(w.Body).Decode(&products); err != nil {
			t.Errorf("Failed to decode products: %v", err)
		}
		
		if len(products) == 0 {
			t.Error("Expected products to be returned")
		}
	})
	
	// Test product by ID endpoint
	t.Run("GET /api/products/1", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/products/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		
		var product Product
		if err := json.NewDecoder(w.Body).Decode(&product); err != nil {
			t.Errorf("Failed to decode product: %v", err)
		}
		
		if product.ID != 1 {
			t.Errorf("Expected product ID 1, got %d", product.ID)
		}
	})
	
	// Test history endpoints
	t.Run("POST /api/history", func(t *testing.T) {
		historyData := map[string]interface{}{
			"description":  "Test history entry",
			"effortHours":  2.5,
			"claudePrompt": "Test prompt",
		}
		
		jsonData, _ := json.Marshal(historyData)
		req := httptest.NewRequest("POST", "/api/history", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}
		
		var entry HistoryEntry
		if err := json.NewDecoder(w.Body).Decode(&entry); err != nil {
			t.Errorf("Failed to decode history entry: %v", err)
		}
		
		if entry.Description != "Test history entry" {
			t.Errorf("Expected description 'Test history entry', got '%s'", entry.Description)
		}
	})
	
	t.Run("GET /api/history", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/history", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		
		var entries []HistoryEntry
		if err := json.NewDecoder(w.Body).Decode(&entries); err != nil {
			t.Errorf("Failed to decode history entries: %v", err)
		}
		
		if len(entries) == 0 {
			t.Error("Expected at least one history entry")
		}
	})
	
	// Test user registration and login
	t.Run("POST /api/users/register", func(t *testing.T) {
		userData := map[string]interface{}{
			"name":        "Test User",
			"address":     "123 Test St",
			"phoneNumber": "123-456-7890",
			"email":       "test@example.com",
			"password":    "testpass123",
		}
		
		jsonData, _ := json.Marshal(userData)
		req := httptest.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}
		
		var user User
		if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
			t.Errorf("Failed to decode user: %v", err)
		}
		
		if user.Name != "Test User" {
			t.Errorf("Expected name 'Test User', got '%s'", user.Name)
		}
		
		if user.Email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
		}
	})
	
	t.Run("POST /api/users/login", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    "test@example.com",
			"password": "testpass123",
		}
		
		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		
		var user User
		if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
			t.Errorf("Failed to decode user: %v", err)
		}
		
		if user.Email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
		}
	})
	
	// Test order creation and retrieval
	t.Run("POST /api/orders", func(t *testing.T) {
		orderData := map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"productId": 1,
					"name":      "Test Product",
					"price":     29.99,
					"quantity":  2,
				},
			},
			"totalAmount": 59.98,
			"shippingInfo": map[string]interface{}{
				"name":        "John Doe",
				"address":     "456 Order St",
				"phoneNumber": "987-654-3210",
			},
		}
		
		jsonData, _ := json.Marshal(orderData)
		req := httptest.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}
		
		var order Order
		if err := json.NewDecoder(w.Body).Decode(&order); err != nil {
			t.Errorf("Failed to decode order: %v", err)
		}
		
		if order.TotalAmount != 59.98 {
			t.Errorf("Expected total amount 59.98, got %f", order.TotalAmount)
		}
		
		if len(order.Items) != 1 {
			t.Errorf("Expected 1 item, got %d", len(order.Items))
		}
		
		if order.ShippingInfo.Name != "John Doe" {
			t.Errorf("Expected shipping name 'John Doe', got '%s'", order.ShippingInfo.Name)
		}
		
		// Verify order was saved to database
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
		if err != nil {
			t.Errorf("Failed to count orders in database: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 order in database, got %d", count)
		}
	})
	
	t.Run("GET /api/orders", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/orders", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		
		var orders []Order
		if err := json.NewDecoder(w.Body).Decode(&orders); err != nil {
			t.Errorf("Failed to decode orders: %v", err)
		}
		
		if len(orders) == 0 {
			t.Error("Expected at least one order")
		}
		
		if len(orders) > 0 {
			order := orders[0]
			if order.TotalAmount != 59.98 {
				t.Errorf("Expected total amount 59.98, got %f", order.TotalAmount)
			}
			
			if len(order.Items) != 1 {
				t.Errorf("Expected 1 item, got %d", len(order.Items))
			}
		}
	})
}

// TestDatabaseConnection tests the database connection
func TestDatabaseConnection(t *testing.T) {
	if db == nil {
		t.Fatal("Database connection is nil")
	}
	
	if err := db.Ping(); err != nil {
		t.Errorf("Failed to ping database: %v", err)
	}
}

// TestMigrations tests that all required tables exist
func TestMigrations(t *testing.T) {
	tables := []string{"update_history", "users", "orders", "order_items"}
	
	for _, table := range tables {
		var exists bool
		query := fmt.Sprintf(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = '%s'
			)`, table)
		
		err := db.QueryRow(query).Scan(&exists)
		if err != nil {
			t.Errorf("Failed to check if table %s exists: %v", table, err)
		}
		
		if !exists {
			t.Errorf("Table %s does not exist", table)
		}
	}
}
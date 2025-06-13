package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// Database connection
var db *sql.DB

// Sample product data - hardcoded in memory
var products = []Product{
	{
		ID:          1,
		Name:        "かわいいぬいぐるみ",
		Price:       2980,
		Description: "ふわふわで抱き心地抜群のぬいぐるみです。お子様へのプレゼントにも最適！",
		ImageURL:    "https://via.placeholder.com/300x300/FFB6C1/FFFFFF?text=ぬいぐるみ",
	},
	{
		ID:          2,
		Name:        "おしゃれなマグカップ",
		Price:       1280,
		Description: "毎日のコーヒータイムを特別にしてくれる、シンプルで上品なデザインのマグカップ。",
		ImageURL:    "https://via.placeholder.com/300x300/87CEEB/FFFFFF?text=マグカップ",
	},
	{
		ID:          3,
		Name:        "手作りキャンドル",
		Price:       1800,
		Description: "自然な香りでリラックス効果抜群。お部屋の雰囲気を優しく演出します。",
		ImageURL:    "https://via.placeholder.com/300x300/DDA0DD/FFFFFF?text=キャンドル",
	},
}

func main() {
	// Initialize database connection
	initDB()
	defer db.Close()

	// Create tables if they don't exist
	createTables()

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/hello", helloHandler).Methods("GET")
	api.HandleFunc("/history", getHistoryHandler).Methods("GET")
	api.HandleFunc("/history", createHistoryHandler).Methods("POST")
	api.HandleFunc("/products", getProductsHandler).Methods("GET")

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

// createTables creates the update_history table if it doesn't exist
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
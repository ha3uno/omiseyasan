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

// Sample product data - hardcoded in memory with Irasutoya images
var products = []Product{
	{
		ID:          1,
		Name:        "かわいいぬいぐるみ",
		Price:       2980,
		Description: "ふわふわで抱き心地抜群のぬいぐるみです。お子様へのプレゼントにも最適！",
		ImageURL:    "https://4.bp.blogspot.com/-Kn7_AO8_50M/VrLCy-NK8GI/AAAAAAAA34s/xhYVeQZdYj8/s400/nuigurumi_kuma.png",
	},
	{
		ID:          2,
		Name:        "おしゃれなマグカップ",
		Price:       1280,
		Description: "毎日のコーヒータイムを特別にしてくれる、シンプルで上品なデザインのマグカップ。",
		ImageURL:    "https://1.bp.blogspot.com/-SZLvBxvLn8k/VXgBfz6JjEI/AAAAAAAAuPE/w1Kn6lXPdEY/s400/mug_cup.png",
	},
	{
		ID:          3,
		Name:        "手作りキャンドル",
		Price:       1800,
		Description: "自然な香りでリラックス効果抜群。お部屋の雰囲気を優しく演出します。",
		ImageURL:    "https://1.bp.blogspot.com/-VgsGrrmGGT0/VeN_IqmzUKI/AAAAAAAAwcg/mBFAZHwX0mI/s400/candle.png",
	},
	{
		ID:          4,
		Name:        "北欧風クッション",
		Price:       3200,
		Description: "お部屋をおしゃれに彩る北欧デザインのクッション。上質な素材で長くご愛用いただけます。",
		ImageURL:    "https://1.bp.blogspot.com/-h-IeZXJGNNs/VeN_H9AH4BI/AAAAAAAAwcQ/2MrCNVKqHZc/s400/cushion.png",
	},
	{
		ID:          5,
		Name:        "ハンドメイド石鹸",
		Price:       680,
		Description: "天然成分100%の優しいハンドメイド石鹸。お肌に優しく、良い香りが特徴です。",
		ImageURL:    "https://2.bp.blogspot.com/-K1IfHIhZAO4/WASJbGDkWvI/AAAAAAAA_zY/G5kXOC8bhNs4eLHb9KeKkfIcFGEw9gY0ACLcB/s400/soap.png",
	},
	{
		ID:          6,
		Name:        "オーガニック紅茶",
		Price:       1450,
		Description: "厳選されたオーガニック茶葉を使用した上品な紅茶。リラックスタイムにぴったりです。",
		ImageURL:    "https://1.bp.blogspot.com/-VaGUUVgPQKI/VFm9nBRtHzI/AAAAAAAApG0/Hzi3lj_7_nE/s400/tea_cup.png",
	},
	{
		ID:          7,
		Name:        "アロマディフューザー",
		Price:       4800,
		Description: "お部屋を良い香りで満たすスタイリッシュなアロマディフューザー。LEDライト付きです。",
		ImageURL:    "https://1.bp.blogspot.com/-aQwGaQnQqJc/WP_WaUQzqWI/AAAAAAABEWo/HdqmQ_mzKJoB6GfLR7c9Sk9kefaJq2FkACLcB/s400/aroma_oil.png",
	},
	{
		ID:          8,
		Name:        "天然木コースター",
		Price:       890,
		Description: "職人が一つ一つ手作りした天然木のコースターセット。温かみのある仕上がりです。",
		ImageURL:    "https://1.bp.blogspot.com/-VJ1L6-Q6z3s/WeL1qPCmZTI/AAAAAAABGsI/C0bTa_4a4x4KU8-i_QE67jHx0MDKnhbqgCLcBGAs/s400/mokuzai_plate.png",
	},
	{
		ID:          9,
		Name:        "ミニ観葉植物",
		Price:       1650,
		Description: "お手入れ簡単で初心者にもおすすめの可愛いミニ観葉植物。デスクにぴったりです。",
		ImageURL:    "https://2.bp.blogspot.com/-ALqjO2pUJ3E/VrYJQ8H4EAI/AAAAAAAA384/5x_3VQ8YlAo/s400/kanyoushokubutsu_pot.png",
	},
	{
		ID:          10,
		Name:        "オーガニックハンドクリーム",
		Price:       1280,
		Description: "乾燥からお肌を守る天然成分のハンドクリーム。さらりとした使い心地が人気です。",
		ImageURL:    "https://1.bp.blogspot.com/-vG-lxj8d-K8/VFm9s1ZSa8I/AAAAAAAApHI/GwS3XHQZ1TE/s400/hand_cream.png",
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
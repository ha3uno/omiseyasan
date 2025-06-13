package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
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

// In-memory storage for history entries
var (
	historyEntries []HistoryEntry
	historyMutex   sync.RWMutex
	nextID         int = 1
)

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

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "hello world from Go!")
}

// getHistoryHandler handles GET /api/history - returns all history entries sorted by timestamp descending
func getHistoryHandler(w http.ResponseWriter, r *http.Request) {
	historyMutex.RLock()
	defer historyMutex.RUnlock()

	// Create a copy of the slice to avoid race conditions
	entriesCopy := make([]HistoryEntry, len(historyEntries))
	copy(entriesCopy, historyEntries)

	// Sort by timestamp descending (newest first)
	sort.Slice(entriesCopy, func(i, j int) bool {
		return entriesCopy[i].Timestamp > entriesCopy[j].Timestamp
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entriesCopy); err != nil {
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
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if newEntry.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	historyMutex.Lock()
	defer historyMutex.Unlock()

	// Create new history entry with generated ID and timestamp
	entry := HistoryEntry{
		ID:           nextID,
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		Description:  newEntry.Description,
		EffortHours:  newEntry.EffortHours,
		ClaudePrompt: newEntry.ClaudePrompt,
	}
	nextID++

	// Add to history
	historyEntries = append(historyEntries, entry)

	// Return the created entry
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(entry); err != nil {
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
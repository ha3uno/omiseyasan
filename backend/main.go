package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"backend/controllers"
	"backend/database"
	"backend/repositories"
	"backend/usecases"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connection
	database.InitDB()
	defer database.DB.Close()

	// Run database migrations
	database.RunMigrations()

	// Create tables if they don't exist (fallback for existing code)
	database.CreateTables()

	// Initialize repositories
	historyRepo := repositories.NewHistoryRepository(database.DB)
	productRepo := repositories.NewProductRepository(database.DB)
	userRepo := repositories.NewUserRepository()
	orderRepo := repositories.NewOrderRepository(database.DB)

	// Initialize usecases
	historyUsecase := usecases.NewHistoryUsecase(historyRepo)
	productUsecase := usecases.NewProductUsecase(productRepo)
	userUsecase := usecases.NewUserUsecase(userRepo)
	orderUsecase := usecases.NewOrderUsecase(orderRepo)

	// Initialize controllers
	controllers := controllers.NewControllers(
		historyUsecase,
		productUsecase,
		userUsecase,
		orderUsecase,
	)

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/hello", controllers.HelloHandler).Methods("GET")
	api.HandleFunc("/history", controllers.GetHistoryHandler).Methods("GET")
	api.HandleFunc("/history", controllers.CreateHistoryHandler).Methods("POST")
	api.HandleFunc("/products", controllers.GetProductsHandler).Methods("GET")
	api.HandleFunc("/products/{id}", controllers.GetProductByIDHandler).Methods("GET")
	api.HandleFunc("/categories", controllers.GetCategoriesHandler).Methods("GET")
	api.HandleFunc("/users/register", controllers.RegisterUserHandler).Methods("POST")
	api.HandleFunc("/users/login", controllers.LoginUserHandler).Methods("POST")
	api.HandleFunc("/orders", controllers.CreateOrderHandler).Methods("POST")
	api.HandleFunc("/orders", controllers.GetOrdersHandler).Methods("GET")

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
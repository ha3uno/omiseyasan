package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// DB holds the database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
}

// RunMigrations runs database migrations using golang-migrate
func RunMigrations() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required for migrations")
	}

	// main.go自身のファイルパスを取得（デプロイ環境で確実にパスを特定）
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Failed to get caller information for database.go")
	}

	// database.goのディレクトリ (backend/database/)
	currentFileDir := filepath.Dir(filename)
	// backendディレクトリ (database/の親ディレクトリ)
	backendDir := filepath.Dir(currentFileDir)
	// プロジェクトルート (backend/の親ディレクトリ)
	projectRoot := filepath.Dir(backendDir)
	// マイグレーションディレクトリの絶対パス
	migrationsPath := "file://" + filepath.Join(projectRoot, "migrations")

	log.Printf("database.go file path: %s", filename)
	log.Printf("backend directory: %s", backendDir)
	log.Printf("Project root: %s", projectRoot)
	log.Printf("Migrations path: %s", migrationsPath)

	// Create postgres driver instance
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
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

// CreateTables creates the update_history table if it doesn't exist (fallback)
func CreateTables() {
	query := `
	CREATE TABLE IF NOT EXISTS update_history (
		id SERIAL PRIMARY KEY,
		timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		description TEXT NOT NULL,
		effort_hours NUMERIC(5, 2) NOT NULL,
		claude_prompt TEXT NOT NULL
	);`

	if _, err := DB.Exec(query); err != nil {
		log.Fatal("Failed to create update_history table:", err)
	}

	log.Println("Database tables initialized successfully")
}
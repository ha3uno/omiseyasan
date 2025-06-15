package repositories

import (
	"backend/data"
	"database/sql"
	"fmt"
)

type ProductRepository interface {
	GetAll(category, search string) ([]data.Product, error)
	GetByID(id int) (*data.Product, error)
	GetCategories() ([]string, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll(category, search string) ([]data.Product, error) {
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
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []data.Product
	for rows.Next() {
		var product data.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.ImageURL, &product.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) GetByID(id int) (*data.Product, error) {
	query := "SELECT id, name, price, description, image_url, category FROM products WHERE id = $1"
	var product data.Product
	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.ImageURL, &product.Category)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetCategories() ([]string, error) {
	query := "SELECT DISTINCT category FROM products ORDER BY category"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
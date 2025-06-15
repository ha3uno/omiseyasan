package usecases

import (
	"backend/data"
	"backend/repositories"
	"database/sql"
	"errors"
	"strconv"
)

type ProductUsecase interface {
	GetAllProducts(category, search string) ([]data.Product, error)
	GetProductByID(productID string) (*data.Product, error)
	GetCategories() ([]string, error)
}

type productUsecase struct {
	productRepo repositories.ProductRepository
}

func NewProductUsecase(productRepo repositories.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (u *productUsecase) GetAllProducts(category, search string) ([]data.Product, error) {
	return u.productRepo.GetAll(category, search)
}

func (u *productUsecase) GetProductByID(productID string) (*data.Product, error) {
	// Convert string ID to int with validation
	id, err := strconv.Atoi(productID)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	if id <= 0 {
		return nil, errors.New("invalid product ID")
	}

	product, err := u.productRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return product, nil
}

func (u *productUsecase) GetCategories() ([]string, error) {
	return u.productRepo.GetCategories()
}
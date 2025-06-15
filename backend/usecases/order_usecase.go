package usecases

import (
	"backend/models"
	"backend/repositories"
	"errors"
)

type OrderUsecase interface {
	CreateOrder(items []models.OrderItem, shippingInfo models.ShippingInfo, userID *int) (*models.Order, error)
	GetAllOrders() ([]models.Order, error)
}

type orderUsecase struct {
	orderRepo repositories.OrderRepository
}

func NewOrderUsecase(orderRepo repositories.OrderRepository) OrderUsecase {
	return &orderUsecase{
		orderRepo: orderRepo,
	}
}

func (u *orderUsecase) CreateOrder(items []models.OrderItem, shippingInfo models.ShippingInfo, userID *int) (*models.Order, error) {
	// Validate required fields
	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	if shippingInfo.Name == "" || shippingInfo.Address == "" {
		return nil, errors.New("shipping name and address are required")
	}

	// Calculate total amount from items to prevent manipulation
	calculatedTotal := 0.0
	for i := range items {
		item := &items[i]
		item.Subtotal = item.Price * float64(item.Quantity)
		calculatedTotal += item.Subtotal
	}

	order := models.Order{
		UserID:       userID,
		Items:        items,
		TotalAmount:  calculatedTotal,
		ShippingInfo: shippingInfo,
	}

	return u.orderRepo.Create(order)
}

func (u *orderUsecase) GetAllOrders() ([]models.Order, error) {
	return u.orderRepo.GetAll()
}
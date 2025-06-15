package usecases

import (
	"backend/models"
	"backend/repositories"
	"errors"
)

type UserUsecase interface {
	RegisterUser(name, address, phoneNumber, email, password string) (*models.User, error)
	LoginUser(email, password string) (*models.User, error)
}

type userUsecase struct {
	userRepo repositories.UserRepository
}

func NewUserUsecase(userRepo repositories.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) RegisterUser(name, address, phoneNumber, email, password string) (*models.User, error) {
	// Validate required fields
	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email, and password are required")
	}

	// Create new user
	user := models.User{
		Name:        name,
		Address:     address,
		PhoneNumber: phoneNumber,
		Email:       email,
		Password:    password,
	}

	createdUser := u.userRepo.Create(user)
	return &createdUser, nil
}

func (u *userUsecase) LoginUser(email, password string) (*models.User, error) {
	// Validate required fields
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// Search for user with matching email and password
	user := u.userRepo.FindByEmailAndPassword(email, password)
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
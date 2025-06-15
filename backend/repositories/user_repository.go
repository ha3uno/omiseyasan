package repositories

import (
	"backend/models"
)

type UserRepository interface {
	GetAll() []models.User
	Create(user models.User) models.User
	FindByEmailAndPassword(email, password string) *models.User
	GetNextID() int
}

type userRepository struct {
	users      []models.User
	nextUserID int
}

func NewUserRepository() UserRepository {
	return &userRepository{
		users:      []models.User{},
		nextUserID: 1,
	}
}

func (r *userRepository) GetAll() []models.User {
	return r.users
}

func (r *userRepository) Create(user models.User) models.User {
	user.ID = r.nextUserID
	r.users = append(r.users, user)
	r.nextUserID++
	return user
}

func (r *userRepository) FindByEmailAndPassword(email, password string) *models.User {
	for _, user := range r.users {
		if user.Email == email && user.Password == password {
			return &user
		}
	}
	return nil
}

func (r *userRepository) GetNextID() int {
	return r.nextUserID
}
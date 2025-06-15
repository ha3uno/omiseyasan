package usecases

import (
	"backend/models"
	"backend/repositories"
	"errors"
)

type HistoryUsecase interface {
	GetAllHistory() ([]models.HistoryEntry, error)
	CreateHistory(description string, effortHours float64, claudePrompt string) (*models.HistoryEntry, error)
}

type historyUsecase struct {
	historyRepo repositories.HistoryRepository
}

func NewHistoryUsecase(historyRepo repositories.HistoryRepository) HistoryUsecase {
	return &historyUsecase{
		historyRepo: historyRepo,
	}
}

func (u *historyUsecase) GetAllHistory() ([]models.HistoryEntry, error) {
	return u.historyRepo.GetAll()
}

func (u *historyUsecase) CreateHistory(description string, effortHours float64, claudePrompt string) (*models.HistoryEntry, error) {
	// Validate required fields
	if description == "" {
		return nil, errors.New("description is required")
	}

	return u.historyRepo.Create(description, effortHours, claudePrompt)
}
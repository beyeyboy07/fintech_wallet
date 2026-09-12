package service

import (
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"
	"math"
)

type TransactionService struct {
	transactionRepository repository.TransactionRepositoryInterface
}

type TransactionHistory struct {
	Items      []model.Transaction `json:"items"`
	Pagination Pagination          `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewTransactionService(
	transactionRepository repository.TransactionRepositoryInterface,
) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
	}
}

func (s *TransactionService) GetByUserID(
	userID uint,
	page int,
	limit int,
) (*TransactionHistory, error) {

	offset := (page - 1) * limit

	transactions, err :=
		s.transactionRepository.FindByUserID(
			userID,
			limit,
			offset,
		)

	if err != nil {
		return nil, err
	}

	total, err :=
		s.transactionRepository.CountByUserID(userID)

	if err != nil {
		return nil, err
	}

	totalPages := int(
		math.Ceil(float64(total) / float64(limit)),
	)

	return &TransactionHistory{
		Items: transactions,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

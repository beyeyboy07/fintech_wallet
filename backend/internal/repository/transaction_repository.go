package repository

import (
	"fintech-wallet/backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepositoryInterface interface {
	Create(transaction *model.Transaction) error
	FindByID(id uint) (*model.Transaction, error)
	FindByUserID(userID uint, limit int, offset int) ([]model.Transaction, error)
	FindByIDForUpdate(id uint) (*model.Transaction, error)
	CountByUserID(userID uint) (int64, error)
	UpdateStatus(id uint, status string) error
	WithTx(tx *gorm.DB) TransactionRepositoryInterface
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) WithTx(
	tx *gorm.DB,
) TransactionRepositoryInterface {
	return &TransactionRepository{
		db: tx,
	}
}

func (r *TransactionRepository) Create(
	transaction *model.Transaction,
) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) FindByID(
	id uint,
) (*model.Transaction, error) {

	var transaction model.Transaction

	err := r.db.
		First(&transaction, id).
		Error

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *TransactionRepository) FindByUserID(
	userID uint,
	limit int,
	offset int,
) ([]model.Transaction, error) {

	var transactions []model.Transaction

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).
		Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) CountByUserID(
	userID uint,
) (int64, error) {

	var total int64

	err := r.db.
		Model(&model.Transaction{}).
		Where("user_id = ?", userID).
		Count(&total).
		Error

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *TransactionRepository) UpdateStatus(
	id uint,
	status string,
) error {
	return r.db.
		Model(&model.Transaction{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

func (r *TransactionRepository) FindByIDForUpdate(
	id uint,
) (*model.Transaction, error) {

	var transaction model.Transaction

	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&transaction, id).
		Error

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

package service

import (
	"errors"
	"time"

	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"

	"gorm.io/gorm"
)

var ErrInvalidTopUpAmount = errors.New("invalid top up amount")

type WalletService struct {
	walletRepository      *repository.WalletRepository
	transactionRepository *repository.TransactionRepository
	db                    *gorm.DB
}

func NewWalletService(
	walletRepository *repository.WalletRepository,
	transactionRepository *repository.TransactionRepository,
	db *gorm.DB,
) *WalletService {
	return &WalletService{
		walletRepository:      walletRepository,
		transactionRepository: transactionRepository,
		db:                    db,
	}
}

func (s *WalletService) GetWallet(
	userID uint,
) (*model.Wallet, error) {

	return s.walletRepository.FindByUserID(userID)
}

func (s *WalletService) TopUp(
	userID uint,
	amount int64,
) (*model.Transaction, error) {

	if amount <= 0 {
		return nil, ErrInvalidTopUpAmount
	}

	var transaction *model.Transaction

	err := s.db.Transaction(func(tx *gorm.DB) error {

		walletRepository := s.walletRepository.WithTx(tx)
		transactionRepository := s.transactionRepository.WithTx(tx)

		// Pastikan wallet user ada
		_, err := walletRepository.FindByUserID(userID)

		if err != nil {
			return err
		}

		// Update saldo
		err = walletRepository.UpdateBalance(
			userID,
			amount,
		)

		if err != nil {
			return err
		}

		// Buat reference sederhana
		reference := "TOPUP-" + time.Now().Format("20060102150405.000000")

		// Buat transaction
		transaction = &model.Transaction{
			UserID:    userID,
			Type:      "TOPUP",
			Amount:    amount,
			Status:    "SUCCESS",
			Reference: reference,
		}

		err = transactionRepository.Create(transaction)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

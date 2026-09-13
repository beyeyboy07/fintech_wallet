package service

import (
	"errors"

	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateWebhook       = errors.New("duplicate webhook event")
	ErrTransactionNotFound    = errors.New("transaction not found")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
	ErrInvalidWebhookAmount   = errors.New("invalid webhook amount")
	ErrTransactionProcessed   = errors.New("transaction already processed")
)

type WebhookService struct {
	db                    *gorm.DB
	webhookRepository     *repository.WebhookRepository
	transactionRepository repository.TransactionRepositoryInterface
	walletRepository      *repository.WalletRepository
}

func NewWebhookService(
	db *gorm.DB,
	webhookRepository *repository.WebhookRepository,
	transactionRepository repository.TransactionRepositoryInterface,
	walletRepository *repository.WalletRepository,
) *WebhookService {
	return &WebhookService{
		db:                    db,
		webhookRepository:     webhookRepository,
		transactionRepository: transactionRepository,
		walletRepository:      walletRepository,
	}
}

func (s *WebhookService) Receive(
	request model.PaymentWebhookRequest,
) error {

	return s.db.Transaction(func(tx *gorm.DB) error {

		webhookRepository := s.webhookRepository.WithTx(tx)
		transactionRepository := s.transactionRepository.WithTx(tx)

		event := &model.WebhookEvent{
			EventID:   request.EventID,
			EventType: "PAYMENT",
			Status:    "PROCESSING",
		}

		err := webhookRepository.Create(event)
		if err != nil {
			var mysqlErr *mysql.MySQLError

			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return ErrDuplicateWebhook
			}

			return err
		}

		// _, err = transactionRepository.FindByIDForUpdate(
		// 	request.TransactionID,
		// )
		// if err != nil {
		// 	return err
		// }

		transaction, err := transactionRepository.FindByIDForUpdate(
			request.TransactionID,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTransactionNotFound
			}

			return err
		}

		if transaction.Type != "TOPUP" {
			return ErrInvalidTransactionType
		}

		if transaction.Amount != request.Amount {
			return ErrInvalidWebhookAmount
		}

		if transaction.Status == "SUCCESS" {
			return ErrTransactionProcessed
		}

		walletRepository := s.walletRepository.WithTx(tx)

		wallet, err := walletRepository.FindByUserIDForUpdate(
			transaction.UserID,
		)
		if err != nil {
			return err
		}

		wallet.Balance += transaction.Amount

		if err := tx.Save(wallet).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&model.Transaction{}).
			Where("id = ?", transaction.ID).
			Update("status", "SUCCESS").
			Error; err != nil {
			return err
		}

		if err := webhookRepository.UpdateStatus(
			event.ID,
			"PROCESSED",
		); err != nil {
			return err
		}

		return nil
	})
}

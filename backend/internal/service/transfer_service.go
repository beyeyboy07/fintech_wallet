package service

import (
	"errors"
	"time"

	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrInvalidTransferAmount = errors.New("invalid transfer amount")
	ErrSelfTransfer          = errors.New("cannot transfer to yourself")
	ErrSenderNotFound        = errors.New("sender wallet not found")
	ErrReceiverNotFound      = errors.New("receiver wallet not found")
	ErrInsufficientBalance   = errors.New("insufficient balance")
)

type TransferService struct {
	walletRepository           *repository.WalletRepository
	transactionRepository      repository.TransactionRepositoryInterface
	transactionEntryRepository *repository.TransactionEntryRepository
	idempotencyRepository      *repository.IdempotencyRepository
	redisService               *RedisService
	db                         *gorm.DB
}

func NewTransferService(
	walletRepository *repository.WalletRepository,
	transactionRepository repository.TransactionRepositoryInterface,
	transactionEntryRepository *repository.TransactionEntryRepository,
	idempotencyRepository *repository.IdempotencyRepository,
	redisService *RedisService,
	db *gorm.DB,
) *TransferService {
	return &TransferService{
		walletRepository:           walletRepository,
		transactionRepository:      transactionRepository,
		transactionEntryRepository: transactionEntryRepository,
		idempotencyRepository:      idempotencyRepository,
		redisService:               redisService,
		db:                         db,
	}
}

func (s *TransferService) Transfer(
	senderID uint,
	receiverID uint,
	amount int64,
	description string,
) (*model.Transaction, error) {

	if amount <= 0 || amount > model.MaxTransferAmount {
		return nil, ErrInvalidTransferAmount
	}

	if senderID == receiverID {
		return nil, ErrSelfTransfer
	}

	var transaction *model.Transaction

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error

		transaction, err = s.transferInTransaction(
			tx,
			senderID,
			receiverID,
			amount,
			description,
		)

		return err
	})

	if err != nil {
		return nil, err
	}

	if s.redisService != nil {
		_ = s.redisService.Delete(walletCacheKey(senderID))
		_ = s.redisService.Delete(walletCacheKey(receiverID))
	}

	return transaction, nil
}

func (s *TransferService) transferInTransaction(
	tx *gorm.DB,
	senderID uint,
	receiverID uint,
	amount int64,
	description string,
) (*model.Transaction, error) {

	walletRepository := s.walletRepository.WithTx(tx)
	transactionRepository := s.transactionRepository.WithTx(tx)
	transactionEntryRepository := s.transactionEntryRepository.WithTx(tx)

	// Cari wallet sender
	senderWallet, err := walletRepository.FindByUserIDForUpdate(senderID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSenderNotFound
		}

		return nil, err
	}

	// Cari wallet receiver
	receiverWallet, err := walletRepository.FindByUserID(receiverID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReceiverNotFound
		}

		return nil, err
	}

	// Cek saldo sender
	if senderWallet.Balance < amount {
		return nil, ErrInsufficientBalance
	}

	// Kurangi saldo sender
	err = walletRepository.UpdateBalance(
		senderID,
		-amount,
	)

	if err != nil {
		return nil, err
	}

	// Tambah saldo receiver
	err = walletRepository.UpdateBalance(
		receiverID,
		amount,
	)

	if err != nil {
		return nil, err
	}

	reference := "TRANSFER-" +
		time.Now().Format("20060102150405.000000")

	senderIDValue := senderID
	receiverIDValue := receiverID

	transaction := &model.Transaction{
		UserID:      senderID,
		SenderID:    &senderIDValue,
		ReceiverID:  &receiverIDValue,
		Type:        "TRANSFER",
		Amount:      amount,
		Status:      "SUCCESS",
		Reference:   reference,
		Description: description,
	}

	err = transactionRepository.Create(transaction)

	if err != nil {
		return nil, err
	}

	err = transactionEntryRepository.Create(&model.TransactionEntry{
		TransactionID: transaction.ID,
		WalletID:      senderID,
		EntryType:     "DEBIT",
		Amount:        amount,
	})

	if err != nil {
		return nil, err
	}

	err = transactionEntryRepository.Create(&model.TransactionEntry{
		TransactionID: transaction.ID,
		WalletID:      receiverID,
		EntryType:     "CREDIT",
		Amount:        amount,
	})

	if err != nil {
		return nil, err
	}

	_ = receiverWallet

	return transaction, nil
}

func (s *TransferService) TransferWithIdempotency(
	senderID uint,
	receiverID uint,
	amount int64,
	description string,
	idempotencyKey string,
) (*model.Transaction, error) {

	if amount <= 0 || amount > model.MaxTransferAmount {
		return nil, ErrInvalidTransferAmount
	}

	if idempotencyKey == "" {
		return nil, ErrIdempotencyKeyRequired
	}

	if amount <= 0 {
		return nil, ErrInvalidTransferAmount
	}

	if senderID == receiverID {
		return nil, ErrSelfTransfer
	}

	requestHash := GenerateRequestHash(
		receiverID,
		amount,
		description,
	)

	var transaction *model.Transaction
	isNewTransfer := false
	err := s.db.Transaction(func(tx *gorm.DB) error {

		idempotencyRepository := s.idempotencyRepository.WithTx(tx)
		transactionRepository := s.transactionRepository.WithTx(tx)

		// Cek apakah key sudah pernah digunakan
		existingKey, err := idempotencyRepository.FindByUserIDAndKey(
			senderID,
			idempotencyKey,
		)

		if err == nil {

			// Key sama, tetapi request berbeda
			if existingKey.RequestHash != requestHash {
				return ErrIdempotencyConflict
			}

			// Request yang sama sudah pernah berhasil
			if existingKey.TransactionID != nil {

				transaction, err = transactionRepository.FindByID(
					*existingKey.TransactionID,
				)

				return err
			}
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
			return err
		}

		// Key belum pernah digunakan
		newKey := &model.IdempotencyKey{
			UserID:      senderID,
			Key:         idempotencyKey,
			RequestHash: requestHash,
		}
		if err := idempotencyRepository.Create(newKey); err != nil {

			var mysqlErr *mysql.MySQLError

			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

				existingKey, findErr :=
					idempotencyRepository.FindByUserIDAndKeyForUpdate(
						senderID,
						idempotencyKey,
					)

				if findErr != nil {
					return findErr
				}

				if existingKey.RequestHash != requestHash {
					return ErrIdempotencyConflict
				}

				if existingKey.TransactionID == nil {
					return errors.New("idempotency transaction is not completed")
				}

				transaction, findErr =
					transactionRepository.FindByID(
						*existingKey.TransactionID,
					)

				if findErr != nil {
					return findErr
				}

				return nil
			}

			return err
		}

		// Jalankan transfer dalam transaction yang sama
		transaction, err = s.transferInTransaction(
			tx,
			senderID,
			receiverID,
			amount,
			description,
		)

		if err != nil {
			return err
		}

		isNewTransfer = true

		// Simpan transaction ID ke idempotency record
		newKey.TransactionID = &transaction.ID

		return idempotencyRepository.Update(newKey)
	})

	if err != nil {
		return nil, err
	}

	if isNewTransfer && s.redisService != nil {
		_ = s.redisService.Delete(walletCacheKey(senderID))
		_ = s.redisService.Delete(walletCacheKey(receiverID))
	}

	return transaction, nil
}

func (s *TransferService) TransferWithIdempotencyByAccountNumber(
	senderID uint,
	receiverAccountNumber string,
	amount int64,
	description string,
	idempotencyKey string,
) (*model.Transaction, error) {

	if len(receiverAccountNumber) != 16 {
		return nil, ErrInvalidAccountNumber
	}

	for _, char := range receiverAccountNumber {
		if char < '0' || char > '9' {
			return nil, ErrInvalidAccountNumber
		}
	}

	receiverWallet, err := s.walletRepository.FindByAccountNumber(
		receiverAccountNumber,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReceiverNotFound
		}

		return nil, err
	}

	return s.TransferWithIdempotency(
		senderID,
		receiverWallet.UserID,
		amount,
		description,
		idempotencyKey,
	)
}

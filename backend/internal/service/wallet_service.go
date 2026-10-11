package service

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"
	"fintech-wallet/backend/internal/worker"

	"gorm.io/gorm"
)

var ErrInvalidTopUpAmount = errors.New("invalid top up amount")
var ErrInvalidAccountNumber = errors.New("invalid account number")

type WalletService struct {
	walletRepository      *repository.WalletRepository
	transactionRepository *repository.TransactionRepository
	outboxRepository      repository.OutboxRepositoryInterface
	redisService          *RedisService
	queueService          *QueueService
	db                    *gorm.DB
}

func NewWalletService(
	walletRepository *repository.WalletRepository,
	transactionRepository *repository.TransactionRepository,
	outboxRepository repository.OutboxRepositoryInterface,
	redisService *RedisService,
	queueService *QueueService,
	db *gorm.DB,
) *WalletService {
	return &WalletService{
		walletRepository:      walletRepository,
		transactionRepository: transactionRepository,
		outboxRepository:      outboxRepository,
		redisService:          redisService,
		queueService:          queueService,
		db:                    db,
	}
}

func (s *WalletService) GetWallet(
	userID uint,
) (*model.Wallet, error) {

	key := walletCacheKey(userID)

	// Coba ambil dari Redis
	cached, err := s.redisService.Get(key)

	if err == nil {
		var wallet model.Wallet

		err := json.Unmarshal(
			[]byte(cached),
			&wallet,
		)

		if err == nil {
			return &wallet, nil
		}
	}

	// Redis miss → ambil dari MySQL
	wallet, err := s.walletRepository.FindByUserID(userID)

	if err != nil {
		return nil, err
	}

	// Simpan hasil MySQL ke Redis
	data, err := json.Marshal(wallet)

	if err == nil {
		_ = s.redisService.Set(
			key,
			string(data),
			5*time.Minute,
		)
	}

	return wallet, nil
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

	// Database transaction sudah berhasil COMMIT.
	// Sekarang hapus wallet cache.
	key := walletCacheKey(userID)

	_ = s.redisService.Delete(key)

	return transaction, nil
}

func walletCacheKey(userID uint) string {
	return "wallet:user:" + strconv.FormatUint(
		uint64(userID),
		10,
	)
}

func (s *WalletService) ProcessTopUp(
	transactionID uint,
	userID uint,
	amount int64,
) error {

	if amount <= 0 {
		return ErrInvalidTopUpAmount
	}

	processed := false

	err := s.db.Transaction(func(tx *gorm.DB) error {

		transactionRepository :=
			s.transactionRepository.WithTx(tx)

		walletRepository :=
			s.walletRepository.WithTx(tx)

		transaction, err :=
			transactionRepository.FindByIDForUpdate(transactionID)

		if err != nil {
			return err
		}

		if transaction.UserID != userID {
			return errors.New("transaction does not belong to user")
		}

		if transaction.Status != "PENDING" {
			return nil
		}

		err = transactionRepository.UpdateStatus(
			transactionID,
			"PROCESSING",
		)

		if err != nil {
			return err
		}

		err = walletRepository.UpdateBalance(
			userID,
			amount,
		)

		if err != nil {
			return err
		}

		processed = true

		return transactionRepository.UpdateStatus(
			transactionID,
			"SUCCESS",
		)
	})

	if err != nil {
		return err
	}

	if processed && s.redisService != nil {
		_ = s.redisService.Delete(
			walletCacheKey(userID),
		)
	}

	return nil
}

func (s *WalletService) TopUpAsync(
	userID uint,
	amount int64,
) (*model.Transaction, error) {

	if amount <= 0 {
		return nil, ErrInvalidTopUpAmount
	}

	var transaction *model.Transaction

	err := s.db.Transaction(func(tx *gorm.DB) error {

		walletRepository :=
			s.walletRepository.WithTx(tx)

		transactionRepository :=
			s.transactionRepository.WithTx(tx)

		outboxRepository :=
			s.outboxRepository.WithTx(tx)

		_, err := walletRepository.FindByUserID(userID)

		if err != nil {
			return err
		}

		reference := "TOPUP-" +
			time.Now().Format("20060102150405.000000")

		transaction = &model.Transaction{
			UserID:    userID,
			Type:      "TOPUP",
			Amount:    amount,
			Status:    "PENDING",
			Reference: reference,
		}

		err = transactionRepository.Create(transaction)

		if err != nil {
			return err
		}

		payload := worker.TopUpPayload{
			UserID:        userID,
			TransactionID: transaction.ID,
			Amount:        amount,
		}

		payloadData, err := json.Marshal(payload)

		if err != nil {
			return err
		}

		outboxEvent := &model.OutboxEvent{
			EventType:   worker.TaskTypeTopUp,
			AggregateID: transaction.ID,
			Payload:     string(payloadData),
			Status:      "PENDING",
		}

		err = outboxRepository.Create(outboxEvent)

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

func (s *WalletService) InquiryAccount(
	accountNumber string,
) (*model.Wallet, error) {
	if len(accountNumber) != 16 {
		return nil, ErrInvalidAccountNumber
	}

	for _, char := range accountNumber {
		if char < '0' || char > '9' {
			return nil, ErrInvalidAccountNumber
		}
	}

	wallet, err := s.walletRepository.FindByAccountNumber(
		accountNumber,
	)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

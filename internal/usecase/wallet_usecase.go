package usecase

import (
	"context"
	"net/http"
	"wallet-service/internal/constants"
	"wallet-service/internal/entity"
	"wallet-service/internal/model"
	"wallet-service/internal/repository"
	"wallet-service/internal/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WalletUseCase struct {
	DB                          *gorm.DB
	Log                         *logrus.Logger
	WalletRepository            *repository.WalletRepository
	WalletTransactionRepository *repository.WalletTransactionRepository
}

func NewWalletUseCase(db *gorm.DB, logger *logrus.Logger, walletRepository *repository.WalletRepository, walletTransactionRepository *repository.WalletTransactionRepository) *WalletUseCase {
	return &WalletUseCase{
		DB:                          db,
		Log:                         logger,
		WalletRepository:            walletRepository,
		WalletTransactionRepository: walletTransactionRepository,
	}
}

func (c *WalletUseCase) GetBalance(ctx context.Context, userID uuid.UUID, walletID uuid.UUID) (*model.WalletBalanceResponse, error) {
	wallet := new(entity.Wallet)

	if err := c.WalletRepository.FindByIDAndUser(c.DB.WithContext(ctx), wallet, walletID, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.Warnf("Wallet not found for user: %s", userID)
			return nil, utils.Error(constants.ErrWalletNotFound, http.StatusNotFound, err)
		}
		c.Log.WithError(err).Error("failed to fetch wallet balance")
		return nil, utils.Error(constants.ErrFetchWalletBalance, http.StatusInternalServerError, err)
	}

	return &model.WalletBalanceResponse{
		Balance:  wallet.Balance,
		Currency: entity.Currency(wallet.Currency),
	}, nil
}

func (c *WalletUseCase) Withdraw(ctx context.Context, userID uuid.UUID, walletID uuid.UUID, amount int64, reference, description string) (*model.WalletWithdrawResponse, error) {
	if amount <= 0 {
		return nil, utils.Error(constants.ErrInvalidAmount, http.StatusBadRequest, nil)
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	wallet := new(entity.Wallet)
	if err := c.WalletRepository.FindByIDAndUserForUpdate(tx, wallet, walletID, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.Warnf("Wallet not found for user: %s", userID)
			return nil, utils.Error(constants.ErrWalletNotFound, http.StatusNotFound, err)
		}
		c.Log.WithError(err).Error("failed to fetch wallet for withdraw")
		return nil, utils.Error(constants.ErrFetchWalletBalance, http.StatusInternalServerError, err)
	}

	if wallet.Balance < amount {
		return nil, utils.Error(constants.ErrInsufficientFunds, http.StatusBadRequest, nil)
	}

	before := wallet.Balance
	after := before - amount

	wallet.Balance = after
	if err := tx.Save(wallet).Error; err != nil {
		c.Log.WithError(err).Error("failed to update wallet balance")
		return nil, utils.Error(constants.ErrFetchWalletBalance, http.StatusInternalServerError, err)
	}

	transaction := &entity.WalletTransaction{
		WalletID:      wallet.ID,
		Type:          "DEBIT",
		Amount:        amount,
		BalanceBefore: before,
		BalanceAfter:  after,
		Reference:     reference,
		Description:   description,
	}

	if err := tx.Create(transaction).Error; err != nil {
		c.Log.WithError(err).Error("failed to record wallet transaction")
		return nil, utils.Error(constants.ErrFetchWalletBalance, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit withdraw transaction")
		return nil, utils.Error(constants.ErrFetchWalletBalance, http.StatusInternalServerError, err)
	}

	return &model.WalletWithdrawResponse{
		Amount:        amount,
		Currency:      wallet.Currency,
		BalanceBefore: before,
		BalanceAfter:  after,
	}, nil
}

func (c *WalletUseCase) List(ctx context.Context, userID uuid.UUID) ([]model.WalletResponse, error) {
	wallets, err := c.WalletRepository.FindAllByUser(c.DB.WithContext(ctx), userID)
	if err != nil {
		c.Log.WithError(err).Error("failed to list wallets")
		return nil, utils.Error(constants.ErrFetchWalletBalance, http.StatusInternalServerError, err)
	}

	responses := make([]model.WalletResponse, 0, len(wallets))
	for _, w := range wallets {
		responses = append(responses, model.WalletResponse{
			ID:       w.ID,
			Currency: entity.Currency(w.Currency),
			Balance:  w.Balance,
		})
	}

	return responses, nil
}

package usecase

import (
	"context"
	"net/http"
	"strings"
	"wallet-service/internal/constants"
	"wallet-service/internal/entity"
	"wallet-service/internal/model"
	"wallet-service/internal/repository"
	"wallet-service/internal/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletUseCase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	WalletRepository *repository.WalletRepository
}

func NewWalletUseCase(db *gorm.DB, logger *logrus.Logger, walletRepository *repository.WalletRepository) *WalletUseCase {
	return &WalletUseCase{
		DB:               db,
		Log:              logger,
		WalletRepository: walletRepository,
	}
}

func (c *WalletUseCase) GetBalance(ctx context.Context, userID uuid.UUID, currency string) (*model.WalletBalanceResponse, error) {
	if strings.TrimSpace(currency) == "" {
		return nil, utils.Error(constants.ErrCurrencyRequired, http.StatusBadRequest, nil)
	}

	currency = strings.ToUpper(currency)

	wallet := new(entity.Wallet)

	tx := c.DB.WithContext(ctx).Clauses(clause.Locking{Strength: "SHARE"})
	if err := c.WalletRepository.FindByUserID(tx, wallet, userID, currency); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.Warnf("Wallet not found for user: %s", userID)
			return nil, utils.Error(constants.ErrWalletNotFound, http.StatusNotFound, err)
		}
		c.Log.WithError(err).Error("failed to fetch wallet balance")
		return nil, utils.Error(constants.ErrFetchWalletBalance, http.StatusInternalServerError, err)
	}

	return &model.WalletBalanceResponse{
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
	}, nil
}

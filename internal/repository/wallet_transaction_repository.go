package repository

import (
	"wallet-service/internal/entity"

	"github.com/sirupsen/logrus"
)

type WalletTransactionRepository struct {
	Repository[entity.WalletTransaction]
	Log *logrus.Logger
}

func NewWalletTransactionRepository(log *logrus.Logger) *WalletTransactionRepository {
	return &WalletTransactionRepository{
		Log: log,
	}
}

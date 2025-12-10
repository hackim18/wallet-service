package repository

import (
	"wallet-service/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func (r *WalletTransactionRepository) FindByWalletIDWithPaging(db *gorm.DB, walletID uuid.UUID, txnType string, limit, offset int) ([]entity.WalletTransaction, int64, error) {
	query := db.Model(&entity.WalletTransaction{}).Where("wallet_id = ?", walletID)
	if txnType != "" {
		query = query.Where("type = ?", txnType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var txs []entity.WalletTransaction
	query = query.Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if err := query.Find(&txs).Error; err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

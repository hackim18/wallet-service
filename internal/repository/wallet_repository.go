package repository

import (
	"wallet-service/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct {
	Repository[entity.Wallet]
	Log *logrus.Logger
}

func NewWalletRepository(log *logrus.Logger) *WalletRepository {
	return &WalletRepository{
		Log: log,
	}
}

func (r *WalletRepository) FindByUserID(db *gorm.DB, wallet *entity.Wallet, userID uuid.UUID, currency string) error {
	return db.Where("user_id = ? AND currency = ?", userID, currency).Take(wallet).Error
}

func (r *WalletRepository) FindByUserIDForUpdate(db *gorm.DB, wallet *entity.Wallet, userID uuid.UUID, currency string) error {
	return db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND currency = ?", userID, currency).
		Take(wallet).Error
}

func (r *WalletRepository) FindByIDAndUserForUpdate(db *gorm.DB, wallet *entity.Wallet, walletID uuid.UUID, userID uuid.UUID) error {
	return db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND user_id = ?", walletID, userID).
		Take(wallet).Error
}

func (r *WalletRepository) FindAllByUser(db *gorm.DB, userID uuid.UUID) ([]entity.Wallet, error) {
	var wallets []entity.Wallet
	if err := db.Where("user_id = ?", userID).Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

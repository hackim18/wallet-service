package migrations

import (
	"wallet-service/internal/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.User{}, &entity.Wallet{}, &entity.WalletTransaction{})
}

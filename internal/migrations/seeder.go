package migrations

import (
	"encoding/json"
	"os"
	"wallet-service/internal/entity"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB, logger *logrus.Logger) error {
	logger.Info("Seeding database...")

	seedFromJSON("internal/migrations/json/users.json", &[]entity.User{}, db, logger)
	seedFromJSON("internal/migrations/json/wallets.json", &[]entity.Wallet{}, db, logger)
	seedFromJSON("internal/migrations/json/wallet_transactions.json", &[]entity.WalletTransaction{}, db, logger)

	return nil
}

func seedFromJSON[T any](filePath string, out *[]T, db *gorm.DB, log *logrus.Logger) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Warnf("Seed file not found: %s", filePath)
		return
	}

	if err := json.Unmarshal(data, out); err != nil {
		log.Warnf("Failed to parse JSON for %s: %v", filePath, err)
		return
	}

	if users, ok := any(out).(*[]entity.User); ok {
		for i := range *users {
			if (*users)[i].PasswordHash == "" {
				hash, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
				(*users)[i].PasswordHash = string(hash)
			}
		}
	}

	var count int64
	if err := db.Model(out).Count(&count).Error; err != nil {
		log.Warnf("Failed to count records for %s: %v", filePath, err)
		return
	}

	if count == 0 {
		if err := db.Create(out).Error; err != nil {
			log.Warnf("Insert failed for %s: %v", filePath, err)
		} else {
			log.Infof("Inserted seed data from %s", filePath)
		}
	} else {
		log.Infof("Skipping insert for %s: table not empty", filePath)
	}
}

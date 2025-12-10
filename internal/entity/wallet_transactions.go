package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletTransaction struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	WalletID      uuid.UUID `gorm:"type:uuid;not null;index:idx_wallet_transactions_wallet_created_at;index:idx_wallet_transactions_wallet" json:"wallet_id"`
	Type          string    `gorm:"type:varchar(20);not null" json:"type"`
	Amount        int64     `gorm:"not null;check:amount > 0" json:"amount"`
	BalanceBefore int64     `gorm:"not null" json:"balance_before"`
	BalanceAfter  int64     `gorm:"not null" json:"balance_after"`
	Reference     string    `gorm:"type:varchar(100)" json:"reference"`
	Description   string    `gorm:"type:text" json:"description"`
	CreatedAt     time.Time `gorm:"type:timestamptz;not null;default:now();index:idx_wallet_transactions_wallet_created_at" json:"created_at"`
	Wallet        Wallet    `gorm:"foreignKey:WalletID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}

func (t *WalletTransaction) BeforeCreate(_ *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}

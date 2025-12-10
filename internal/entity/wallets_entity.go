package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Currency string

const (
	IDR Currency = "IDR"
	USD Currency = "USD"
)

type Wallet struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;index:idx_wallets_user_currency" json:"user_id"`
	Currency  Currency  `gorm:"type:varchar(10);not null;index:idx_wallets_user_currency" json:"currency"`
	Balance   int64     `gorm:"type:bigint;not null;default:0;check:balance >= 0" json:"balance"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli" json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (w *Wallet) UniqueIndex() []interface{} {
	return []interface{}{"user_id", "currency"}
}

func (w *Wallet) TableName() string {
	return "wallets"
}

func (w *Wallet) BeforeCreate(_ *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return
}

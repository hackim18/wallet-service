package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime:milli" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli" json:"updated_at"`
	Wallets      []Wallet  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	return
}

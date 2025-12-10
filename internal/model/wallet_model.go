package model

import (
	"github.com/google/uuid"
	"wallet-service/internal/entity"
)

type WalletResponse struct {
	ID       uuid.UUID       `json:"id"`
	Currency entity.Currency `json:"currency"`
	Balance  int64           `json:"balance"`
}

type WalletBalanceResponse struct {
	Balance  int64           `json:"balance"`
	Currency entity.Currency `json:"currency"`
}

type WalletWithdrawRequest struct {
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	Reference   string `json:"reference,omitempty" validate:"max=100"`
	Description string `json:"description,omitempty" validate:"max=255"`
}

type WalletWithdrawResponse struct {
	Amount        int64           `json:"amount"`
	Currency      entity.Currency `json:"currency"`
	BalanceBefore int64           `json:"balance_before"`
	BalanceAfter  int64           `json:"balance_after"`
}

package model

import (
	"wallet-service/internal/entity"

	"github.com/google/uuid"
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

type WalletDepositRequest struct {
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	Reference   string `json:"reference,omitempty" validate:"max=100"`
	Description string `json:"description,omitempty" validate:"max=255"`
}

type WalletDepositResponse struct {
	WalletID      uuid.UUID       `json:"wallet_id"`
	Amount        int64           `json:"amount"`
	Currency      entity.Currency `json:"currency"`
	BalanceBefore int64           `json:"balance_before"`
	BalanceAfter  int64           `json:"balance_after"`
}

type WalletTransactionResponse struct {
	ID            uuid.UUID   `json:"id"`
	Type          entity.Type `json:"type"`
	Amount        int64       `json:"amount"`
	BalanceBefore int64       `json:"balance_before"`
	BalanceAfter  int64       `json:"balance_after"`
	Reference     string      `json:"reference,omitempty"`
	Description   string      `json:"description,omitempty"`
	CreatedAt     int64       `json:"created_at"`
}

type WalletTransactionsListRequest struct {
	Page int    `json:"page" validate:"omitempty,min=1"`
	Size int    `json:"size" validate:"omitempty,min=1,max=100"`
	Type string `json:"type,omitempty" validate:"omitempty,oneof=DEBIT CREDIT"`
}

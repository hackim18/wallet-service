package model

import "wallet-service/internal/entity"

type WalletBalanceResponse struct {
	Balance  int64           `json:"balance"`
	Currency entity.Currency `json:"currency"`
}

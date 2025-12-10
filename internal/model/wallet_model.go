package model

type WalletBalanceResponse struct {
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

package converter

import (
	"wallet-service/internal/entity"
	"wallet-service/internal/model"
)

func WalletToResponse(wallet *entity.Wallet) model.WalletResponse {
	return model.WalletResponse{
		ID:       wallet.ID,
		Currency: entity.Currency(wallet.Currency),
		Balance:  wallet.Balance,
	}
}

func WalletToBalanceResponse(wallet *entity.Wallet) *model.WalletBalanceResponse {
	return &model.WalletBalanceResponse{
		Balance:  wallet.Balance,
		Currency: entity.Currency(wallet.Currency),
	}
}

func WalletToWithdrawResponse(wallet *entity.Wallet, amount, before, after int64) *model.WalletWithdrawResponse {
	return &model.WalletWithdrawResponse{
		Amount:        amount,
		Currency:      entity.Currency(wallet.Currency),
		BalanceBefore: before,
		BalanceAfter:  after,
	}
}

func WalletToDepositResponse(wallet *entity.Wallet, amount, before, after int64) *model.WalletDepositResponse {
	return &model.WalletDepositResponse{
		WalletID:      wallet.ID,
		Amount:        amount,
		Currency:      entity.Currency(wallet.Currency),
		BalanceBefore: before,
		BalanceAfter:  after,
	}
}

func WalletTransactionsToResponse(txs []entity.WalletTransaction) []model.WalletTransactionResponse {
	responses := make([]model.WalletTransactionResponse, 0, len(txs))
	for _, tx := range txs {
		responses = append(responses, model.WalletTransactionResponse{
			ID:            tx.ID,
			Type:          tx.Type,
			Amount:        tx.Amount,
			BalanceBefore: tx.BalanceBefore,
			BalanceAfter:  tx.BalanceAfter,
			Reference:     tx.Reference,
			Description:   tx.Description,
			CreatedAt:     tx.CreatedAt.Unix(),
		})
	}
	return responses
}

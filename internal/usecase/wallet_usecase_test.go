package usecase_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/hackim18/wallet-service/internal/constants"
	"github.com/hackim18/wallet-service/internal/entity"
	"github.com/hackim18/wallet-service/internal/migrations"
	"github.com/hackim18/wallet-service/internal/model"
	"github.com/hackim18/wallet-service/internal/repository"
	"github.com/hackim18/wallet-service/internal/usecase"
	"github.com/hackim18/wallet-service/internal/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Integration tests run against a real PostgreSQL database (row-level locking cannot be faked).
// Set WALLET_TEST_DATABASE_DSN, e.g.
//
//	host=localhost port=5432 user=postgres password=postgres dbname=wallet_test sslmode=disable
//
// Without it, only the tests that do not touch the database run.
const dsnEnv = "WALLET_TEST_DATABASE_DSN"

func quietLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(io.Discard)
	return l
}

func newUseCase(db *gorm.DB) *usecase.WalletUseCase {
	log := quietLogger()
	return usecase.NewWalletUseCase(db, log, repository.NewWalletRepository(log), repository.NewWalletTransactionRepository(log))
}

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv(dsnEnv)
	if dsn == "" {
		t.Skipf("%s not set; skipping database test", dsnEnv)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("connect database: %v", err)
	}
	if err := migrations.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := db.Exec("TRUNCATE wallet_transactions, wallets, users CASCADE").Error; err != nil {
		t.Fatalf("clean tables: %v", err)
	}
	sqlDB, _ := db.DB()
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}

func createWallet(t *testing.T, db *gorm.DB, balance int64) (userID, walletID uuid.UUID) {
	t.Helper()
	user := &entity.User{Name: "Test User", Email: uuid.NewString() + "@example.com", PasswordHash: "hash"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	wallet := &entity.Wallet{UserID: user.ID, Currency: entity.IDR, Balance: balance}
	if err := db.Create(wallet).Error; err != nil {
		t.Fatalf("create wallet: %v", err)
	}
	return user.ID, wallet.ID
}

func walletBalance(t *testing.T, db *gorm.DB, walletID uuid.UUID) int64 {
	t.Helper()
	var w entity.Wallet
	if err := db.First(&w, "id = ?", walletID).Error; err != nil {
		t.Fatalf("load wallet: %v", err)
	}
	return w.Balance
}

func countTransactions(t *testing.T, db *gorm.DB, walletID uuid.UUID, txType entity.Type) int64 {
	t.Helper()
	var n int64
	if err := db.Model(&entity.WalletTransaction{}).Where("wallet_id = ? AND type = ?", walletID, txType).Count(&n).Error; err != nil {
		t.Fatalf("count transactions: %v", err)
	}
	return n
}

func assertHTTPError(t *testing.T, err error, status int, message string) {
	t.Helper()
	var httpErr utils.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected an HTTPError, got %v", err)
	}
	if httpErr.Status() != status || httpErr.Message() != message {
		t.Fatalf("got status=%d message=%q, want status=%d message=%q", httpErr.Status(), httpErr.Message(), status, message)
	}
}

// --- No database needed: amount validation happens before any query. ---

func TestDepositAndWithdraw_RejectNonPositiveAmount(t *testing.T) {
	uc := newUseCase(nil)
	for _, amount := range []int64{0, -100} {
		_, err := uc.Deposit(context.Background(), uuid.New(), uuid.New(), amount, "", "")
		assertHTTPError(t, err, http.StatusBadRequest, constants.ErrInvalidAmount)

		_, err = uc.Withdraw(context.Background(), uuid.New(), uuid.New(), amount, "", "")
		assertHTTPError(t, err, http.StatusBadRequest, constants.ErrInvalidAmount)
	}
}

// --- Integration tests (PostgreSQL). ---

func TestDeposit_IncreasesBalanceAndRecordsCredit(t *testing.T) {
	db := setupDB(t)
	uc := newUseCase(db)
	userID, walletID := createWallet(t, db, 1_000)

	res, err := uc.Deposit(context.Background(), userID, walletID, 500, "ref-1", "top up")
	if err != nil {
		t.Fatalf("Deposit returned error: %v", err)
	}
	if res.BalanceBefore != 1_000 || res.BalanceAfter != 1_500 || res.Amount != 500 {
		t.Errorf("unexpected response: %+v", res)
	}
	if got := walletBalance(t, db, walletID); got != 1_500 {
		t.Errorf("stored balance = %d, want 1500", got)
	}
	if n := countTransactions(t, db, walletID, entity.CREDIT); n != 1 {
		t.Errorf("CREDIT transactions = %d, want 1", n)
	}
}

func TestWithdraw_DecreasesBalanceAndRecordsDebit(t *testing.T) {
	db := setupDB(t)
	uc := newUseCase(db)
	userID, walletID := createWallet(t, db, 1_000)

	res, err := uc.Withdraw(context.Background(), userID, walletID, 300, "ref-2", "cash out")
	if err != nil {
		t.Fatalf("Withdraw returned error: %v", err)
	}
	if res.BalanceBefore != 1_000 || res.BalanceAfter != 700 {
		t.Errorf("unexpected response: %+v", res)
	}
	if got := walletBalance(t, db, walletID); got != 700 {
		t.Errorf("stored balance = %d, want 700", got)
	}
	if n := countTransactions(t, db, walletID, entity.DEBIT); n != 1 {
		t.Errorf("DEBIT transactions = %d, want 1", n)
	}
}

func TestWithdraw_InsufficientFundsLeavesStateUnchanged(t *testing.T) {
	db := setupDB(t)
	uc := newUseCase(db)
	userID, walletID := createWallet(t, db, 100)

	_, err := uc.Withdraw(context.Background(), userID, walletID, 101, "", "")
	assertHTTPError(t, err, http.StatusBadRequest, constants.ErrInsufficientFunds)

	if got := walletBalance(t, db, walletID); got != 100 {
		t.Errorf("stored balance = %d, want 100 (unchanged)", got)
	}
	if n := countTransactions(t, db, walletID, entity.DEBIT); n != 0 {
		t.Errorf("DEBIT transactions = %d, want 0", n)
	}
}

func TestWallet_OtherUsersWalletIsNotFound(t *testing.T) {
	db := setupDB(t)
	uc := newUseCase(db)
	_, walletID := createWallet(t, db, 1_000)
	strangerID, _ := createWallet(t, db, 0)

	_, err := uc.Withdraw(context.Background(), strangerID, walletID, 10, "", "")
	assertHTTPError(t, err, http.StatusNotFound, constants.ErrWalletNotFound)

	_, err = uc.GetBalance(context.Background(), strangerID, walletID)
	assertHTTPError(t, err, http.StatusNotFound, constants.ErrWalletNotFound)

	if got := walletBalance(t, db, walletID); got != 1_000 {
		t.Errorf("stored balance = %d, want 1000 (unchanged)", got)
	}
}

// Concurrent withdrawals must never overdraw the wallet: SELECT ... FOR UPDATE serializes them,
// so exactly balance/amount requests succeed and the rest fail with insufficient funds.
func TestWithdraw_ConcurrentRequestsNeverOverdraw(t *testing.T) {
	db := setupDB(t)
	uc := newUseCase(db)
	userID, walletID := createWallet(t, db, 1_000)

	const requests, amount = 25, 100
	var (
		wg           sync.WaitGroup
		mu           sync.Mutex
		succeeded    int
		insufficient int
		otherErrs    []error
	)
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := uc.Withdraw(context.Background(), userID, walletID, amount, "", "")
			mu.Lock()
			defer mu.Unlock()
			var httpErr utils.HTTPError
			switch {
			case err == nil:
				succeeded++
			case errors.As(err, &httpErr) && httpErr.Message() == constants.ErrInsufficientFunds:
				insufficient++
			default:
				otherErrs = append(otherErrs, err)
			}
		}()
	}
	wg.Wait()

	if len(otherErrs) > 0 {
		t.Fatalf("unexpected errors: %v", otherErrs)
	}
	if succeeded != 10 || insufficient != requests-10 {
		t.Errorf("succeeded=%d insufficient=%d, want 10 and %d", succeeded, insufficient, requests-10)
	}
	if got := walletBalance(t, db, walletID); got != 0 {
		t.Errorf("stored balance = %d, want 0", got)
	}
	if n := countTransactions(t, db, walletID, entity.DEBIT); n != 10 {
		t.Errorf("DEBIT transactions = %d, want 10", n)
	}
}

func TestListTransactions_PagingAndTypeFilter(t *testing.T) {
	db := setupDB(t)
	uc := newUseCase(db)
	userID, walletID := createWallet(t, db, 0)
	ctx := context.Background()

	for i := 0; i < 4; i++ {
		if _, err := uc.Deposit(ctx, userID, walletID, 100, "", ""); err != nil {
			t.Fatalf("Deposit returned error: %v", err)
		}
	}
	if _, err := uc.Withdraw(ctx, userID, walletID, 50, "", ""); err != nil {
		t.Fatalf("Withdraw returned error: %v", err)
	}

	items, page, err := uc.ListTransactions(ctx, userID, walletID, model.WalletTransactionsListRequest{Page: 1, Size: 2})
	if err != nil {
		t.Fatalf("ListTransactions returned error: %v", err)
	}
	if len(items) != 2 || page.TotalItem != 5 || page.TotalPage != 3 || !page.HasNext || page.HasPrevious {
		t.Errorf("unexpected page: items=%d meta=%+v", len(items), page)
	}

	credits, page, err := uc.ListTransactions(ctx, userID, walletID, model.WalletTransactionsListRequest{Type: "credit"})
	if err != nil {
		t.Fatalf("ListTransactions (CREDIT) returned error: %v", err)
	}
	if len(credits) != 4 || page.TotalItem != 4 {
		t.Errorf("CREDIT filter: items=%d total=%d, want 4 and 4", len(credits), page.TotalItem)
	}
	for _, c := range credits {
		if c.Type != entity.CREDIT {
			t.Errorf("CREDIT filter returned a %s transaction", c.Type)
		}
	}

	_, _, err = uc.ListTransactions(ctx, userID, walletID, model.WalletTransactionsListRequest{Type: "REFUND"})
	assertHTTPError(t, err, http.StatusBadRequest, constants.ErrInvalidTransactionType)
}

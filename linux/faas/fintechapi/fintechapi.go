package fintechapi

import (
	"log"
	"os"
)

// FintechAPI defines the interface for the fintech API
type AccountAPI interface {
	CreateAccount(accountID string, initialBalance float64) error
	GetAccountBalance(accountID string) (float64, error)
	TransferFunds(fromAccountID, toAccountID string, amount float64) error
	ListTransactions(accountID string) ([]Transaction, error)
}

// Transaction represents a financial transaction
type Transaction struct {
	ID        string
	AccountID string
	Amount    float64
	Date      string
	Type      string
}

// StockAPI defines the interface for stock price queries
type IStockAPI interface {
	GetStockPrice(ticker string) (float32, error)
}

type StockAPI struct {
}

func (s StockAPI) GetStockPrice(ticker string) (float32, error) {
	api := NewYHFinanceCompleteAPI(log.New(os.Stdout, "", log.Ltime))

	resp, err := api.GetStockPrice(ticker); if err != nil {
		return 0, err
	}

	return resp.Price, nil
}

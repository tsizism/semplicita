package fintechapi

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
type StockAPI interface {
	GetStockPrice(ticker string) (float64, error)
	ListStockPrices(tickers []string) (map[string]float64, error)
}

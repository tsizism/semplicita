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

	resp, err := api.GetStockPrice(ticker)
	if err != nil {
		return 0, err
	}

	return resp.Price, nil
}

// RPCServer defines the server for handling RPC calls
type RPCServer struct {
	stockAPI IStockAPI
}

// NewRPCServer creates a new RPC server
func NewRPCServer(stockAPI IStockAPI) *RPCServer {
	return &RPCServer{stockAPI: stockAPI}
}

// GetStockPriceRequest represents the request for getting stock price
type GetStockPriceRequest struct {
	Ticker string
}

// GetStockPriceResponse represents the response for getting stock price
type GetStockPriceResponse struct {
	Price float32
	Error string
}

// GetStockPrice handles the RPC call to get stock price
func (s *RPCServer) GetStockPrice(req GetStockPriceRequest, res *GetStockPriceResponse) error {
	price, err := s.stockAPI.GetStockPrice(req.Ticker)
	if err != nil {
		res.Error = err.Error()
		return err
	}
	res.Price = price
	return nil
}

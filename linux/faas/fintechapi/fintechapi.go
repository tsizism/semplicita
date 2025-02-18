package fintechapi

import (
	"fmt"
	"log"
	"strings"
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

func NewStockAPI(logger *log.Logger) IStockAPI {
	return StockAPI{
		yahooApi: NewYHFinanceCompleteAPI(logger),
	}
}

// StockAPI defines the interface for stock price queries
type IStockAPI interface {
	GetSingleStockPriceNum(ticker string) (float32, error)
	GetStocksPriceCSV(tickersCSV string) (string, error)
	GetStocksFullPriceCSV(tickersCSV string) (string, error)
}

type StockAPI struct {
	yahooApi YHFinanceCompleteAPI
}

func (s StockAPI) GetStocksFullPriceCSV(tickersCSV string) (string, error) {
	s.yahooApi.logger.Printf("StockAPI.GetStocksFullPriceCSV stock prices for %s", tickersCSV)

	tickers := strings.Split(tickersCSV, ",")

	result := ""
	stockPriceFmt := ""

	for _, ticker := range tickers {
		resp, err := s.yahooApi.GetSingleStockFullPrice(strings.TrimSpace(ticker)); if err != nil {
			return "", err
		}

		if resp.Price.RegularMarketChange < 0 { stockPriceFmt = "%s:%.2f %.2f (%.2f%%),"} else { stockPriceFmt = "%s:%.2f +%.2f (+%.2f%%)," }
		stockInfo := fmt.Sprintf(stockPriceFmt, ticker, resp.Price.RegularMarketPrice, resp.Price.RegularMarketChange, resp.Price.RegularMarketChangePercent*100)
		result += stockInfo
	}

	return result, nil
}


func (s StockAPI) GetStocksPriceCSV(tickersCSV string) (string, error) {
	s.yahooApi.logger.Printf("StockAPI.GetStocksPriceCSV stock prices for %s", tickersCSV)

	tickers := strings.Split(tickersCSV, ",")

	result := ""

	for _, ticker := range tickers {
		resp, err := s.yahooApi.GetSingleStockPrice(strings.TrimSpace(ticker));	if err != nil {
			return "", err
		}
		result += ticker + ":" + fmt.Sprintf("%.2f",resp.Price) + ","
	}

	return result, nil
}

func (s StockAPI) GetSingleStockPriceNum(ticker string) (float32, error) {
	s.yahooApi.logger.Printf("StockAPI.GetSingleStockPriceNum stock prices for %s", ticker)
	resp, err := s.yahooApi.GetSingleStockPrice(ticker)
	if err != nil {
		return 0, err
	}
	return resp.Price, nil
}

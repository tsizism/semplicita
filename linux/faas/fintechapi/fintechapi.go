package fintechapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
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
	Shutdown()
}

type StockAPI struct {
	yahooApi YHFinanceCompleteAPI
}

func (s StockAPI) Shutdown() {
	s.yahooApi.logger.Printf("Shutdown")
}

// func (s StockAPI) GetStocksFullPriceCSV(tickersCSV string) (string, error) {
// 	s.yahooApi.logger.Printf("StockAPI.GetStocksFullPriceCSV stock prices for %s", tickersCSV)

// 	tickers := strings.Split(tickersCSV, ",")

// 	result := ""
// 	stockPriceFmt := ""

// 	for _, ticker := range tickers {
// 		resp, err := s.yahooApi.getSingleStockFullPrice(strings.TrimSpace(ticker)); if err != nil {
// 			return "", err
// 		}

// 		if resp.Price.RegularMarketChange < 0 { stockPriceFmt = "%s:%.2f %.2f (%.2f%%),"} else { stockPriceFmt = "%s:%.2f +%.2f (+%.2f%%)," }
// 		stockInfo := fmt.Sprintf(stockPriceFmt, ticker, resp.Price.RegularMarketPrice, resp.Price.RegularMarketChange, resp.Price.RegularMarketChangePercent*100)
// 		result += stockInfo
// 	}

// 	return result, nil
// }

func (s StockAPI) GetStocksFullPriceCSV(tickersCSV string) (string, error) {
	s.yahooApi.logger.Printf("StockAPI.GetStocksFullPriceCSV stock prices for %s", tickersCSV)

	tickers := strings.Split(tickersCSV, ",")

	result := ""
	var lastErr error

	// for _, ticker := range tickers {
	// 	resp, err := s.yahooApi.getSingleStockPrice(strings.TrimSpace(ticker));	if err != nil {
	// 		return "", err
	// 	}
	// 	result += ticker + ":" + fmt.Sprintf("%.2f",resp.Price) + ","
	// }

	// count := len(tickers) time.Duration(count)
	ctx, cancel := context.WithTimeout(context.Background(), 100 *time.Second)
	defer cancel()
	stockPriceFmt := ""

	var wg sync.WaitGroup

	wg.Add(len(tickers))
	go func() {
		for i := 0; i < len(tickers); i++ {
			var stockFullPrice YffullstockpriceResponse
			select {
			case stockFullPrice = <-s.yahooApi.stockFullPriceCh:
			case lastErr = <-s.yahooApi.stockFullPriceErrCh:
			case <-ctx.Done():
				println("done")
				lastErr = errors.New("time out")
			}

			if lastErr == nil {
				if stockFullPrice.Price.RegularMarketChange < 0 {
					stockPriceFmt = "%s:%.2f %.2f (%.2f%%),"
				} else {
					stockPriceFmt = "%s:%.2f +%.2f (+%.2f%%),"
				}
				stockInfo := fmt.Sprintf(stockPriceFmt, stockFullPrice.Price.Symbol, stockFullPrice.Price.RegularMarketPrice, stockFullPrice.Price.RegularMarketChange, stockFullPrice.Price.RegularMarketChangePercent*100)
				result += stockInfo
			}
			wg.Done()
		}

		// for stockFullPrice := range s.yahooApi.stockFullPriceCh {
		// }
	}()

	for _, ticker := range tickers {
		println(ticker)
		s.yahooApi.symbolCh <- Symbol(ticker)
		println("send")
	}

	wg.Wait()

	for result == "" && lastErr != nil {
		return "", lastErr
	}

	return result, nil
}

func (s StockAPI) GetStocksPriceCSV(tickersCSV string) (string, error) {
	s.yahooApi.logger.Printf("StockAPI.GetStocksPriceCSV stock prices for %s", tickersCSV)

	tickers := strings.Split(tickersCSV, ",")

	result := ""

	for _, ticker := range tickers {
		resp, err := s.yahooApi.getSingleStockPrice(strings.TrimSpace(ticker))
		if err != nil {
			return "", err
		}
		result += ticker + ":" + fmt.Sprintf("%.2f", resp.Price) + ","
	}

	return result, nil
}

// func (s StockAPI) GetSingleStockPriceNum2(ticker string) (float32, error) {
// 	s.yahooApi.logger.Printf("StockAPI.GetSingleStockPriceNum stock prices for %s", ticker)
// 	resp, err := s.yahooApi.getSingleStockPrice(ticker)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return resp.Price, nil
// }

func (s StockAPI) GetSingleStockPriceNum(ticker string) (float32, error) {
	s.yahooApi.logger.Printf("StockAPI.GetSingleStockPriceNum stock prices for %s", ticker)

	// resp, err := s.yahooApi.GetSingleStockPrice(ticker)

	var resp YfpriceResponse
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		println("go YfpriceResponse")
		select {
		case resp = <-s.yahooApi.priceResponseCh:
			println("resp YfpriceResponse")
		case err = <-s.yahooApi.priceResponseErrCh:
			println("priceResponseErrCh")
			fmt.Printf("%v\n", err)
		case <-ctx.Done():
			println("resp YfpriceResponse TO")
			err = errors.New("time out")
		}
		println("go YfpriceResponse done")
		wg.Done()
	}()

	s.yahooApi.tickerCh <- Ticker(ticker)
	wg.Wait()

	if err != nil {
		return 0, err
	}

	return resp.Price, nil
}

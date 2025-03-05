package fintechapi

import (
	"context"
	"fmt"
	"log"
	"runtime"
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
		reqTimeoutSec: 2,
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
	reqTimeoutSec int
	yahooApi YHFinanceCompleteAPI
}

func (s StockAPI) Shutdown() {
	s.yahooApi.logger.Printf("Shutdown StockAPI")
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

	count := len(tickers)
	stockPriceFmt := ""

	var mutext sync.Mutex 
	var wg sync.WaitGroup
	wg.Add(count)
	var stockFullPrice YffullstockpriceResponse
	// sem := semaphore.NewWeighted(2)
	for i := 0; i < count; i++ {
		fmt.Println(runtime.NumGoroutine())
		// err := sem.Acquire(context.Background(), 1); if err != nil {
		//  	log.Fatal(err)
		// }

		go func() {
			// defer sem.Release(1)

			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.reqTimeoutSec) * time.Second); defer cancel()

			select {
			case stockFullPrice = <-s.yahooApi.stockFullPriceCh:
				// fmt.Printf("stockFullPrice for= %s\n", stockFullPrice.Price.Symbol)
				if stockFullPrice.Price.RegularMarketChange < 0 {
					stockPriceFmt = "%s:%.2f %.2f (%.2f%%),"
				} else {
					stockPriceFmt = "%s:%.2f +%.2f (+%.2f%%),"
				}
				stockInfo := fmt.Sprintf(stockPriceFmt, stockFullPrice.Price.Symbol, stockFullPrice.Price.RegularMarketPrice, stockFullPrice.Price.RegularMarketChange, stockFullPrice.Price.RegularMarketChangePercent*100)
				mutext.Lock()
				result += stockInfo
				mutext.Unlock()
			case lastErr = <-s.yahooApi.stockFullPriceErrCh:
				fmt.Printf("lastErr %s\n", lastErr)
			case <-ctx.Done():
				deadline,_:=ctx.Deadline()
				lastErr = ctx.Err()
				fmt.Printf("====>GetStocksFullPriceCSV Deadline err=%v, val=%v\n", lastErr, time.Until(deadline))
			}
			fmt.Printf("GetStocksFullPriceCSV i=%d\n", i)
		}()
	}
	// println("Finished loop")


	for _, ticker := range tickers {
		// println(ticker)
		s.yahooApi.stockFullPriceTickerCh <- Ticker(strings.TrimSpace(ticker))
		// println("ticker sent")
	}

	// println("waiting ... ")
	wg.Wait()
	// println("waiting done")

	for result == "" && lastErr != nil {
		return "", lastErr
	}

	return result, nil
}

func (s StockAPI) GetSingleStockPriceNum(ticker string) (float32, error) {
	s.yahooApi.logger.Printf("StockAPI.GetSingleStockPriceNum stock prices for %s", ticker)

	// resp, err := s.yahooApi.GetSingleStockPrice(ticker)

	var resp YfpriceResponse
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.reqTimeoutSec) * time.Second); defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// println("go YfpriceResponse")
		select {
		case resp = <-s.yahooApi.priceResponseCh:
			// println("resp YfpriceResponse")
		case err = <-s.yahooApi.priceResponseErrCh:
			// println("priceResponseErrCh")
			fmt.Printf("%v\n", err)
		case <-ctx.Done():
			deadline,_:=ctx.Deadline()
			err = ctx.Err()
			fmt.Printf("====>GetStocksFullPriceCSV Deadline err=%v, val=%v\n", err, time.Until(deadline))
	}
		wg.Done()
	}()

	s.yahooApi.tickerCh <- Ticker(ticker)
	wg.Wait()

	if err != nil {
		return 0, err
	}

	return resp.Price, nil
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


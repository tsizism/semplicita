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
		reqTimeoutSec: 1,
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

	// consumer

	countdown := count
	unclaimed := 0
	for i := 0; i < count; i++ {
		// err := sem.Acquire(context.Background(), 1); if err != nil {
		//  	log.Fatal(err)
		// }

		go func(id int) {
			// defer sem.Release(1)
			defer wg.Done()
			defer fmt.Printf("done GetStocksFullPriceCSV id=%d, countdown=%d\n", id, countdown)
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.reqTimeoutSec) * time.Second); defer cancel()
			countdown--

			// select will block waiting for any messages on the channels. 
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
					return
				case lastErr = <-s.yahooApi.stockFullPriceErrCh:
					fmt.Printf("lastErr %s\n", lastErr)
					return
				case <-ctx.Done():
					deadline,_:=ctx.Deadline()
					lastErr = ctx.Err()
					unclaimed++
					fmt.Printf("====>consumer GetStocksFullPriceCSV Deadline err=%v, val=%v, id=%d unclaimed=%d\n", lastErr, time.Until(deadline), id, unclaimed)
					return
				}
		}(i)
	}

	// producer
	for _, ticker := range tickers {
		println("sending " + ticker)
		s.yahooApi.stockFullPriceTickerCh <- Ticker(strings.TrimSpace(ticker))
		println("ticker sent " + ticker)
	}
	println("waiting ... ")
	wg.Wait()
	println("waiting done")

	fmt.Printf("unclaimed=%d, countdown=%d\n", unclaimed, countdown)
	// for i:=0; i < unclaimed; i++ {
	// 	<-s.yahooApi.stockFullPriceCh
	// }

	for result == "" && lastErr != nil {
		return "", lastErr
	}
	fmt.Println(runtime.NumGoroutine())

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


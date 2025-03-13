package fintechapi

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const TIMEOUT_EXCHANGE_SEC int = 5
const TIMOUT_HTTPREQ_MS int = 5000

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
		yahooApi:           NewYHFinanceCompleteAPI(logger),
		exchangeTimeoutSec: TIMEOUT_EXCHANGE_SEC,
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
	exchangeTimeoutSec int
	yahooApi           YHFinanceCompleteAPI
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

	var resultMutext sync.Mutex
	var wg sync.WaitGroup
	wg.Add(count)
	sem := semaphore.NewWeighted(10)

	fmt.Printf("______________________>NumGoroutine before=%d\n", runtime.NumGoroutine())
	countdown := count
	unclaimed := 0

	for i := 0; i < count; i++ {
		go func(id int) {
			println("Acquiring semafore id=%d", id)
			if err := sem.Acquire(context.Background(), 1); err != nil {
				log.Fatal(err)
			}
			defer sem.Release(1)
			defer wg.Done()
			defer fmt.Printf("done GetStocksFullPriceCSV id=%d, countdown=%d\n", id, countdown)

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.exchangeTimeoutSec)*time.Second)
			defer cancel()
			countdown--

			// select will block waiting for any messages on the channels.
			select {
			case stockFullPrice := <-s.yahooApi.stockFullPriceCh:
				// if stockFullPrice.Price.RegularMarketPrice > 0 {
				resultMutext.Lock()
				result += createStockInfoLine(stockFullPrice, id)
				resultMutext.Unlock()
				return
			case lastErr = <-s.yahooApi.stockFullPriceErrCh:
				fmt.Printf("lastErr %s\n", lastErr)
				return
			case <-ctx.Done():
				deadline, _ := ctx.Deadline()
				lastErr = ctx.Err()
				unclaimed++
				fmt.Printf("====>consumer GetStocksFullPriceCSV Deadline err=%v, val=%v, id=%d unclaimed=%d\n", lastErr, time.Until(deadline), id, unclaimed)
				return
			}
		}(i)
	}

	// producer
	for _, ticker := range tickers {
		// println("-------->sending " + ticker)
		s.yahooApi.stockFullPriceTickerCh <- Ticker(strings.TrimSpace(ticker))
		// println("ticker sent " + ticker)
	}
	// println("waiting ... ")
	wg.Wait()
	// println("waiting done")

	fmt.Printf("unclaimed=%d, countdown=%d\n", unclaimed, countdown)
	// for i:=0; i < unclaimed; i++ {
	// 	<-s.yahooApi.stockFullPriceCh
	// }

	fmt.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<NumGoroutine after=%d\n", runtime.NumGoroutine())
	for result == "" && lastErr != nil {
		return "", lastErr
	}

	fmt.Println("result->" + result)

	fmt.Println("((((((((((((((sessionEndCh))))))))))))))")
	s.yahooApi.sessionEndCh <- true

	return result, nil
}

func createStockInfoLine(stockFullPrice YffullstockpriceResponse, id int) string {
	stockPriceFmt := ""
	if stockFullPrice.Price.RegularMarketChange < 0 {
		stockPriceFmt = "%s:%.2f %.2f (%.2f%%)"
	} else {
		stockPriceFmt = "%s:%.2f +%.2f (+%.2f%%)"
	}
	stockInfo := fmt.Sprintf(stockPriceFmt, stockFullPrice.Price.Symbol, stockFullPrice.Price.RegularMarketPrice, stockFullPrice.Price.RegularMarketChange, stockFullPrice.Price.RegularMarketChangePercent*100)

	// t := time.Now()
	// stockInfo += fmt.Sprintf(" %02d:%02d", t.Second(), t.Nanosecond())

	stockInfo += ","

	fmt.Printf("id=%d createStockInfoLine=%s\n", id, stockInfo)
	return stockInfo
}

func (s StockAPI) GetSingleStockPriceNum(ticker string) (float32, error) {
	s.yahooApi.logger.Printf("StockAPI.GetSingleStockPriceNum stock prices for %s", ticker)

	// resp, err := s.yahooApi.GetSingleStockPrice(ticker)

	var resp YfpriceResponse
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.exchangeTimeoutSec)*time.Second)
	defer cancel()

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
			deadline, _ := ctx.Deadline()
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

package fintechapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://query1.finance.yahoo.com/v7/finance/quote?symbols="

type QuoteResponse struct {
	QuoteResponse struct {
		Result []struct {
			Symbol             string  `json:"symbol"`
			RegularMarketPrice float64 `json:"regularMarketPrice"`
			RegularMarketTime  int64   `json:"regularMarketTime"`
		} `json:"result"`
	} `json:"quoteResponse"`
}

func getQuote(symbol string) (*QuoteResponse, error) {
	resp, err := http.Get(baseURL + symbol)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var quoteResponse QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quoteResponse); err != nil {
		return nil, err
	}

	return &quoteResponse, nil
}

func main() {
	symbol := "AAPL"
	quote, err := getQuote(symbol)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(quote.QuoteResponse.Result) > 0 {
		fmt.Printf("Symbol: %s\n", quote.QuoteResponse.Result[0].Symbol)
		fmt.Printf("Price: %.2f\n", quote.QuoteResponse.Result[0].RegularMarketPrice)
		fmt.Printf("Time: %d\n", quote.QuoteResponse.Result[0].RegularMarketTime)
	} else {
		fmt.Println("No data found for symbol:", symbol)
	}
}
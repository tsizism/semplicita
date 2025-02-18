package main

import (
	"faas/fintechapi"
	"log"
)

// RPCServer defines the server for handling RPC calls
type RPCServer struct {
	stockAPI fintechapi.IStockAPI
}

// NewRPCServer creates a new RPC server
func NewRPCServer(logger *log.Logger) *RPCServer {
	return &RPCServer{stockAPI: fintechapi.NewStockAPI(logger)}
}

// GetSingleStockPriceRequest represents the request for getting stock price
type GetSingleStockPriceRequest struct {
	Ticker string
}

// GetSingleStockPriceResponse represents the response for getting stock price
type GetSingleStockPriceResponse struct {
	Price float32
	Error string
}

type GetStocksPriceRequest struct {
	TickersCSV string
}

// GetSingleStockPriceResponse represents the response for getting stock price
type GetStocksPriceResponse struct {
	TickerPriceCSV string
	Error string
}

type GetStocksPriceReqResp struct {
	TickerPriceCSV string
	Error string
}


// GetSingleStockPrice handles the RPC call to get stock price
func (s *RPCServer) GetSingleStockPrice(req GetSingleStockPriceRequest, res *GetSingleStockPriceResponse) error {
	price, err := s.stockAPI.GetSingleStockPrice(req.Ticker); if err != nil {
		res.Error = err.Error()
		return err
	}

	res.Price = price
	return nil
}

func (s *RPCServer) GetStocksPrice(req GetStocksPriceReqResp, res *GetStocksPriceReqResp) error {
	tickerPriceCSV, err := s.stockAPI.GetStocksPrice(req.TickerPriceCSV); if err != nil {
		res.Error = err.Error()
		return err
	}

	res.TickerPriceCSV = tickerPriceCSV
	return nil
}

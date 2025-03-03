package main

import (
	"faas/fintechapi"
	"fmt"
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

type GetSingleStockPriceReqResp struct {
	TickerPrice string
	Error       string
}

// type GetStocksPriceRequest struct {
// 	TickersCSV string
// }

// // GetSingleStockPriceResponse represents the response for getting stock price
// type GetStocksPriceResponse struct {
// 	TickerPriceCSV string
// 	Error string
// }

type GetStocksPriceReqResp struct {
	TickerPriceCSV string
	Error          string
}

func (s *RPCServer) Shutdown() {
	s.stockAPI.Shutdown()
}

// GetSingleStockPrice handles the RPC call to get stock price
func (s *RPCServer) GetSingleStockPriceTxt(req GetSingleStockPriceReqResp, res *GetSingleStockPriceReqResp) error {
	price, err := s.stockAPI.GetSingleStockPriceNum(req.TickerPrice)
	if err != nil {
		res.Error = err.Error()
		return err
	}

	res.TickerPrice = fmt.Sprintf("%.2f", price)
	return nil
}

func (s *RPCServer) GetStocksPriceCSV(req GetStocksPriceReqResp, res *GetStocksPriceReqResp) error {
	tickerPriceCSV, err := s.stockAPI.GetStocksPriceCSV(req.TickerPriceCSV)
	if err != nil {
		res.Error = err.Error()
		return err
	}

	res.TickerPriceCSV = tickerPriceCSV
	return nil
}

func (s *RPCServer) GetStocksFullPriceCSV(req GetStocksPriceReqResp, res *GetStocksPriceReqResp) error {
	tickerPriceCSV, err := s.stockAPI.GetStocksFullPriceCSV(req.TickerPriceCSV)
	if err != nil {
		res.Error = err.Error()
		return err
	}

	res.TickerPriceCSV = tickerPriceCSV
	return nil
}

package main

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

const (
	address = "localhost:50051"
)

type applicationContext struct{
	logger *log.Logger 
}

func unwrapErrors(prfix string, logger *log.Logger, err error, deep int) {
	fmt := "%s(%d):%v\n"
	i := 0
	logger.Printf(fmt, prfix, i, err)
	for {
		if deep != 0 && i >= deep-1 {
			break
		}

		err = errors.Unwrap(err); if err != nil {
			i++
			logger.Printf(fmt, prfix, i, err)
		} else {
			break;
		}
	}
}

type GetSingleStockPriceReqResp struct {
	TickerPrice string
	Error string
}

type GetStocksPriceReqRespCSV struct {
	TickerPriceCSV string
	Error string
}


func main() {
	appCtx := applicationContext {
		logger: log.New(os.Stdout, "", log.Ltime),
	}

	hit1(appCtx)
}

func hit1(appCtx applicationContext) {
	txt, err := appCtx.getStocksFullPricCSV(); if err != nil {
		unwrapErrors("getStocksFullPricCSV", appCtx.logger, err, 1)
	}
	fmt.Printf("Full price CSV: %s", txt)
}

func hit2(appCtx applicationContext) {
	txt, err := appCtx.getSingleStockPriceTxt(); if err != nil {
		unwrapErrors("getSingleStockPriceTxt", appCtx.logger, err, 1)
	}
	fmt.Printf("Stock Price: %s", txt)
}

func (appCtx *applicationContext) getStocksFullPricCSV() (string, error){
	tickers := "BCE.TO,BCE, CM, CM.TO"
	appCtx.logger.Printf("getStocksFullPriceCSV: payload=%+v", tickers)
	client, err := rpc.Dial("tcp", "localhost:5003") // faas-service:5003

	if err != nil {
		// shared.ErrorJSON(w, err)
		// return
		return "", fmt.Errorf("%w", err)
	}

	type GetStocksPriceReqRespCSV struct {
		TickerPriceCSV string
		Error string
	}

	var result, rpcPayload GetStocksPriceReqRespCSV
	rpcPayload.TickerPriceCSV = tickers
	
	err = client.Call("RPCServer.GetStocksFullPriceCSV", rpcPayload, &result);	if err != nil {
		return "", fmt.Errorf("GetStocksFullPrice: Failed to call RPC %w", err)
		// appCtx.logger.Printf("GetStocksFullPrice: Failed to call RPC %s", err)
		// shared.ErrorJSON(w, err)
		// return
	}

	appCtx.logger.Printf("GetStocksFullPriceCSV for %s: result=%+v", tickers, result)

	return result.TickerPriceCSV, nil

	// payload := shared.JsonResponse{
	// 	Error:   false,
	// 	Message: fmt.Sprintf("TickersPriceCSV: %s", result.TickerPriceCSV),
	// }

	// shared.WriteJSON(w, http.StatusAccepted, payload)
}


func (appCtx applicationContext) getSingleStockPriceTxt() (string,error) {
	ticker := "AAPL"
	txt := ""
	appCtx.logger.Printf("getSingleStockPriceTxt: payload=%+v", ticker)
	
	// faas-service:5003
	client, err := rpc.Dial("tcp", "localhost:5003"); if err != nil {
		errW := fmt.Errorf("failed to create RPC err=%w", err)
		return txt, errW
	}

	var rpcPayload,result GetSingleStockPriceReqResp
	rpcPayload.TickerPrice = ticker
	
	err = client.Call("RPCServer.GetSingleStockPriceTxt", rpcPayload, &result)
	
	if err != nil {
		errW := fmt.Errorf("failed to call RPC %w", err)
		return txt, errW
	}

	appCtx.logger.Printf("GetSingleStockPriceTxt for%s: result=%+v", ticker, result)

	// payload := shared.JsonResponse{
	// 	Error:   false,
	// 	Message: fmt.Sprintf("Stock Price: %f", result.TickerPrice),
	// }

	txt = result.TickerPrice
	// shared.WriteJSON(w, http.StatusAccepted, payload)

	return txt, nil
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"slices"
	"strings"

	"github.com/fatih/color"
)

const (
	address = "localhost:50051"
)

type applicationContext struct{
	logger *log.Logger 
}

func unwrapErrors(prfix string, logger *log.Logger, err error, deep int) {
	txtFmt := "Error: %s(%d):%v\n"
	i := 0
	txt := fmt.Sprintf(txtFmt, prfix, i, err)
	color.HiMagenta(txt)
	// logger.Print(txt)
	for {
		if deep != 0 && i >= deep-1 {
			break
		}

		err = errors.Unwrap(err); if err != nil {
			i++
			txt = fmt.Sprintf(txtFmt, prfix, i, err)
			color.HiMagenta(txt)
			// logger.Print(txt)
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

	var sw int
	flag.IntVar(&sw, "sw", 0, "0,1,2")
	flag.Parse()
	fmt.Printf("sw=%v\n", sw)

	m := "CSCO,BB.TO,BRK-A,BTC-USD"
	l1 := "AAPL, BCE.TO,BCE, CM, CM.TO, ENB, ENB.TO, AVGO, A, T, V, META,X, CAD=X, CADUSD=X, INTC, IBM, SHOP.TO, FCAU, GRMN"
	l2 := "GM, F, AMZN, CSCO, DELL, GE, GOOG, LOGI, MSFT, NVDA, QCOM, TSLA, ZM, GARM, YNDX"

	switch sw {
		case 1:  one(appCtx)
		case 2:  many(appCtx, m)
		case 3:  many(appCtx, l1)
		case 4:  many(appCtx, l2)
		default: one(appCtx); many(appCtx, m)
	}
}

func many(appCtx applicationContext, tickers string) {
	println("many " + tickers)
	numCommas := strings.Count(tickers, ",")
	res, err := appCtx.getStocksFullPricCSV(tickers); if err != nil {
		unwrapErrors("getStocksFullPricCSV", appCtx.logger, err, 1)
		return
	}
	// fmt.Printf("Full price CSV: %s", txt)
	res = strings.TrimSuffix(res, ",")
	lines := strings.Split(res, ",")
	slices.Sort(lines)

	for _, line := range lines {
		if strings.Contains(line, "+" ) {
			color.Green(line)
		} else {
			color.Red(line)
		}
	}

	color.Cyan("Total=%d/%d", len(lines), numCommas+1)
}

func one(appCtx applicationContext) {
	println("one")
	txt, err := appCtx.getSingleStockPriceTxt(); if err != nil {
		unwrapErrors("getSingleStockPriceTxt", appCtx.logger, err, 1)
	}
	fmt.Printf("Stock Price: %s\n", txt)
}

func (appCtx *applicationContext) getStocksFullPricCSV(tickers string) (string, error){
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

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

func main() {
	
	appCtx := applicationContext {
		logger: log.New(os.Stdout, "", log.Ltime),
	}

	txt, err := appCtx.getSingleStockPriceTxt(); if err != nil {
		unwrapErrors("getSingleStockPriceTxt", appCtx.logger, err, 1)
		return 
	}

	fmt.Printf("Stock Price: %s", txt)
}

type GetSingleStockPriceReqResp struct {
	TickerPrice string
	Error string
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

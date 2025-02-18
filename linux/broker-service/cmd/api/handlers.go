package main

import (
	"broker/event"
	"broker/trace"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/rpc"
	"time"

	shared "github.com/tsizism/semplicita/linux/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func setupHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

/*{
    "action": "auth",
    "auth": {
	    "email": "admin@example.com",
	    "password": "verysecret"
    }
}*/

type RequestPayload struct {
	Action string       `json:"action"`
	Auth   AuthPayload  `json:"auth,omitempty"`
	Trace  TracePayload `json:"trace,omitempty"`
	Mail   MailPayload  `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TracePayload struct {
	Src  string `json:"src"`
	Via  string `json:"via"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (appCtx *applicationContext) handleSubmission(w http.ResponseWriter, r *http.Request) {
	setupHeader(w)
	appCtx.logger.Printf("Hit broker handleSubmission http.Request=%+v\n", r)

	var requestPayload RequestPayload

	err := shared.ReadJSON(w, r, &requestPayload)

	if err != nil {
		appCtx.logger.Printf("handleSubmission: Failed to read payload %s", err)
		shared.ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		appCtx.authenticate(w, requestPayload.Auth)
	case "mail":
		appCtx.sendMail(w, requestPayload.Mail)
	case "trace":
		requestPayload.Trace.Via = "RPC"
		// appCtx.logEventRPC(w, requestPayload.Trace)
		//appCtx.getSingleStockPrice(w)
		// appCtx.getStocksPrice(w)
		appCtx.getStocksFullPricCSV(w)
	case "traceMq":
		requestPayload.Trace.Via = "MQ"
		appCtx.logEventMQ(w, requestPayload.Trace)
	case "traceEvent":
		appCtx.traceEvent(w, requestPayload.Trace)
	default:
		appCtx.logger.Printf("handleSubmission: bad request - unknown action= '%s'", requestPayload.Action)
		shared.ErrorJSON(w, errors.New("bad request"), http.StatusBadRequest)
	}
}

type RPCPayload struct {
	Src  string
	Via  string
	Data string
}

func (appCtx *applicationContext) logEventRPC(w http.ResponseWriter, tracePayload TracePayload) {
	appCtx.logger.Printf("logEventRPC: payload=%+v", tracePayload)
	client, err := rpc.Dial("tcp", "trace-service:5001")

	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Src:  tracePayload.Src,
		Via:  tracePayload.Via,
		Data: tracePayload.Data,
	}

	var result string
	err = client.Call("RPCServer.InsertTraceEvent", rpcPayload, &result)

	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	payload := shared.JsonResponse{
		Error:   false,
		Message: "Logged via RPC",
	}

	shared.WriteJSON(w, http.StatusAccepted, payload)
}

func (appCtx *applicationContext) getSingleStockPriceTxt(w http.ResponseWriter) {
	ticker := "AAPL"
	appCtx.logger.Printf("getSingleStockPriceTxt: payload=%+v", ticker)
	client, err := rpc.Dial("tcp", "localhost:5003") // faas-service:5003

	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	type GetSingleStockPriceReqResp struct {
		TickerPrice string
		Error string
	}

	var rpcPayload,result GetSingleStockPriceReqResp
	rpcPayload.TickerPrice = ticker
	err = client.Call("RPCServer.GetSingleStockPriceTxt", rpcPayload, &result)
	if err != nil {
		appCtx.logger.Printf("GetSingleStockPriceTxt: Failed to call RPC %s", err)
		shared.ErrorJSON(w, err)
		return
	}

	appCtx.logger.Printf("GetSingleStockPriceTxt for%s: result=%+v", ticker, result)

	payload := shared.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Stock Price: %f", result.TickerPrice),
	}

	shared.WriteJSON(w, http.StatusAccepted, payload)
}

func (appCtx *applicationContext) getStocksFullPricCSV(w http.ResponseWriter) {
	tickers := "BCE.TO,BCE, CM, CM.TO"
	appCtx.logger.Printf("getStocksFullPriceCSV: payload=%+v", tickers)
	client, err := rpc.Dial("tcp", "localhost:5003") // faas-service:5003

	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	type GetStocksPriceReqRespCSV struct {
		TickerPriceCSV string
		Error string
	}

	var result, rpcPayload GetStocksPriceReqRespCSV
	rpcPayload.TickerPriceCSV = tickers
	
	err = client.Call("RPCServer.GetStocksFullPriceCSV", rpcPayload, &result);	if err != nil {
		appCtx.logger.Printf("GetStocksFullPrice: Failed to call RPC %s", err)
		shared.ErrorJSON(w, err)
		return
	}

	appCtx.logger.Printf("GetStocksFullPriceCSV for %s: result=%+v", tickers, result)

	payload := shared.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("TickersPriceCSV: %s", result.TickerPriceCSV),
	}

	shared.WriteJSON(w, http.StatusAccepted, payload)
}

func (appCtx *applicationContext) getStocksPriceCSV(w http.ResponseWriter) {
	tickers := "BCE.TO,BCE"
	appCtx.logger.Printf("getStocksPriceCSV: payload=%+v", tickers)
	client, err := rpc.Dial("tcp", "localhost:5003") // faas-service:5003

	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	type GetStocksPriceReqRespCSV struct {
		TickerPriceCSV string
		Error string
	}

	var result, rpcPayload GetStocksPriceReqRespCSV
	rpcPayload.TickerPriceCSV = tickers
	
	err = client.Call("RPCServer.GetStocksPriceCSV", rpcPayload, &result);	if err != nil {
		appCtx.logger.Printf("GetStocksPriceCSV: Failed to call RPC %s", err)
		shared.ErrorJSON(w, err)
		return
	}

	appCtx.logger.Printf("GetStocksPriceCSV for%s: result=%+v", tickers, result)

	payload := shared.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("TickerPriceCSV: %s", result.TickerPriceCSV),
	}

	shared.WriteJSON(w, http.StatusAccepted, payload)
}

func (appCtx *applicationContext) traceGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := shared.ReadJSON(w, r, &requestPayload)

	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	appCtx.logger.Printf("traceGRPC: requestPayload=%+v", requestPayload)

	gprsCh, err := grpc.NewClient(fmt.Sprintf("trace-service:%d", appCtx.cfg.gRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials())) //, grpc.WithBlock())

	if err != nil {
		appCtx.logger.Println("traceGRPC: Failed NewClient")
		shared.ErrorJSON(w, err)
		return
	}

	defer gprsCh.Close()

	serviceClient := trace.NewTraceServiceClient(gprsCh)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	resp, err := serviceClient.TraceEvent(ctx, &trace.TraceRequest{
		TraceEntry: &trace.Trace{Src: requestPayload.Trace.Src, Data: requestPayload.Trace.Data},
	})

	if err != nil {
		appCtx.logger.Println("traceGRPC: Failed TraceEvent")
		shared.ErrorJSON(w, err)
		return
	}

	appCtx.logger.Printf("traceGRPC: TraceEvent=%+v", resp)

	payload := shared.JsonResponse{
		Error:   false,
		Message: "Logged trace via gRPC",
	}

	shared.WriteJSON(w, http.StatusAccepted, payload)
}

func (appCtx *applicationContext) logEventMQ(w http.ResponseWriter, t TracePayload) {
	appCtx.logger.Printf("logEventMQ: payload=%+v", t)
	err := appCtx.pushToQueue(t.Src, t.Data)

	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	payload := shared.JsonResponse{
		Error:   false,
		Message: "Logged via RabbitMQ",
	}

	shared.WriteJSON(w, http.StatusAccepted, payload)
}

func (appCtx *applicationContext) traceEvent(w http.ResponseWriter, t TracePayload) {
	appCtx.logger.Printf("traceEvent: payload=%+v", t)

	traceServiceEndpoint := "http://trace-service/trace"

	jsonData, _ := json.MarshalIndent(t, "", "\t") // Marshal in prod

	request, err := http.NewRequest("POST", traceServiceEndpoint, bytes.NewBuffer(jsonData))

	if err != nil {
		appCtx.logger.Printf("traceEvent: Failed to create endpoint request %s", err)
		shared.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		appCtx.logger.Printf("traceEvent: Failed to call endpoint %s", err)
		shared.ErrorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		appCtx.logger.Print("traceEvent:failed calling trace service")
		shared.ErrorJSON(w, errors.New("failed calling trace service"))
		return
	}

	payload := shared.JsonResponse{
		Error:   false,
		Message: "Inserted",
	}

	shared.WriteJSON(w, http.StatusAccepted, payload)
}

func (appCtx *applicationContext) pushToQueue(name, msg string) error {
	appCtx.logger.Println("pushToQueue")
	emitter, err := event.NewEventEmitter(appCtx.connMQ, appCtx.logger)

	if err != nil {
		return err
	}

	payload := TracePayload{
		Src:  name,
		Via:  "MQ",
		Data: msg,
	}

	jsonAsBytes, _ := json.MarshalIndent(payload, "", "\t") // Marshal in prod

	err = emitter.Push(string(jsonAsBytes), "log.INFO")

	if err != nil {
		return err
	}

	return nil
}

func (appCtx *applicationContext) sendMail(w http.ResponseWriter, p MailPayload) {
	setupHeader(w)
	appCtx.logger.Printf("sendMail: payload=%+v", p)

	endpoint := "http://mail-service/send"

	jsonAsBytes, _ := json.MarshalIndent(p, "", "\t") // Marshal in prod

	// appCtx.logger.Printf("sendMail: jsonData=%+v", jsonData)

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonAsBytes))

	if err != nil {
		appCtx.logger.Printf("sendMail: Failed to create endpoint request %s", err)
		shared.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		appCtx.logger.Printf("sendMail: Failed to call endpoint %s", err)
		shared.ErrorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		appCtx.logger.Print("sendMail:failed calling mail service")
		shared.ErrorJSON(w, errors.New("failed calling mail service"))
		return
	}

	payload := shared.JsonResponse{
		Error:   false,
		Message: "Mail Sent",
	}

	shared.WriteJSON(w, http.StatusAccepted, payload)
}

func (appCtx *applicationContext) authenticate(w http.ResponseWriter, a AuthPayload) {
	appCtx.logger.Printf("authenticate: payload=%+v", a)

	// >docker exec -it  project-authentication-service-1  sh
	// nslookup authentication-service
	// docker exec -it  project-authentication-service-1 nslookup authentication-service
	// docker exec -it  project-broker-service-1 nslookup  broker-service
	// docker exec -it  project-broker-service-1 sh
	// wget http://authentication-service/authenticate
	// Connecting to authentication-service (172.18.0.2:80)
	// wget: server returned error: HTTP/1.1 400 Bad Request

	authServiceEndpoint := "http://authentication-service:8081/authenticate"

	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", authServiceEndpoint, bytes.NewBuffer(jsonData))

	if err != nil {
		appCtx.logger.Printf("authenticate: Failed to create endpoint request %s", err)
		shared.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		appCtx.logger.Printf("authenticate: Failed to call endpoint %s", err)
		shared.ErrorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		appCtx.logger.Print("authenticate: invalid credentials")
		shared.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		appCtx.logger.Print("authenticate: failed calling auth service")
		shared.ErrorJSON(w, errors.New("failed calling auth service"))
		return
	}

	var jsonAuthResponse shared.JsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonAuthResponse)

	if err != nil {
		appCtx.logger.Printf("authenticate: Failed to decode response %s", err)
		shared.ErrorJSON(w, err)
		return
	}

	if jsonAuthResponse.Error {
		appCtx.logger.Println("authenticate: StatusUnauthorized(401)")
		shared.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	payload := shared.JsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonAuthResponse,
	}

	shared.WriteJSON(w, http.StatusAccepted, payload)

}

func (appCtx *applicationContext) handleRoot(w http.ResponseWriter, r *http.Request) {
	setupHeader(w)
	// appCtx.logger.Printf("Hit broker http.Request=%+v\n", r)

	res, err := httputil.DumpRequest(r, true)

	if err != nil {
		appCtx.logger.Fatal(err)
	}

	appCtx.logger.Println()
	appCtx.logger.Println()
	appCtx.logger.Println()
	appCtx.logger.Println("rootHandler: DumpRequest")
	appCtx.logger.Print(string(res))

	appCtx.logger.Printf("User-Agent:%s\n", r.UserAgent())

	payload := shared.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	setupHeader(w)
	_ = shared.WriteJSON(w, http.StatusOK, payload)

	// out , _ := json.MarshalIndent(payload, "", "\t")
	// w.Header().Set("ContentType", "application/json")
	// w.WriteHeader(http.StatusAccepted)
	// w.Write(out)

	// curl
	// 2024/12/04 11:35:11 DumpRequest
	// 2024/12/04 11:35:11 GET / HTTP/1.1
	// Host: localhost:4000
	// User-Agent: Mozilla/5.0 (Windows NT; Windows NT 10.0; en-US) WindowsPowerShell/5.1.19041.5129

	// browser

	// 2024/12/04 11:35:48 DumpRequest
	// 2024/12/04 11:35:48 GET /favicon.ico HTTP/1.1
	// Host: localhost:4000
	// Accept: image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8
	// Accept-Encoding: gzip, deflate, br, zstd
	// Accept-Language: en-US,en;q=0.9,ru;q=0.8
	// Connection: keep-alive
	// Referer: http://localhost:4000/
	// Sec-Ch-Ua: "Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"
	// Sec-Ch-Ua-Mobile: ?0
	// Sec-Ch-Ua-Platform: "Windows"
	// Sec-Fetch-Dest: image
	// Sec-Fetch-Mode: no-cors
	// Sec-Fetch-Site: same-origin
	// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36

}

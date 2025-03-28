package fintechapi

// https://rapidapi.com/belchiorarkad-FqvHs2EDOtP/api/yh-finance-complete
// https://algotrading101.com/learn/yahoo-finance-api-guide/
// https://rapidapi.com/
// Basic $9.99 - 14,986 / Month

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// export YHFINCOMPLETE_APIKEY_FN=~/github/yhfincomplete_apikey.txt
// echo $YHFINCOMPLETE_APIKEY_FN
// $env:YHFINCOMPLETE_APIKEY_FN="C:\github\yhfincomplete_apikey.txt"
// $env:YHFINCOMPLETE_APIKEY_FN
// set YHFINCOMPLETE_APIKEY_FN="C:\github\yhfincomplete_apikey.txt"
// echo %BROKER_URL%

const (
	UrlDomain_YHFinanceCompleteAPI = "https://yh-finance-complete.p.rapidapi.com"
	ApiHost_YHFinanceCompleteAPI   = "yh-finance-complete.p.rapidapi.com"
	ApiHost_ENVVAR                 = "YHFINCOMPLETE_APIKEY_FN"
)

type YHFinanceCompleteAPI struct {
	// ctx 					context.Context
	// cancel 					context.CancelFunc
	httpClient             *http.Client
	logger                 *log.Logger
	urlDomain              string // "https://yh-finance-complete.p.rapidapi.com"
	apiKey                 string
	apiHost                string
	cacheFileNameFmt       string
	requestCache           map[string]*http.Request
	invalidTickers         map[string]int
	tickerCh               chan Ticker
	priceResponseCh        chan YfpriceResponse
	priceResponseErrCh     chan error
	stockFullPriceTickerCh chan Ticker
	stockFullPriceCh       chan YffullstockpriceResponse
	stockFullPriceErrCh    chan error
	sessionEndCh           chan bool
}

type Ticker string

func NewYHFinanceCompleteAPI(logger *log.Logger) YHFinanceCompleteAPI {
	fn := os.Getenv(ApiHost_ENVVAR)
	apiKey_YHFinanceCompleteAPI, err := os.ReadFile(fn)
	if err != nil {
		log.Fatalf("os.ReadFile error: %v defined by envvar %s", err, ApiHost_ENVVAR)
	}
	// println("--------------------->" + string(apiKey_YHFinanceCompleteAPI))

	// tickerCh, priceResponseCh := runGetSingleStockPrice()

	// Use context.WithTimeout for per-request timeouts!

	yHFinanceCompleteAPI := YHFinanceCompleteAPI{
		// ctx: 				ctx,
		// cancel: 			cancel,
		httpClient:       http.DefaultClient,
		urlDomain:        UrlDomain_YHFinanceCompleteAPI,
		apiHost:          ApiHost_YHFinanceCompleteAPI,
		apiKey:           strings.TrimSpace(string(apiKey_YHFinanceCompleteAPI)),
		logger:           log.New(os.Stdout, "", log.Ltime),
		cacheFileNameFmt: "%s.%s.json",
		requestCache:     make(map[string]*http.Request),
		invalidTickers:   make(map[string]int),
		sessionEndCh:     make(chan bool),
	}

	const httpReqTimoutMs = TIMOUT_HTTPREQ_MS
	yHFinanceCompleteAPI.httpClient.Timeout = time.Duration(httpReqTimoutMs) * time.Millisecond

	yHFinanceCompleteAPI.runGetSingleStockPrice()
	yHFinanceCompleteAPI.runGetSingleStockFullPrice()

	return yHFinanceCompleteAPI
}

func (api *YHFinanceCompleteAPI) runGetSingleStockPrice() {
	// url := "https://yh-finance-complete.p.rapidapi.com/yhprice?ticker=BCE.TO"
	api.tickerCh = make(chan Ticker)
	api.priceResponseCh = make(chan YfpriceResponse)
	api.priceResponseErrCh = make(chan error, 1)

	api.logger.Println("Kicking getSingleStockPrice goroutine")

	go func() {
		for ticker := range api.tickerCh {
			jsonResponse, err := api.getSingleStockPrice(string(ticker))
			if err != nil {
				api.priceResponseErrCh <- err
			} else {
				api.priceResponseCh <- jsonResponse
			}
		}
	}()
}

// GetSingleStockPrice retrieves the stock price for the given ticker symbol.
// It sends a request to the YHFinanceComplete API and decodes the response into a YfpriceResponse struct.
//
// Parameters:
//   - ticker: The stock ticker symbol for which the price is to be retrieved.
//
// Returns:
//   - YfpriceResponse: The response containing the stock price information.
//   - error: An error if the request fails or the response is invalid.
//
// Errors:
//   - Returns an error if the ticker is empty.
//   - Returns an error if there is an issue building the request.
//   - Returns an error if the HTTP request fails.
//   - Returns an error if the response cannot be decoded.
//   - Returns an error if the symbol in the response is empty.
func (api YHFinanceCompleteAPI) getSingleStockPrice(ticker string) (YfpriceResponse, error) {
	api.logger.Println("GetSingleStockPrice: ticker=", ticker)
	// url := "https://yh-finance-complete.p.rapidapi.com/yhprice?ticker=BCE.TO"

	var jsonResponse YfpriceResponse

	if ticker == "" {
		return jsonResponse, fmt.Errorf("ticker is empty")
	}

	queryParams := url.Values{"ticker": {ticker}}

	req, err := api.buildRequest("yhprice", queryParams)
	if err != nil {
		return jsonResponse, fmt.Errorf("buildRequest error: %w", err)
	}

	// url := fmt.Sprintf("%s/yhprice?ticker=%s", urlDomain, ticker)

	// client := http.DefaultClient
	// client.Timeout = 500 * time.Millisecond

	// var myHttpClient = &http.Client{
	// 	Timeout: time.Second * 10,
	// }

	// t := http.DefaultTransport.(*http.Transport).Clone()
	// t.MaxIdleConns = 100
	// t.MaxConnsPerHost = 100
	// t.MaxIdleConnsPerHost = 100
		
	// myHttpClient := &http.Client{
	//   Timeout:   10 * time.Second,
	//   Transport: t,
	// }	
	// resp, err := myHttpClient.Do(req)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return jsonResponse, fmt.Errorf("DefaultClient.Do error: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		return jsonResponse, fmt.Errorf("json.Decode error: %w", err)
	}

	if jsonResponse.Symbol == "" {
		return jsonResponse, fmt.Errorf("symbol in response is empty")
	}

	fmt.Printf("jresp=%+v\n", jsonResponse)

	return jsonResponse, nil
}

func (api *YHFinanceCompleteAPI) runGetSingleStockFullPrice() {
	api.stockFullPriceTickerCh = make(chan Ticker)
	api.stockFullPriceCh = make(chan YffullstockpriceResponse)
	api.stockFullPriceErrCh = make(chan error, 1)

	respCache := []YffullstockpriceResponse{}
	errCache := []error{}

	respMutex := sync.Mutex{}
	errMutex := sync.Mutex{}

	api.logger.Println("Kicking getSingleStockFullPrice goroutine")
	go func() {
		var resp YffullstockpriceResponse
		var err error
		for {
			respMutex.Lock()
			if len(respCache) > 0 {
				resp, respCache = respCache[0], respCache[1:] // pop from queue;pop from stack x, a = a[len(a)-1], a[:len(a)-1]
				// fmt.Printf("Sending resp %v\n", resp)
				println("********************Sending to resp stockFullPriceCh")
				api.stockFullPriceCh <- resp
				fmt.Printf("+++++++++++++++++++++++Resp sent to stockFullPriceCh for Symbol=%s\n", resp.Price.Symbol)
			} else {
			 	time.Sleep(10 * time.Millisecond)
			}
			respMutex.Unlock()

			if len(errCache) > 0 {
				errMutex.Lock()
				err, errCache = errCache[0], errCache[1:]
				errMutex.Unlock()
				api.stockFullPriceErrCh <- err
			}
		}
	}()

	go func() {
		respMutex2 := sync.Mutex{}
		for ticker := range api.stockFullPriceTickerCh {
			go func(ticker Ticker) {
				fmt.Printf("Processing ticker recived from stockFullPriceTickerCh=%s\n", ticker)
				respMutex2.Lock()
				resp, err := api.getSingleStockFullPrice(string(ticker))
				respMutex2.Unlock()

				if err != nil {
					fmt.Printf("Saving err %v\n", err)
					errMutex.Lock()
					errCache = append(errCache, err)
					errMutex.Unlock()
					fmt.Println("Error saved")

				} else {
					fmt.Printf("Saving resp %v\n", resp)
					respMutex.Lock()
					respCache = append(respCache, resp)
					fmt.Printf("Resp saved to respCache cache size %d\n", len(respCache))
					respMutex.Unlock()
				}
			}(ticker)
		}
	}()

	go func() {
		for range api.sessionEndCh {
			fmt.Printf("sessionEndCh --------------->: cache size=%d\n", len(respCache))

			respMutex.Lock()
			// close(api.stockFullPriceCh)
			respCache = nil
			errCache  = nil
			respMutex.Unlock()

			// for _, v := range respCache {
			// 	fmt.Printf("%v\n", v.Price.Symbol)
			// }

		}
	}()
}

// GetFullSingleStockPrice retrieves the full stock price information for a given stock symbol.
// It sends a request to the YH Finance Complete API and decodes the JSON response into a YffullstockpriceResponse struct.
//
// Parameters:
//   - symbol: The stock symbol for which to retrieve the price information.
//
// Returns:
//   - YffullstockpriceResponse: The response containing the stock price information.
//   - error: An error if the request fails or the response is invalid.
//
// The function performs the following steps:
//  1. Logs the request with the provided stock symbol.
//  2. Checks if the symbol is empty and returns an error if it is.
//  3. Builds the request with the given symbol as a query parameter.
//  4. Sends the request using the default HTTP client.
//  5. Decodes the JSON response into the YffullstockpriceResponse struct.
//  6. Checks if the symbol in the response is empty and returns an error if it is.
//  7. Returns the decoded response and any error encountered during the process.
func (api YHFinanceCompleteAPI) getSingleStockFullPrice(ticker string) (YffullstockpriceResponse, error) {
	api.logger.Println("getSingleStockFullPrice: ticker=", ticker)
	//  "https://yh-finance-complete.p.rapidapi.com/price?symbol=cm.to"

	var fullStockPrice YffullstockpriceResponse
	// jsonResponse := YffullstockpriceResponse{}
	// jsonResponse.Price.Symbol = ticker

	invalidTickerCnt, ok := api.invalidTickers[ticker]
	if ok { //&& invalidTickerCnt % 100 != 0
		invalidTickerCnt++
		api.invalidTickers[ticker] = invalidTickerCnt
		api.logger.Printf("getSingleStockFullPrice: '%s' ticker (from cache) identified as invalid cnt=%d, skipping retrival\n", ticker, invalidTickerCnt)
		api.logger.Printf("================Invalid ticker cache %v\n", api.invalidTickers)
		return fullStockPrice, fmt.Errorf("invalid ticker %s(from cache) %d", ticker, invalidTickerCnt)
	}

	if ticker == "" {
		return fullStockPrice, fmt.Errorf("ticker is empty")
	}

	queryParams := url.Values{"symbol": {ticker}}

	// ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second);	defer cancel()
	req, err := api.buildRequest("price", queryParams)
	if err != nil {
		api.logger.Printf("buildRequest error: %v", err)
		return fullStockPrice, fmt.Errorf("buildRequest error: %w", err)
		// return jsonResponse, nil
	}

	fmt.Printf("---------------------------->httpClient.Do: URL = %v\n", req.URL)

	// var myHttpClient = &http.Client{
	//  	Timeout: time.Second * 100,
	// }

	// t := http.DefaultTransport.(*http.Transport).Clone()
	// t.MaxIdleConns = 100
	// t.MaxConnsPerHost = 100
	// t.MaxIdleConnsPerHost = 100
		
	// myHttpClient := &http.Client{
	//   Timeout:  100 * time.Second,
	//   Transport: t,
	// }	
	// resp, err := myHttpClient.Do(req)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		fmt.Printf("httpClient Do error: %v\n", err)
		return fullStockPrice, fmt.Errorf("httpClient Do error: %w", err)
	}

	defer resp.Body.Close()

	dump, _ := httputil.DumpResponse(resp, true)
	fmt.Printf("----------------------->url=%s, %q\n", req.URL, dump)  // %q to safely escape a string and add quotes to it.


	if err := json.NewDecoder(resp.Body).Decode(&fullStockPrice); err != nil {
		// dump, _ := httputil.DumpResponse(resp, true)
		// fmt.Printf("----------------------->%q\n", dump)

		body, _ := io.ReadAll(resp.Body)

		if len(body) == 0 {
			api.invalidTickers[ticker] = 1
			err = fmt.Errorf("invalid ticker %s", ticker)
			api.logger.Printf("--------------->invalid ticker %s", ticker)

		} else {
			api.logger.Printf("json.Decode response body error %s", ticker)
			err = fmt.Errorf("json.Decode response body error: %w", err)
		}

		return fullStockPrice, err
	}

	if fullStockPrice.Price.Symbol == "" {
		return fullStockPrice, fmt.Errorf("symbol in response is empty")
	}

	// fmt.Printf("getSingleStockFullPrice resp=%v\n", fullStockPrice)
	api.logger.Printf("getSingleStockFullPrice resp for symbol=%s, ticker=%s\n", fullStockPrice.Price.Symbol, ticker)

	return fullStockPrice, nil
}

// buildRequest constructs an HTTP request for the YHFinanceCompleteAPI.
// It takes a sub-directory path and query parameters as inputs and returns
// an HTTP request pointer and an error if any occurs during the request creation.
//
// Parameters:
//   - subDir: A string representing the sub-directory path for the request.
//   - queryParams: A url.Values object containing the query parameters for the request.
//
// Returns:
//   - *http.Request: A pointer to the constructed HTTP request.
//   - error: An error object if an error occurs during the request creation.
func (api YHFinanceCompleteAPI) buildRequest(subDir string, queryParams url.Values) (*http.Request, error) { // context.CancelFunc
	// api.logger.Printf("buildRequest: subDir=%s %+v\n", subDir, queryParams)

	// fmt.Printf("YHFinanceCompleteAPI.go: yhfhistoricalDecode ticker=%s,sdate=%s,edate=%s\n", ticker, sdate, edate)

	// url := "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=TSLA&sdate=2025-02-10&edate=2025-02-11"
	// url := fmt.Sprintf("https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=%s&sdate=%s&edate=%s", ticker, sdate, edate)
	// url := fmt.Sprintf("https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=%s&sdate=%s&edate=%s", ticker, sdate, edate)

	var request *http.Request
	var exists bool
	// ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second) ///;	defer cancel()
	if request, exists = api.requestCache[subDir]; !exists {
		// api.logger.Printf("buildRequest: creating new request for =%s\n", subDir)
		var err error

		// ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second);	defer cancel()
		// request, err = http.NewRequestWithContext(context.Background(), http.MethodGet, "", nil)

		request, err = http.NewRequest("GET", "", nil)
		// api.logger.Printf("buildRequest w/ctx -------------->: request=%+v\n", request)
		if err != nil {
			return nil, fmt.Errorf("http.NewRequest error: %w", err) //, cancel
		}

		api.requestCache[subDir] = request
		// api.logger.Printf("buildRequest -------------->: requestCache=%+v\n", api.requestCache)

		// fmt.Printf("request=%+v\n", request)

		request.Header.Add("x-rapidapi-host", api.apiHost)
		request.Header.Add("x-rapidapi-key", api.apiKey)
		// api.logger.Printf("x-rapidapi-key='%s'\n", api.apiKey)
	// } else {
	// 	api.logger.Printf("buildRequest: pulled request from cache") //=%v\n", request)
	}

	requestUrl := fmt.Sprintf("%s/%s?", api.urlDomain, subDir)
	// api.logger.Printf("buildRequest: queryParams=%v str=%s\n", queryParams, queryParams.Encode())
	requestUrl += queryParams.Encode()

	// api.logger.Printf("buildRequest: requestUrl=%s\n", requestUrl)

	var err error
	// api.logger.Printf("buildRequest -------------->: request=%+v\n", request)
	request.URL, err = url.Parse(requestUrl)

	if err != nil {
		return nil, fmt.Errorf("url.Parse error: %w", err) //, cancel
	}
	return request, nil //, cancel
}

// GetHistoricalWithUnmarshal retrieves historical financial data for a given ticker symbol
// between the specified start date (sdate) and end date (edate). It first attempts to read
// the data from a cached file. If the file does not exist, it calls an external API to fetch
// the data, saves the response to a cache file, and then unmarshals the JSON response into
// a slice of YfhistoricalResponse structs.
//
// Parameters:
//   - ticker: The stock ticker symbol (e.g., "TSLA").
//   - sdate: The start date for the historical data in the format "YYYY-MM-DD".
//   - edate: The end date for the historical data in the format "YYYY-MM-DD".
//
// Returns:
//   - A slice of YfhistoricalResponse structs containing the historical financial data.
//   - An error if any issues occur during the process of reading the file, making the API
//     request, or unmarshaling the JSON response.
func (api YHFinanceCompleteAPI) GetHistoricalWithUnmarshal(ticker, sdate, edate string) ([]YfhistoricalResponse, error) {
	// [map[adjclose:350.7300109863281 close:350.7300109863281 date:2025-02-10T14:30:00.000Z high:362.70001220703125 low:350.510009765625 open:356.2099914550781 volume:7.75149e+07]] 1 1
	// url := "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=TSLA&sdate=2024-01-10&edate=2024-02-16"

	var jsonMapArr []YfhistoricalResponse
	api.logger.Printf("YHFinanceCompleteAPI.go: yhfhistoricalUnmarshal ticker=%s,sdate=%s,edate=%s\n", ticker, sdate, edate)

	fname := fmt.Sprintf(api.cacheFileNameFmt, ticker, sdate)
	bodyTxt := "" // string

	buf, err := os.ReadFile(fname)
	if err != nil {
		println("File does not exist. Calling API", fname)

		queryParams := url.Values{"ticker": {ticker}, "sdate": {sdate}, "edate": {edate}}

		req, err := api.buildRequest("yhfhistorical", queryParams)
		if err != nil {
			return jsonMapArr, fmt.Errorf("buildRequest error: %w", err)
		}

		resp, err := api.httpClient.Do(req)
		if err != nil {
			api.logger.Println("DefaultClient.Do error==============>", err)
			return jsonMapArr, fmt.Errorf("DefaultClient.Do error: %w", err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return jsonMapArr, fmt.Errorf("io.ReadAll error: %w", err)
		}

		bodyTxt = string(body)
		api.logger.Println("bodyTxt------------->", bodyTxt)

		f, err := os.Create(fname)
		if err != nil {
			return jsonMapArr, fmt.Errorf("os.Create error: %w", err)
		}

		_, err = f.WriteString(bodyTxt)
		if err != nil {
			return jsonMapArr, fmt.Errorf("f.WriteString error: %w", err)
		}
		f.Sync()
		f.Close()
	} else {
		println("Read JSON data from file", fname)
		bodyTxt = string(buf)
	}

	err = json.Unmarshal([]byte(bodyTxt), &jsonMapArr)
	if err != nil {
		api.logger.Println("json.Unmarshal error==============>", err)
		return jsonMapArr, fmt.Errorf("json.Unmarshal error: %w", err)
	}

	fmt.Println(jsonMapArr, len(jsonMapArr), cap(jsonMapArr))
	fmt.Printf("%+v\n", jsonMapArr)

	return jsonMapArr, nil
}

func (api YHFinanceCompleteAPI) GetHistoricalWitDecode(ticker, sdate, edate string) ([]YfhistoricalResponse, error) {
	// [map[adjclose:350.7300109863281 close:350.7300109863281 date:2025-02-10T14:30:00.000Z high:362.70001220703125 low:350.510009765625 open:356.2099914550781 volume:7.75149e+07]] 1 1
	var jsonMapArr []YfhistoricalResponse

	// url := "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=TSLA&sdate=2025-02-10&edate=2025-02-11"

	queryParams := url.Values{
		"ticker": {ticker},
		"sdate":  {sdate},
		"edate":  {edate},
	}

	req, err := api.buildRequest("yhfhistorical", queryParams)
	if err != nil {
		return jsonMapArr, fmt.Errorf("buildRequest error: %w", err)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return jsonMapArr, fmt.Errorf("DefaultClient.Do error: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&jsonMapArr); err != nil {
		return jsonMapArr, fmt.Errorf("json.Decode error: %w", err)
	}

	fmt.Println(jsonMapArr, len(jsonMapArr), cap(jsonMapArr))
	fmt.Printf("resp=%+v\n", jsonMapArr)

	return jsonMapArr, nil
}

// GetStockSummaryDetail retrieves the stock summary details for a given ticker symbol.
// It sends a request to the YH Finance Complete API and decodes the response into a YfResponse struct.
//
// Parameters:
//   - ticker: The stock ticker symbol for which to retrieve the summary details.
//
// Returns:
//   - YfResponse: The response containing the stock summary details.
//   - error: An error if the request fails or the response is invalid.
//
// Errors:
//   - Returns an error if the ticker is empty.
//   - Returns an error if there is an issue building the request.
//   - Returns an error if there is an issue with the HTTP request.
//   - Returns an error if there is an issue decoding the JSON response.
//   - Returns an error if the symbol in the response is empty.
func (api YHFinanceCompleteAPI) GetStockSummaryDetail(ticker string) (YfResponse, error) {
	api.logger.Println("GetSummaryDetail: ticker=", ticker)
	// url := "https://yh-finance-complete.p.rapidapi.com/yhsummary?ticker=BCE.TO"

	var jsonResponse YfResponse

	if ticker == "" {
		return jsonResponse, fmt.Errorf("ticker is empty")
	}

	queryParams := url.Values{"ticker": {ticker}}

	req, err := api.buildRequest("yhf", queryParams)
	if err != nil {
		return jsonResponse, fmt.Errorf("buildRequest error: %w", err)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return jsonResponse, fmt.Errorf("DefaultClient.Do error: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		return jsonResponse, fmt.Errorf("json.Decode error: %w", err)
	}

	if jsonResponse.Price.Symbol == "" {
		return jsonResponse, fmt.Errorf("symbol in response is empty")
	}

	fmt.Printf("GetStockSummaryDetail resp=%+v\n", jsonResponse)

	return jsonResponse, nil
}

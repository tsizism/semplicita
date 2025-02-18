package fintechapi

// https://rapidapi.com/belchiorarkad-FqvHs2EDOtP/api/yh-finance-complete
// https://algotrading101.com/learn/yahoo-finance-api-guide/

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	UrlDomain_YHFinanceCompleteAPI = "https://yh-finance-complete.p.rapidapi.com"
	ApiKey_YHFinanceCompleteAPI    = "9b405718ddmsh954d4191ebcf658p148c17jsn58521162b938"
	ApiHost_YHFinanceCompleteAPI   = "yh-finance-complete.p.rapidapi.com"
)

type YHFinanceCompleteAPI struct {
	logger           *log.Logger
	urlDomain        string // "https://yh-finance-complete.p.rapidapi.com"
	apiKey           string
	apiHost          string
	cacheFileNameFmt string
}

func NewYHFinanceCompleteAPI(logger *log.Logger) YHFinanceCompleteAPI {
	return YHFinanceCompleteAPI{
		urlDomain:        UrlDomain_YHFinanceCompleteAPI,
		apiHost:          ApiHost_YHFinanceCompleteAPI,
		apiKey:           ApiKey_YHFinanceCompleteAPI,
		logger:           log.New(os.Stdout, "", log.Ltime),
		cacheFileNameFmt: "%s.%s.json",
	}
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
func (api YHFinanceCompleteAPI) buildRequest(subDir string, queryParams url.Values) (*http.Request, error) {
	api.logger.Printf("buildRequest: subDir=%s %+v\n", subDir, queryParams)

	// fmt.Printf("YHFinanceCompleteAPI.go: yhfhistoricalDecode ticker=%s,sdate=%s,edate=%s\n", ticker, sdate, edate)

	// url := "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=TSLA&sdate=2025-02-10&edate=2025-02-11"
	// url := fmt.Sprintf("https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=%s&sdate=%s&edate=%s", ticker, sdate, edate)
	// url := fmt.Sprintf("https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=%s&sdate=%s&edate=%s", ticker, sdate, edate)

	requestUrl := fmt.Sprintf("%s/%s?", api.urlDomain, subDir)
	requestUrl += queryParams.Encode()

	api.logger.Printf("buildRequest: requestUrl=%s\n", requestUrl)

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest error: %w", err)
	}

	fmt.Printf("req=%+v\n", req)

	// req.Header.Add("x-rapidapi-host", "yh-finance-complete.p.rapidapi.com")
	// req.Header.Add("x-rapidapi-key", "9b405718ddmsh954d4191ebcf658p148c17jsn58521162b938")

	req.Header.Add("x-rapidapi-host", api.apiHost)
	req.Header.Add("x-rapidapi-key", api.apiKey)

	return req, nil
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

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			api.logger.Println("DefaultClient.Do error==============>", err)
			return jsonMapArr, fmt.Errorf("DefaultClient.Do error: %w", err)
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return jsonMapArr, fmt.Errorf("DefaultClient.Do error: %w", err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&jsonMapArr); err != nil {
		return jsonMapArr, fmt.Errorf("json.Decode error: %w", err)
	}

	fmt.Println(jsonMapArr, len(jsonMapArr), cap(jsonMapArr))
	fmt.Printf("resp=%+v\n", jsonMapArr)

	return jsonMapArr, nil
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
func (api YHFinanceCompleteAPI) GetSingleStockPrice(ticker string) (YfpriceResponse, error) {
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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return jsonResponse, fmt.Errorf("DefaultClient.Do error: %w", err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		return jsonResponse, fmt.Errorf("json.Decode error: %w", err)
	}

	if jsonResponse.Symbol == "" {
		return jsonResponse, fmt.Errorf("symbol in response is empty")
	}

	fmt.Printf("resp=%+v\n", jsonResponse)

	return jsonResponse, nil
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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return jsonResponse, fmt.Errorf("DefaultClient.Do error: %w", err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		return jsonResponse, fmt.Errorf("json.Decode error: %w", err)
	}

	if jsonResponse.Price.Symbol == "" {
		return jsonResponse, fmt.Errorf("symbol in response is empty")
	}

	fmt.Printf("GetStockSummaryDetail resp=%+v\n", jsonResponse)

	return jsonResponse, nil
}

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
	"time"
)

const (
	UrlDomain_YHFinanceCompleteAPI = "https://yh-finance-complete.p.rapidapi.com"
	ApiKey_YHFinanceCompleteAPI    = "9b405718ddmsh954d4191ebcf658p148c17jsn58521162b938"
	ApiHost_YHFinanceCompleteAPI   = "yh-finance-complete.p.rapidapi.com"
)

type YfhistoricalResponse struct {
	Date     string  `json:"date"`
	Adjclose float32 `json:"adjclose"`
	Close    float32 `json:"close"`
	High     float32 `json:"high"`
	Low      float32 `json:"low"`
	Open     float32 `json:"open"`
	Volume   float32 `json:"volume"`
}

// url := "https://yh-finance-complete.p.rapidapi.com/yhprice?ticker=BCE.TO"

type YHFinanceCompleteAPI struct {
	logger           *log.Logger
	urlDomain        string // "https://yh-finance-complete.p.rapidapi.com"
	apiKey           string
	apiHost          string
	cacheFileNameFmt string
}

func weekStartDate(date time.Time) time.Time {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}

// NewYHFinanceCompleteAPI creates a new instance of YHFinanceCompleteAPI with the provided logger.
// It initializes the API with the necessary URL domain, API host, and API key.
// The logger is set to output to standard output with time-stamped log entries.
//
// Parameters:
//   - logger: A pointer to a log.Logger instance for logging purposes.
//
// Returns:
//   - YHFinanceCompleteAPI: A new instance of YHFinanceCompleteAPI configured with the provided logger and default settings.
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

func (api YHFinanceCompleteAPI) GetHistoricalWithUnmarshal(ticker, sdate, edate string) ([]YfhistoricalResponse, error) {
	// [map[adjclose:350.7300109863281 close:350.7300109863281 date:2025-02-10T14:30:00.000Z high:362.70001220703125 low:350.510009765625 open:356.2099914550781 volume:7.75149e+07]] 1 1

	// url := "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=TSLA&sdate=2024-01-10&edate=2024-02-16"
	// url := "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?edate=2025-02-01&sdate=2025-01-27&ticker=TSLA"

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

	// var resp YfhistoricalResponse
	// var jsonMapArr2 []YfhistoricalResponse
	// if err := json.NewDecoder(res.Body).Decode(&jsonMapArr2); err != nil {
	// 	fmt.Println("json.Decode error:", err)
	// }

	// fmt.Printf("resp=%+v\n",jsonMapArr2)

	return jsonMapArr, nil
}

func (api YHFinanceCompleteAPI) YHFHistoricalWithDecode(ticker, sdate, edate string) ([]YfhistoricalResponse, error) {
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

	fmt.Printf("resp=%+v\n", jsonMapArr)

	return jsonMapArr, nil
}

// func (api YHFinanceCompleteAPI) GetStockPrice(ticker string) (float64, error) {
// 	api.logger.Println("GetStockPrice: ticker=", ticker)

// 	// url := "https://yh-finance-complete.p.rapidapi.com/yhprice?ticker=BCE.TO"

// 	pararms := map[string]string{"ticker": ticker}

// 	req, err := api.buildRequest("yhprice", pararms); if err != nil {
// 		api.logger.Println("buildRequest error=", err)
// 		return 0, fmt.Errorf("buildRequest error: %w", err)
// 	}

// 	// url := fmt.Sprintf("%s/yhprice?ticker=%s", urlDomain, ticker)

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("http.NewRequest error: %w", err)
// 	}

// 	req.Header.Add("x-rapidapi-key", apiKey)
// 	req.Header.Add("x-rapidapi-host", apiHost)

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("DefaultClient.Do error: %w", err)
// 	}

// 	defer res.Body.Close()

// 	var quoteResponse QuoteResponse
// 	if err := json.NewDecoder(res.Body).Decode(&quoteResponse); err != nil {
// 		return nil, err
// 	}

// 	return &quoteResponse, nil
// }

// func buildRequestForYhfhistorical(ticker, sdate, edate string) (*http.Request, error) {
// 	fmt.Printf("YHFinanceCompleteAPI.go: yhfhistoricalDecode ticker=%s,sdate=%s,edate=%s\n", ticker, sdate, edate)

// 	// url := "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=TSLA&sdate=2025-02-10&edate=2025-02-11"
// 	// url := fmt.Sprintf("https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=%s&sdate=%s&edate=%s", ticker, sdate, edate)
// 	url := fmt.Sprintf("https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=%s&sdate=%s&edate=%s", ticker, sdate, edate)
// 	println(url)

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("http.NewRequest error: %w", err)
// 	}

// 	fmt.Printf("req=%+v\n", req)

// 	req.Header.Add("x-rapidapi-key", "9b405718ddmsh954d4191ebcf658p148c17jsn58521162b938")
// 	req.Header.Add("x-rapidapi-host", "yh-finance-complete.p.rapidapi.com")

// 	return req, nil
// }

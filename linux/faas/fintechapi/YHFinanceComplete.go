package fintechapi

// https://rapidapi.com/belchiorarkad-FqvHs2EDOtP/api/yh-finance-complete
// https://algotrading101.com/learn/yahoo-finance-api-guide/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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

func buildRequest(ticker, sdate, edate string) (*http.Request, error) {
	fmt.Printf("YHFinanceComplete.go: yhfhistoricalDecode ticker=%s,sdate=%s,edate=%s\n", ticker, sdate, edate)

	// url := "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=TSLA&sdate=2025-02-10&edate=2025-02-11"
	// url := fmt.Sprintf("https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=%s&sdate=%s&edate=%s", ticker, sdate, edate)
	url := fmt.Sprintf("https://yh-finance-complete.p.rapidapi.com/yhfhistorical?ticker=%s&sdate=%s&edate=%s", ticker, sdate, edate)
	println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest error: %w", err)
	}

	fmt.Printf("req=%+v\n", req)

	req.Header.Add("x-rapidapi-key", "9b405718ddmsh954d4191ebcf658p148c17jsn58521162b938")
	req.Header.Add("x-rapidapi-host", "yh-finance-complete.p.rapidapi.com")

	return req, nil
}

func yhfhistoricalUnmarshal(ticker, sdate, edate string) ([]YfhistoricalResponse, error) {
	// [map[adjclose:350.7300109863281 close:350.7300109863281 date:2025-02-10T14:30:00.000Z high:362.70001220703125 low:350.510009765625 open:356.2099914550781 volume:7.75149e+07]] 1 1
	var jsonMapArr []YfhistoricalResponse
	fmt.Printf("YHFinanceComplete.go: yhfhistoricalUnmarshal ticker=%s,sdate=%s,edate=%s\n", ticker, sdate, edate)

	fname := sdate + ".json"
	bodyTxt := "" // string

	buf, err := os.ReadFile(fname)
	if err != nil {
		println("File does not exist. Clling API", fname)

		req, err := buildRequest(ticker, sdate, edate)
		if err != nil {
			return jsonMapArr, fmt.Errorf("buildRequest error: %w", err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return jsonMapArr, fmt.Errorf("DefaultClient.Do error: %w", err)
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return jsonMapArr, fmt.Errorf("io.ReadAll error: %w", err)
		}

		bodyTxt = string(body)
		// fmt.Println("body=", bodyTxt)

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

func yhfhistoricalDecode(ticker, sdate, edate string) ([]YfhistoricalResponse, error) {
	// [map[adjclose:350.7300109863281 close:350.7300109863281 date:2025-02-10T14:30:00.000Z high:362.70001220703125 low:350.510009765625 open:356.2099914550781 volume:7.75149e+07]] 1 1
	var jsonMapArr []YfhistoricalResponse

	req, err := buildRequest(ticker, sdate, edate)
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

func weekStartDate(date time.Time) time.Time {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}

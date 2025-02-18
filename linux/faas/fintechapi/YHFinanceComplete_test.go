package fintechapi

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// go test .\fintechapi\. -v
//go test YHFinanceComplete_test.go YHFinanceComplete.go .\YHFinanceResponse.go   -v 

// func buildRequest(t *testing.T) {
// 	ticker, sdate, edate := "TSLA", "2025-02-10", "2025-02-10"
// 	req, err := buildRequestForYhfhistorical(ticker, sdate, edate)
// 	assert.NotNil(t, req)
// 	assert.Nil(t, err)

//		// req, err = buildRequest("", "", "")
//		// assert.NotNil(t, req)
//		// assert.Nil(t, err)
//	}

// util_weekStartDate calculates the start date of the week for a given date.
// The start of the week is considered to be Monday.
//
// Parameters:
//   - date: The input date for which the start of the week is to be calculated.
//
// Returns:
//   - time.Time: The start date of the week (Monday) for the given date.
func util_weekStartDate(date time.Time) time.Time {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}

func util_getTestWeekStartAndEndDates() (string, string) {
	// const daysBack = 2
	// sdate := time.Now().AddDate(0,0,(-1 * (daysBack+3))).Format("2006-01-02")
	// edate := time.Now().AddDate(0,0,-2).Format("2006-01-02")
	// fmt.Println("Test_yhfhistorical()", daysBack, sdate, edate)

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	weekStart := util_weekStartDate(firstOfMonth)
	weekEnd := weekStart.AddDate(0, 0, 5)

	sdate := weekStart.Format("2006-01-02")
	edate := weekEnd.Format("2006-01-02")

	fmt.Printf("util_getTestWeekStartAndEndDates ---> firstOfMonth=%s weekStart=%s weekEnd=%s\n", firstOfMonth.Format("2006-01-02 Mon"), sdate, edate)

	return sdate, edate
}
func TestWeekStartDate(t *testing.T) {
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 7; i++ {
		weekStart := util_weekStartDate(date)
		// fmt.Printf("%s %s\n", date.Format("2006-01-02 Mon"), weekStart.Format("2006-01-02 Mon"))
		assert.NotNil(t, weekStart)
		assert.Equal(t, time.Monday, weekStart.Weekday())
		date = date.Add(24 * time.Hour)
	}
	date = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 7; i++ {
		weekStart := util_weekStartDate(date)
		// fmt.Printf("%s %s\n", date.Format("2006-01-02 Mon"), weekStart.Format("2006-01-02 Mon"))
		assert.NotNil(t, weekStart)
		assert.Equal(t, time.Monday, weekStart.Weekday())
		date = date.Add(24 * time.Hour)
	}
}

func Test_GetHistoricalWithUnmarshal(t *testing.T) {
	yahooApi := NewYHFinanceCompleteAPI(log.New(os.Stdout, "", log.Ltime))
	ticker := "BCE.TO"
	sdate, edate := util_getTestWeekStartAndEndDates()
	fname := fmt.Sprintf(yahooApi.cacheFileNameFmt, ticker, sdate)

	t.Run("GetHistoricalWithUnmarshal BCE.TO Valid request without existing file", func(t *testing.T) {
		// arrange
		os.Remove(fname)

		// act
		jsonMapArr, err := yahooApi.GetHistoricalWithUnmarshal(ticker, sdate, edate)

		// assert
		assert.Nil(t, err)
		assert.NotEmpty(t, jsonMapArr)
		assert.Equal(t, 5, len(jsonMapArr)) // 5 business days in a week
		assert.Equal(t, "2025-01-27T14:30:00.000Z", jsonMapArr[0].Date)
		assert.Equal(t, float32(33.89), jsonMapArr[0].Low)
	})

	t.Run("GetHistoricalWithUnmarshal BCE.TO Valid request with existing file", func(t *testing.T) {
		// arrange
		expectedContent := `[{"date":"2025-01-27T14:30:00.000Z","high":34.79999923706055,"volume":3596500,"open":33.88999938964844,"low":33.88999938964844,"close":34.560001373291016,"adjclose":34.560001373291016},{"date":"2025-01-28T14:30:00.000Z","high":35.130001068115234,"volume":3729500,"open":34.970001220703125,"low":34.369998931884766,"close":34.40999984741211,"adjclose":34.40999984741211},{"date":"2025-01-29T14:30:00.000Z","high":34.65999984741211,"volume":1919000,"open":34.33000183105469,"low":34.150001525878906,"close":34.209999084472656,"adjclose":34.209999084472656},{"date":"2025-01-30T14:30:00.000Z","high":34.86000061035156,"volume":2790600,"open":34.290000915527344,"low":34.029998779296875,"close":34.61000061035156,"adjclose":34.61000061035156},{"date":"2025-01-31T14:30:00.000Z","high":34.970001220703125,"volume":3340000,"open":34.59000015258789,"low":34.40999984741211,"close":34.61000061035156,"adjclose":34.61000061035156}]`
		err := os.WriteFile(fname, []byte(expectedContent), 0644)
		assert.Nil(t, err)
		defer os.Remove(fname)

		// act
		jsonMapArr, err := yahooApi.GetHistoricalWithUnmarshal(ticker, sdate, edate)

		// assert
		assert.Nil(t, err)
		assert.NotEmpty(t, jsonMapArr)
		assert.Equal(t, 5, len(jsonMapArr))
		assert.Equal(t, "2025-01-27T14:30:00.000Z", jsonMapArr[0].Date)
		assert.Equal(t, float32(34.8), jsonMapArr[0].High)
	})

	t.Run("GetHistoricalWithUnmarshal Invalid request with incorrect dates", func(t *testing.T) {
		// arrange
		sdate, edate := "1900-01-01", "1900-01-01"

		// act
		jsonMapArr, err := yahooApi.GetHistoricalWithUnmarshal(ticker, sdate, edate)

		// assert
		assert.NotNil(t, err)
		assert.Empty(t, jsonMapArr)
	})
}

func Test_GetHistoricalWithDecode(t *testing.T) {
	// arrange
	yahooApi := NewYHFinanceCompleteAPI(log.New(os.Stdout, "", log.Ltime))
	ticker := "BCE.TO"
	sdate, edate := util_getTestWeekStartAndEndDates()
	// fname := fmt.Sprintf(api.cacheFileNameFmt, ticker, sdate)

	t.Run("GetHistoricalWithDecode BCE.TO Valid request without existing file", func(t *testing.T) {
		// act
		jsonMapArr, err := yahooApi.GetHistoricalWithUnmarshal(ticker, sdate, edate)

		// assert
		assert.Nil(t, err)
		assert.NotEmpty(t, jsonMapArr)
		assert.Equal(t, 5, len(jsonMapArr)) // 5 business days in a week
		assert.Equal(t, "2025-01-27T14:30:00.000Z", jsonMapArr[0].Date)
		assert.Equal(t, float32(33.89), jsonMapArr[0].Low)
		assert.Equal(t, float32(34.8), jsonMapArr[0].High)
	})

	t.Run("GetHistoricalWithDecode Invalid request with incorrect dates", func(t *testing.T) {
		// arrange
		sdate, edate := "1900-01-01", "1900-01-01"

		// act
		jsonMapArr, err := yahooApi.GetHistoricalWitDecode(ticker, sdate, edate)

		// assert
		assert.NotNil(t, err)
		assert.Empty(t, jsonMapArr)
	})
}

func Test_GetSingleStockPrice(t *testing.T) {
	yahooApi := NewYHFinanceCompleteAPI(log.New(os.Stdout, "", log.Ltime))
	ticker := "BCE.TO"

	t.Run("GetSingleStockPrice Valid request", func(t *testing.T) {
		// act
		priceResponse, err := yahooApi.GetSingleStockPrice(ticker)

		// assert
		assert.Nil(t, err)
		assert.NotEmpty(t, priceResponse)
		assert.Equal(t, ticker, priceResponse.Symbol)
		assert.NotZero(t, priceResponse.Price)
		assert.NotEmpty(t, priceResponse.Currency)
		assert.NotZero(t, priceResponse.MarketCap)
	})

	t.Run("GetSingleStockPrice Invalid request with empty ticker", func(t *testing.T) {
		// act
		priceResponse, err := yahooApi.GetSingleStockPrice("")

		// assert
		assert.NotNil(t, err)
		assert.Empty(t, priceResponse)
	})

	t.Run("GetSingleStockPrice Invalid request with non-existent ticker", func(t *testing.T) {
		// act
		priceResponse, err := yahooApi.GetSingleStockPrice("INVALID_TICKER")

		// assert
		assert.NotNil(t, err)
		assert.Empty(t, priceResponse)
	})
}

func Test_GetStockSummaryDetail(t *testing.T) {
	yahooApi := NewYHFinanceCompleteAPI(log.New(os.Stdout, "", log.Ltime))
	ticker := "BCE.TO"

	t.Run("GetStockSummaryDetail Valid request", func(t *testing.T) {
		// act
		summaryResponse, err := yahooApi.GetStockSummaryDetail(ticker)

		// assert
		assert.Nil(t, err)
		assert.NotEmpty(t, summaryResponse)
		assert.Equal(t, ticker, summaryResponse.Price.Symbol)
		assert.NotZero(t, summaryResponse.Price.RegularMarketPrice)
		assert.NotEmpty(t, summaryResponse.Price.Currency)
		assert.NotZero(t, summaryResponse.Price.MarketCap)
	})

	t.Run("GetStockSummaryDetail Invalid request with empty ticker", func(t *testing.T) {
		// act
		summaryResponse, err := yahooApi.GetStockSummaryDetail("")

		// assert
		assert.NotNil(t, err)
		assert.Empty(t, summaryResponse)
	})

	t.Run("GetStockSummaryDetail Invalid request with non-existent ticker", func(t *testing.T) {
		// act
		summaryResponse, err := yahooApi.GetStockSummaryDetail("INVALID_TICKER")

		// assert
		assert.NotNil(t, err)
		assert.Empty(t, summaryResponse)
	})
}

// func Test_YHFinanceCompleteAPI_buildRequest_AI(t *testing.T) {
// 	api := NewYHFinanceCompleteAPI(log.New(os.Stdout, "", log.Ltime))

// 	t.Run("Valid request", func(t *testing.T) {
// 		// arrange
// 		queryParams := url.Values{"ticker": {"BCE.TO"}, "sdate": {"2025-02-10"}, "edate": {"2025-02-11"}}

// 		// act
// 		req, err := api.buildRequest("yhfhistorical", queryParams)

// 		// assert
// 		assert.NotNil(t, req)
// 		assert.Nil(t, err)
// 		assert.Equal(t, "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?edate=2025-02-11&sdate=2025-02-10&ticker=BCE.TO", req.URL.String())
// 		assert.Equal(t, "yh-finance-complete.p.rapidapi.com", req.Header.Get("x-rapidapi-host"))
// 		assert.Equal(t, "9b405718ddmsh954d4191ebcf658p148c17jsn58521162b938", req.Header.Get("x-rapidapi-key"))
// 	})

// 	t.Run("Invalid subDir", func(t *testing.T) {
// 		// arrange
// 		queryParams := url.Values{"ticker": {"BCE.TO"}, "sdate": {"2025-02-10"}, "edate": {"2025-02-11"}}

// 		// act
// 		req, err := api.buildRequest("", queryParams)

// 		// assert
// 		assert.NotNil(t, req)
// 		assert.Nil(t, err)
// 		assert.Equal(t, "https://yh-finance-complete.p.rapidapi.com/?edate=2025-02-11&sdate=2025-02-10&ticker=BCE.TO", req.URL.String())
// 	})

// 	t.Run("Empty queryParams", func(t *testing.T) {
// 		// arrange
// 		queryParams := url.Values{}

// 		// act
// 		req, err := api.buildRequest("yhfhistorical", queryParams)

// 		// assert
// 		assert.NotNil(t, req)
// 		assert.Nil(t, err)
// 		assert.Equal(t, "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?", req.URL.String())
// 	})

// 	t.Run("Nil queryParams", func(t *testing.T) {
// 		// arrange
// 		var queryParams url.Values

// 		// act
// 		req, err := api.buildRequest("yhfhistorical", queryParams)

// 		// assert
// 		assert.NotNil(t, req)
// 		assert.Nil(t, err)
// 		assert.Equal(t, "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?", req.URL.String())
// 	})
// }

// func Test_YHFinanceCompleteAPI_buildRequest(t *testing.T) {
// 	// arrange
// 	api := NewYHFinanceCompleteAPI(log.New(os.Stdout, "", log.Ltime))
// 	queryParams := url.Values{"ticker": {"BCE.TO"}, "sdate": {"2025-02-10"}, "edate": {"2025-02-11"}}

// 	// act
// 	req, err := api.buildRequest("yhfhistorical", queryParams)

// 	// assert
// 	assert.NotNil(t, req)
// 	assert.Nil(t, err)

// 	// fmt.Printf("req=%+v\n", req)
// 	// fmt.Printf("url=%+v\n", req.URL)

// 	toMatch := "https://yh-finance-complete.p.rapidapi.com/yhfhistorical?edate=2025-02-11&sdate=2025-02-10&ticker=BCE.TO"
// 	assert.Equal(t, req.URL.String(), toMatch)

// 	// Fairly simple test, no need for a negative one
// }

package fintechapi

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// go test .\fintechapi\. -v
// go test YHFinanceComplete_test.go .\YHFinanceComplete.go  -v

func getStertAndEndDSates() (string, string) {
	// const daysBack = 2
	// sdate := time.Now().AddDate(0,0,(-1 * (daysBack+3))).Format("2006-01-02")
	// edate := time.Now().AddDate(0,0,-2).Format("2006-01-02")
	// fmt.Println("Test_yhfhistorical()", daysBack, sdate, edate)

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	weekStart := weekStartDate(firstOfMonth)
	weekEnd := weekStart.AddDate(0, 0, 5)

	sdate := weekStart.Format("2006-01-02")
	edate := weekEnd.Format("2006-01-02")

	fmt.Printf("firstOfMonth=%s weekStart=%s weekEnd=%s\n", firstOfMonth.Format("2006-01-02 Mon"), sdate, edate)

	return sdate, edate
}

func Test_buildRequest(t *testing.T) {
	ticker, sdate, edate := "TSLA", "2025-02-10", "2025-02-10"
	req, err := buildRequest(ticker, sdate, edate)
	assert.NotNil(t, req)
	assert.Nil(t, err)

	// req, err = buildRequest("", "", "")
	// assert.NotNil(t, req)
	// assert.Nil(t, err)
}

func Test_yhfhistoricalUnmarshal(t *testing.T) {
	// arrange
	sdate, edate := getStertAndEndDSates()

	// act
	jsonMapArrUnmarshal, errUnmarshal := yhfhistoricalUnmarshal("TSLA", sdate, edate)

	// assert
	assert.NotEmpty(t, jsonMapArrUnmarshal)

	assert.Nil(t, errUnmarshal)

	// act
	jsonMapArrUnmarshal, errUnmarshal = yhfhistoricalUnmarshal("TSLA", "01-01-1900", "01-01-1900")	

	// assert
	assert.Empty(t, jsonMapArrUnmarshal)

	assert.NotNil(t, errUnmarshal)

	fmt.Println("errUnmarshal 	--->", errUnmarshal)
}


func yhfhistoricalFull(t *testing.T) {
	// arrange
	sdate, edate := getStertAndEndDSates()

	// act
	jsonMapArrUnmarshal, errUnmarshal := yhfhistoricalUnmarshal("TSLA", sdate, edate)
	jsonMapArrDecode, errDecode := yhfhistoricalDecode("TSLA", sdate, edate)

	// assert
	assert.NotEmpty(t, jsonMapArrUnmarshal)
	assert.NotEmpty(t, jsonMapArrDecode)

	assert.Nil(t, errUnmarshal)
	assert.Nil(t, errDecode)

	// act
	jsonMapArrUnmarshal, errUnmarshal = yhfhistoricalUnmarshal("TSLA", "01-01-1900", "01-01-1900")	
	jsonMapArrDecode, errDecode = yhfhistoricalDecode("TSLA", "01-01-1900", "01-01-1900")

	// assert
	assert.Empty(t, jsonMapArrUnmarshal)
	assert.Empty(t, jsonMapArrDecode)

	assert.NotNil(t, errUnmarshal)
	assert.NotNil(t, errDecode)

	fmt.Println("errUnmarshal 	--->", errUnmarshal)
	fmt.Println("errDecode		--->", errDecode)

}

func TestWeekStartDate(t *testing.T) {
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 7; i++ {
		weekStart := weekStartDate(date)
		// fmt.Printf("%s %s\n", date.Format("2006-01-02 Mon"), weekStart.Format("2006-01-02 Mon"))
		assert.NotNil(t, weekStart)
		assert.Equal(t, time.Monday, weekStart.Weekday())
		date = date.Add(24 * time.Hour)
	}
	date = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 7; i++ {
		weekStart := weekStartDate(date)
		// fmt.Printf("%s %s\n", date.Format("2006-01-02 Mon"), weekStart.Format("2006-01-02 Mon"))
		assert.NotNil(t, weekStart)
		assert.Equal(t, time.Monday, weekStart.Weekday())
		date = date.Add(24 * time.Hour)
	}
}

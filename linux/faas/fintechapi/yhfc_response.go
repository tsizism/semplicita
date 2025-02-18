package fintechapi

// https://rapidapi.com/belchiorarkad-FqvHs2EDOtP/api/yh-finance-complete
// https://algotrading101.com/learn/yahoo-finance-api-guide/
// Genrate go struct with json tags from json response

type YfhistoricalResponse struct {
	Date     string  `json:"date"`
	Adjclose float32 `json:"adjclose"`
	Close    float32 `json:"close"`
	High     float32 `json:"high"`
	Low      float32 `json:"low"`
	Open     float32 `json:"open"`
	Volume   uint64  `json:"volume"`
	Symbol   string  `json:"symbol"`
}

// symbol:"BCE.TO"
// price:33.78
// currency:"CAD"
// marketCap:30816919552

type YfpriceResponse struct {
	Symbol    string  `json:"symbol"`
	Price     float32 `json:"price"`
	Currency  string  `json:"currency"`
	MarketCap uint64  `json:"MarketCap"`
}

type YffullstockpriceResponse struct {
	Price Price `json:"price"`
}

type Price struct {
	MaxAge                    int     `json:"maxAge"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
	RegularMarketChange        float64 `json:"regularMarketChange"`
	RegularMarketTime          string  `json:"regularMarketTime"`
	PriceHint                  int     `json:"priceHint"`
	RegularMarketPrice         float64 `json:"regularMarketPrice"`
	RegularMarketDayHigh       float64 `json:"regularMarketDayHigh"`
	RegularMarketDayLow        float64 `json:"regularMarketDayLow"`
	RegularMarketVolume        int     `json:"regularMarketVolume"`
	RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
	RegularMarketSource        string  `json:"regularMarketSource"`
	RegularMarketOpen          float64 `json:"regularMarketOpen"`
	Exchange                   string  `json:"exchange"`
	ExchangeName               string  `json:"exchangeName"`
	ExchangeDataDelayedBy      int     `json:"exchangeDataDelayedBy"`
	MarketState                string  `json:"marketState"`
	QuoteType                  string  `json:"quoteType"`
	Symbol                     string  `json:"symbol"`
	UnderlyingSymbol           *string `json:"underlyingSymbol"`
	ShortName                  string  `json:"shortName"`
	LongName                   string  `json:"longName"`
	Currency                   string  `json:"currency"`
	QuoteSourceName            string  `json:"quoteSourceName"`
	CurrencySymbol             string  `json:"currencySymbol"`
	FromCurrency               *string `json:"fromCurrency"`
	ToCurrency                 *string `json:"toCurrency"`
	LastMarket                 *string `json:"lastMarket"`
	MarketCap                  int64   `json:"marketCap"`
}

/*
{"price": {
    "maxAge": 1,
    "regularMarketChangePercent": 0.00079536065,
    "regularMarketChange": 0.069999695,
    "regularMarketTime": "2025-02-18T17:45:30.000Z",
    "priceHint": 2,
    "regularMarketPrice": 88.08,
    "regularMarketDayHigh": 88.43,
    "regularMarketDayLow": 87.44,
    "regularMarketVolume": 442930,
    "regularMarketPreviousClose": 88.01,
    "regularMarketSource": "FREE_REALTIME",
    "regularMarketOpen": 87.67,
    "exchange": "TOR",
    "exchangeName": "Toronto",
    "exchangeDataDelayedBy": 15,
    "marketState": "REGULAR",
    "quoteType": "EQUITY",
    "symbol": "CM.TO",
    "underlyingSymbol": null,
    "shortName": "CANADIAN IMPERIAL BANK OF COMME",
    "longName": "Canadian Imperial Bank of Commerce",
    "currency": "CAD",
    "quoteSourceName": "Free Realtime Quote",
    "currencySymbol": "$",
    "fromCurrency": null,
    "toCurrency": null,
    "lastMarket": null,
    "marketCap": 83004391424
}}
*/

type YfSummaryDetail struct {
	MaxAge                        int     `json:"maxAge"`
	PriceHint                     int     `json:"priceHint"`
	PreviousClose                 float32 `json:"previousClose"`
	Open                          float32 `json:"open"`
	DayLow                        float32 `json:"dayLow"`
	DayHigh                       float32 `json:"dayHigh"`
	RegularMarketPreviousClose    float32 `json:"regularMarketPreviousClose"`
	RegularMarketOpen             float32 `json:"regularMarketOpen"`
	RegularMarketDayLow           float32 `json:"regularMarketDayLow"`
	RegularMarketDayHigh          float32 `json:"regularMarketDayHigh"`
	DividendRate                  float32 `json:"dividendRate"`
	DividendYield                 float32 `json:"dividendYield"`
	ExDividendDate                string  `json:"exDividendDate"`
	PayoutRatio                   float32 `json:"payoutRatio"`
	FiveYearAvgDividendYield      float32 `json:"fiveYearAvgDividendYield"`
	Beta                          float32 `json:"beta"`
	TrailingPE                    float32 `json:"trailingPE"`
	ForwardPE                     float32 `json:"forwardPE"`
	Volume                        uint64  `json:"volume"`
	RegularMarketVolume           uint64  `json:"regularMarketVolume"`
	AverageVolume                 uint64  `json:"averageVolume"`
	AverageVolume10days           uint64  `json:"averageVolume10days"`
	AverageDailyVolume10Day       uint64  `json:"averageDailyVolume10Day"`
	Bid                           float32 `json:"bid"`
	Ask                           float32 `json:"ask"`
	BidSize                       int     `json:"bidSize"`
	AskSize                       int     `json:"askSize"`
	MarketCap                     uint64  `json:"marketCap"`
	FiftyTwoWeekLow               float32 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh              float32 `json:"fiftyTwoWeekHigh"`
	PriceToSalesTrailing12Months  float32 `json:"priceToSalesTrailing12Months"`
	FiftyDayAverage               float32 `json:"fiftyDayAverage"`
	TwoHundredDayAverage          float32 `json:"twoHundredDayAverage"`
	TrailingAnnualDividendRate    float32 `json:"trailingAnnualDividendRate"`
	TrailingAnnualDividendYield   float32 `json:"trailingAnnualDividendYield"`
	Currency                      string  `json:"currency"`
	FromCurrency                  *string `json:"fromCurrency"`
	ToCurrency                    *string `json:"toCurrency"`
	LastMarket                    *string `json:"lastMarket"`
	CoinMarketCapLink             *string `json:"coinMarketCapLink"`
	Algorithm                     *string `json:"algorithm"`
	Tradeable                     bool    `json:"tradeable"`
}

type YfPrice struct {
	MaxAge                      int     `json:"maxAge"`
	RegularMarketChangePercent  float32 `json:"regularMarketChangePercent"`
	RegularMarketChange         float32 `json:"regularMarketChange"`
	RegularMarketTime           string  `json:"regularMarketTime"`
	PriceHint                   int     `json:"priceHint"`
	RegularMarketPrice          float32 `json:"regularMarketPrice"`
	RegularMarketDayHigh        float32 `json:"regularMarketDayHigh"`
	RegularMarketDayLow         float32 `json:"regularMarketDayLow"`
	RegularMarketVolume         uint64  `json:"regularMarketVolume"`
	AverageDailyVolume10Day     uint64  `json:"averageDailyVolume10Day"`
	AverageDailyVolume3Month    uint64  `json:"averageDailyVolume3Month"`
	RegularMarketPreviousClose  float32 `json:"regularMarketPreviousClose"`
	RegularMarketSource         string  `json:"regularMarketSource"`
	RegularMarketOpen           float32 `json:"regularMarketOpen"`
	Exchange                    string  `json:"exchange"`
	ExchangeName                string  `json:"exchangeName"`
	ExchangeDataDelayedBy       int     `json:"exchangeDataDelayedBy"`
	MarketState                 string  `json:"marketState"`
	QuoteType                   string  `json:"quoteType"`
	Symbol                      string  `json:"symbol"`
	UnderlyingSymbol            *string `json:"underlyingSymbol"`
	ShortName                   string  `json:"shortName"`
	LongName                    string  `json:"longName"`
	Currency                    string  `json:"currency"`
	QuoteSourceName             string  `json:"quoteSourceName"`
	CurrencySymbol              string  `json:"currencySymbol"`
	FromCurrency                *string `json:"fromCurrency"`
	ToCurrency                  *string `json:"toCurrency"`
	LastMarket                  *string `json:"lastMarket"`
	MarketCap                   uint64  `json:"marketCap"`
}

type YfResponse struct {
	SummaryDetail YfSummaryDetail `json:"summaryDetail"`
	Price         YfPrice         `json:"price"`
}



/*
{
    "summaryDetail": {
        "maxAge": 1,
        "priceHint": 2,
        "previousClose": 33.52,
        "open": 33.51,
        "dayLow": 33.33,
        "dayHigh": 33.85,
        "regularMarketPreviousClose": 33.52,
        "regularMarketOpen": 33.51,
        "regularMarketDayLow": 33.33,
        "regularMarketDayHigh": 33.85,
        "dividendRate": 3.99,
        "dividendYield": 0.1181,
        "exDividendDate": "2025-03-14T00:00:00.000Z",
        "payoutRatio": 44,
        "fiveYearAvgDividendYield": 6.57,
        "beta": 0.435,
        "trailingPE": 375.3333,
        "forwardPE": 11.608247,
        "volume": 4411591,
        "regularMarketVolume": 4411591,
        "averageVolume": 4471525,
        "averageVolume10days": 5283050,
        "averageDailyVolume10Day": 5283050,
        "bid": 33.7,
        "ask": 33.82,
        "bidSize": 0,
        "askSize": 0,
        "marketCap": 30816919552,
        "fiftyTwoWeekLow": 31.43,
        "fiftyTwoWeekHigh": 51.57,
        "priceToSalesTrailing12Months": 1.2598904,
        "fiftyDayAverage": 34.1994,
        "twoHundredDayAverage": 42.35395,
        "trailingAnnualDividendRate": 3.96,
        "trailingAnnualDividendYield": 0.118138425,
        "currency": "CAD",
        "fromCurrency": null,
        "toCurrency": null,
        "lastMarket": null,
        "coinMarketCapLink": null,
        "algorithm": null,
        "tradeable": false
    },
    "price": {
        "maxAge": 1,
        "regularMarketChangePercent": 0.0077565135,
        "regularMarketChange": 0.25999832,
        "regularMarketTime": "2025-02-14T21:00:00.000Z",
        "priceHint": 2,
        "regularMarketPrice": 33.78,
        "regularMarketDayHigh": 33.85,
        "regularMarketDayLow": 33.33,
        "regularMarketVolume": 4411591,
        "averageDailyVolume10Day": 5283050,
        "averageDailyVolume3Month": 4471525,
        "regularMarketPreviousClose": 33.52,
        "regularMarketSource": "FREE_REALTIME",
        "regularMarketOpen": 33.51,
        "exchange": "TOR",
        "exchangeName": "Toronto",
        "exchangeDataDelayedBy": 15,
        "marketState": "CLOSED",
        "quoteType": "EQUITY",
        "symbol": "BCE.TO",
        "underlyingSymbol": null,
        "shortName": "BCE INC.",
        "longName": "BCE Inc.",
        "currency": "CAD",
        "quoteSourceName": "Delayed Quote",
        "currencySymbol": "$",
        "fromCurrency": null,
        "toCurrency": null,
        "lastMarket": null,
        "marketCap": 30816919552
    }
}
	
*/


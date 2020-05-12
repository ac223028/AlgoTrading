package alphaVantage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"
)

//func PrettyPrint(v interface{}) (err error) {
//	b, err := json.MarshalIndent(v, "", "  ")
//	if err == nil {
//		fmt.Println(string(b))
//	}
//	return
//}

const baseURL = "https://www.alphavantage.co"
const httpDelayPerRequest = time.Second * 2 // may have to do something about this

// Client represents a new alphavantage client
type Client struct {
	apiKey          string
	httpClient      *http.Client
	httpNextRequest time.Time
	sync.Mutex
}

// New creates new Client instance
func New(apiKey string) *Client {
	const httpTimeout = time.Second * 30

	httpClient := &http.Client{
		Timeout: httpTimeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 5,
		},
	}

	return &Client{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c *Client) MakeHTTPRequest(url string) ([]byte, error) {
	c.Lock()
	defer c.Unlock()

	// Run request only every x seconds (determined by httpNextRequest)
	now := time.Now()
	if now.Before(c.httpNextRequest) {
		ticker := time.NewTicker(c.httpNextRequest.Sub(now))
		<-ticker.C
	}
	defer func(c *Client) {
		c.httpNextRequest = time.Now().Add(httpDelayPerRequest)
	}(c)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("building http request failed: %w", err)
	}
	req.Header.Set("User-Agent", "Go client: Anon")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close() // not sure what this does

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: expected %d, got %d",
			http.StatusOK, resp.StatusCode)
	}

	return body, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// IndicatorRSI represents the overall struct for RSI indicator
// Example https://www.alphavantage.co/query?function=RSI&symbol=MSFT&interval=daily&apikey=demo
type IndicatorRSI struct {
	Metadata          IndicatorRSIMetadata            `json:"Meta Data"`
	TechnicalAnalysis map[string]TechnicalRSIAnalysis `json:"Technical Analysis: RSI"`
}

// IndicatorRSIMetadata is the metadata subset of IndicatorRSI
type IndicatorRSIMetadata struct {
	Symbol        string `json:"1: Symbol"`
	Indicator     string `json:"2: Indicator"`
	LastRefreshed string `json:"3: Last Refreshed"`
	Interval      string `json:"4: Interval"`
	TimePeriod    int    `json:"5: Time Period"`
	SeriesType    string `json:"6: Series Type"`
	TimeZone      string `json:"7: Time Zone"`
}

// TechnicalRSIAnalysis is the RSI indicator subset of IndicatorRSI
type TechnicalRSIAnalysis struct {
	RSI float64 `json:",string"`
}

func toIndicatorRSI(buf []byte) (*IndicatorRSI, error) {
	indicatorRSI := &IndicatorRSI{}
	if err := json.Unmarshal(buf, indicatorRSI); err != nil {
		return nil, err
	}
	return indicatorRSI, nil
}

// IndicatorRSI fetches the "RSI" indicators for given symbol from API.
// The order of dates in TechnicalAnalysis is random because it's a map.
func (c *Client) IndicatorRSI(symbol string, interval string, timePeriod string, seriesType string) (*IndicatorRSI, error) { // come back to make them enums?
	// the daily RSI 14 close is Alpaca's Wilder's 1 year RSI 14

	url := fmt.Sprintf("%s/query?function=%s&symbol=%s&interval=%s&time_period=%s&series_type=%s&apikey=%s",
		baseURL, "RSI", symbol, interval, timePeriod, seriesType, c.apiKey)
	body, err := c.MakeHTTPRequest(url)
	if err != nil {
		return nil, err
	}
	indicator, err := toIndicatorRSI(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return indicator, nil
}

// Latest returns the most recent TechnicalRSIAnalysis for given RSI.
func (RSI *IndicatorRSI) Latest() (date string, latest *TechnicalRSIAnalysis) { // this should work regardless of time
	dates := make([]string, len(RSI.TechnicalAnalysis))
	for date := range RSI.TechnicalAnalysis {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	date = dates[len(dates)-1]
	latestVal, _ := RSI.TechnicalAnalysis[date]
	latest = &latestVal
	return
}

// Today returns TechnicalRSIAnalysis for today.
func (RSI *IndicatorRSI) Today() *TechnicalRSIAnalysis {
	today := time.Now()
	return RSI.ByDate(today)
}

// ByDate returns TechnicalRSIAnalysis for the given date.
func (RSI *IndicatorRSI) ByDate(date time.Time) *TechnicalRSIAnalysis {
	day := date.Format("2006-01-02") // not sure what this string does
	//print(day, "\n")
	item, exists := RSI.TechnicalAnalysis[day]
	if !exists {
		return nil
	}
	return &item
}
func (RSI *IndicatorRSI) ByMinute(date time.Time) *TechnicalRSIAnalysis {
	day := date.Format("2006-01-02 15:04") // this string needs to be changed
	//print(day, "\n")
	item, exists := RSI.TechnicalAnalysis[day]
	if !exists {
		return nil
	}
	return &item
}

func (RSI *IndicatorRSI) GetRSI() (float32, []float32) { // this is daily avg?

	var latest float32
	//var array []float32
	var array []float32

	if len(RSI.TechnicalAnalysis) < 90 {
		array = make([]float32, len(RSI.TechnicalAnalysis))
	} else {
		array = make([]float32, 90)
	}
	var dates []string
	for date := range RSI.TechnicalAnalysis {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	date := dates[len(dates)-1]
	latestVal, _ := RSI.TechnicalAnalysis[date]
	latest = float32(latestVal.RSI)

	//for i := range dates {
	//	print(i, " ", dates[i], "\n")
	//}

	for i := range array {
		r := RSI.TechnicalAnalysis[dates[i]].RSI
		//print(i, " ", dates[i], " ", r, "\n")
		//array = append(array, float32(r))
		array[i] = float32(r)
	}

	//for i := range array {
	//	print(i, " ", array[i], "\n")
	//}

	//print(latest, " ")
	//for i := range array {
	//	print(array[i], " ")
	//}

	return latest, array

}

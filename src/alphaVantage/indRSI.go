package alphaVantage

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

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
func (c *Client) IndicatorRSI(symbol string, interval string, timePeriod string, seriesType string) (*IndicatorRSI, error) {

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

	if len(indicator.TechnicalAnalysis) == 0 {
		return nil, fmt.Errorf("there is no indicator data: %w", err)
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

func (RSI *IndicatorRSI) GetRSI() (float32, []float32) { // this is a mess, test as is

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

	for i := range array {
		r := RSI.TechnicalAnalysis[dates[i]].RSI
		array[i] = float32(r)
	}

	return latest, array

}

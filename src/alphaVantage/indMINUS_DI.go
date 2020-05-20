package alphaVantage

import (
	"encoding/json"
	"fmt"
	"sort"
)

type IndicatorMINUS_DI struct {
	Metadata          IndicatorMINUS_DIMetadata            `json:"Meta Data"`
	TechnicalAnalysis map[string]TechnicalMINUS_DIAnalysis `json:"Technical Analysis: MINUS_DI"`
}

// IndicatorMINUS_DIMetadata is the metadata subset of IndicatorMINUS_DI
type IndicatorMINUS_DIMetadata struct {
	Symbol        string `json:"1: Symbol"`
	Indicator     string `json:"2: Indicator"`
	LastRefreshed string `json:"3: Last Refreshed"`
	Interval      string `json:"4: Interval"`
	TimePeriod    int    `json:"5: Time Period"`
	TimeZone      string `json:"6: Time Zone"`
}

// TechnicalMINUS_DIAnalysis is the MINUS_DI indicator subset of IndicatorMINUS_DI
type TechnicalMINUS_DIAnalysis struct {
	MINUS_DI float64 `json:",string"`
}

func toIndicatorMINUS_DI(buf []byte) (*IndicatorMINUS_DI, error) {
	indicatorMINUS_DI := &IndicatorMINUS_DI{}
	if err := json.Unmarshal(buf, indicatorMINUS_DI); err != nil {
		return nil, err
	}
	return indicatorMINUS_DI, nil
}

// IndicatorMINUS_DI fetches the "MINUS_DI" indicators for given symbol from API.
// The order of dates in TechnicalAnalysis is random because it's a map.
func (c *Client) IndicatorMINUS_DI(symbol string, interval string, timePeriod string) (*IndicatorMINUS_DI, error) {

	url := fmt.Sprintf("%s/query?function=%s&symbol=%s&interval=%s&time_period=%s&apikey=%s",
		baseURL, "MINUS_DI", symbol, interval, timePeriod, c.apiKey)
	body, err := c.MakeHTTPRequest(url)
	if err != nil {
		return nil, err
	}
	indicator, err := toIndicatorMINUS_DI(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(indicator.TechnicalAnalysis) == 0 {
		return nil, fmt.Errorf("there is no indicator data: %w", err)
	}

	return indicator, nil
}

func (MINUS_DI *IndicatorMINUS_DI) Latest() (date string, latest *TechnicalMINUS_DIAnalysis) { // this should work regardless of time
	dates := make([]string, len(MINUS_DI.TechnicalAnalysis))
	for date := range MINUS_DI.TechnicalAnalysis {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	date = dates[len(dates)-1]
	latestVal, _ := MINUS_DI.TechnicalAnalysis[date]
	latest = &latestVal
	return
}

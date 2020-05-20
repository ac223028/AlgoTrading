package alphaVantage

import (
	"encoding/json"
	"fmt"
	"sort"
)

type IndicatorPLUS_DI struct {
	Metadata          IndicatorPLUS_DIMetadata            `json:"Meta Data"`
	TechnicalAnalysis map[string]TechnicalPLUS_DIAnalysis `json:"Technical Analysis: PLUS_DI"`
}

// IndicatorPLUS_DIMetadata is the metadata subset of IndicatorPLUS_DI
type IndicatorPLUS_DIMetadata struct {
	Symbol        string `json:"1: Symbol"`
	Indicator     string `json:"2: Indicator"`
	LastRefreshed string `json:"3: Last Refreshed"`
	Interval      string `json:"4: Interval"`
	TimePeriod    int    `json:"5: Time Period"`
	TimeZone      string `json:"6: Time Zone"`
}

// TechnicalPLUS_DIAnalysis is the PLUS_DI indicator subset of IndicatorPLUS_DI
type TechnicalPLUS_DIAnalysis struct {
	PLUS_DI float64 `json:",string"`
}

func toIndicatorPLUS_DI(buf []byte) (*IndicatorPLUS_DI, error) {
	indicatorPLUS_DI := &IndicatorPLUS_DI{}
	if err := json.Unmarshal(buf, indicatorPLUS_DI); err != nil {
		return nil, err
	}
	return indicatorPLUS_DI, nil
}

// IndicatorPLUS_DI fetches the "PLUS_DI" indicators for given symbol from API.
// The order of dates in TechnicalAnalysis is random because it's a map.
func (c *Client) IndicatorPLUS_DI(symbol string, interval string, timePeriod string) (*IndicatorPLUS_DI, error) {

	url := fmt.Sprintf("%s/query?function=%s&symbol=%s&interval=%s&time_period=%s&apikey=%s",
		baseURL, "PLUS_DI", symbol, interval, timePeriod, c.apiKey)
	body, err := c.MakeHTTPRequest(url)
	if err != nil {
		return nil, err
	}
	indicator, err := toIndicatorPLUS_DI(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(indicator.TechnicalAnalysis) == 0 {
		return nil, fmt.Errorf("there is no indicator data: %w", err)
	}

	return indicator, nil
}

func (PLUS_DI *IndicatorPLUS_DI) Latest() (date string, latest *TechnicalPLUS_DIAnalysis) { // this should work regardless of time
	dates := make([]string, len(PLUS_DI.TechnicalAnalysis))
	for date := range PLUS_DI.TechnicalAnalysis {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	date = dates[len(dates)-1]
	latestVal, _ := PLUS_DI.TechnicalAnalysis[date]
	latest = &latestVal
	return
}

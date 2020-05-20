package alphaVantage

import (
	"encoding/json"
	"fmt"
	"sort"
)

type IndicatorADX struct {
	Metadata          IndicatorADXMetadata            `json:"Meta Data"`
	TechnicalAnalysis map[string]TechnicalADXAnalysis `json:"Technical Analysis: ADX"`
}

// IndicatorADXMetadata is the metadata subset of IndicatorADX
type IndicatorADXMetadata struct {
	Symbol        string `json:"1: Symbol"`
	Indicator     string `json:"2: Indicator"`
	LastRefreshed string `json:"3: Last Refreshed"`
	Interval      string `json:"4: Interval"`
	TimePeriod    int    `json:"5: Time Period"`
	TimeZone      string `json:"6: Time Zone"`
}

// TechnicalADXAnalysis is the ADX indicator subset of IndicatorADX
type TechnicalADXAnalysis struct {
	ADX float64 `json:",string"`
}

func toIndicatorADX(buf []byte) (*IndicatorADX, error) {
	indicatorADX := &IndicatorADX{}
	if err := json.Unmarshal(buf, indicatorADX); err != nil {
		return nil, err
	}
	return indicatorADX, nil
}

// IndicatorADX fetches the "ADX" indicators for given symbol from API.
// The order of dates in TechnicalAnalysis is random because it's a map.
func (c *Client) IndicatorADX(symbol string, interval string, timePeriod string) (*IndicatorADX, error) {

	url := fmt.Sprintf("%s/query?function=%s&symbol=%s&interval=%s&time_period=%s&apikey=%s",
		baseURL, "ADX", symbol, interval, timePeriod, c.apiKey)
	body, err := c.MakeHTTPRequest(url)
	if err != nil {
		return nil, err
	}
	indicator, err := toIndicatorADX(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(indicator.TechnicalAnalysis) == 0 {
		return nil, fmt.Errorf("there is no indicator data: %w", err)
	}

	return indicator, nil
}

func (ADX *IndicatorADX) Latest() (date string, latest *TechnicalADXAnalysis) { // this should work regardless of time
	dates := make([]string, len(ADX.TechnicalAnalysis))
	for date := range ADX.TechnicalAnalysis {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	date = dates[len(dates)-1]
	latestVal, _ := ADX.TechnicalAnalysis[date]
	latest = &latestVal
	return
}

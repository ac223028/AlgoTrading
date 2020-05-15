package alphaVantage

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

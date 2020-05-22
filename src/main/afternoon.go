package main

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func getDataPrices(alpacaAPI *alpaca.Client, timeT time.Time) {
	fromFile := timeT.Format("2006-01-02")                         // TODO: fix this naming scheme
	dat, err := ioutil.ReadFile("C:/Trading/" + fromFile + ".txt") // reads assets form file
	if err != nil {
		dat, err = ioutil.ReadFile(fromFile + ".txt") // reads assets form file
		if err != nil {
			panic(err)
		}
	}
	//print(string(dat))
	tips := strings.Split(string(dat), "\n")

	param := alpaca.ListBarParams{
		Timeframe: "day",
		StartDt:   &timeT,
		EndDt:     nil,
		Limit:     nil,
	}

	file, err := os.Create("C:/Trading/" + fromFile + "-sell-data.txt")
	if err != nil {
		file, _ = os.Create(fromFile + "data.txt")
	}

	for j, i := range tips {
		print(j, " ", i)
	}

	for _, symbol := range tips {
		bars, err := alpacaAPI.GetSymbolBars(symbol, param)
		if err != nil {
			panic(err)
		}
		o := float32(0)
		c := float32(0)

		for _, i := range bars {
			if time.Unix(i.Time, 0).Format("2006-01-02") == timeT.Format("2006-01-02") {
				c = i.Close
				o = i.Open
				break
			}
		}

		op := fmt.Sprintf("%f", o)
		cl := fmt.Sprintf("%f", c)

		output := symbol + " " + op + " " + cl + "\n"
		file.WriteString(output)
	}

	file.Close()
}

func todaysBuyPrices(alpacaAPI *alpaca.Client, timeT time.Time) string {

}

func yesterdaysSellPrices(alpacaAPI *alpaca.Client, timeT time.Time) string {

}

func main() {
	// TODO: This is where manager would be implemented
	API_KEY := "PK8SJ2WJQT189BUFK9MJ"
	API_SECRET := "mCn1jZIs3Eu2GFVgT44O/aIiFLGswYQKiBh0w1ll"
	BASE_URL := "https://paper-api.alpaca.markets"
	alpaca.SetBaseUrl(BASE_URL)

	if common.Credentials().ID == "" {
		os.Setenv(common.EnvApiKeyID, API_KEY)
	}
	if common.Credentials().Secret == "" {
		os.Setenv(common.EnvApiSecretKey, API_SECRET)
	}

	os.Setenv(common.EnvApiSecretKey, API_SECRET)
	//AlpClient := alpaca.NewClient(common.Credentials())

	//	t := time.Now()
	//	today := todaysBuyPrices(AlpClient, t)
	//	yesterday := yesterdaysSellPrices(AlpClient, t)
}

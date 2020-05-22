package main

import (
	"../alphaVantage"
	"../trendFollowing"
	"encoding/json"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// TODO: work on error messages...

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func checkError(e error) {
	if e != nil {
		PrettyPrint(e.Error())
	}
}

func getIndData(alpacaAPI *alpaca.Client, avAPI *alphaVantage.Client, date time.Time, fileName string,
	interval string, timePeriod string, seriesType string) { // need to get EMA(RSI) // think about returning string
	// still broken for some reason
	// date in "YYYY-MM-DD" format or "2006-01-02"
	//interval := "weekly"
	//timePeriod := "14"
	//seriesType := "close"

	dat, err := ioutil.ReadFile(fileName + "-tipsheet.txt") // reads assets form file // CHANGE TO BE DATE UNIQUE
	if err != nil {
		panic(err)
	}
	affordableAssets := strings.Split(strings.TrimSpace(string(dat)), "\n") // this is a set of tickers and prices

	outputName := date.Format("2006-01-02") + "-indicator-data.txt"

	file, e := os.Create("C:/Trading/" + outputName)
	if e != nil { // if not there make relative
		file, _ = os.Create(outputName)
	}

	for _, tuple := range affordableAssets {

		a := strings.Split(tuple, " ")
		symbol := a[0]
		//price := a[1] // this is old data

		INDrsi, err := avAPI.IndicatorRSI(symbol, interval, timePeriod, seriesType)
		checkError(err)
		INDadx, err := avAPI.IndicatorADX(symbol, interval, timePeriod)
		checkError(err)
		INDplus, err := avAPI.IndicatorPLUS_DI(symbol, interval, timePeriod)
		checkError(err)
		INDminus, err := avAPI.IndicatorMINUS_DI(symbol, interval, timePeriod)
		checkError(err)

		RSI, arr := INDrsi.GetRSI()
		EMA := trendFollowing.EMA(arr, 10)

		rsi := fmt.Sprintf("%f", RSI)
		ema := fmt.Sprintf("%f", EMA)

		_, A := INDadx.Latest()
		_, B := INDplus.Latest()
		_, C := INDminus.Latest()

		adx := strings.Trim(fmt.Sprintf("%f", *A), "{}")
		plus := strings.Trim(fmt.Sprintf("%f", *B), "{}")
		minus := strings.Trim(fmt.Sprintf("%f", *C), "{}")

		outString := "" + symbol + " " + rsi + " " + ema + " " + adx + " " + plus + " " + minus + "\n"
		print(outString)
		file.WriteString(outString)
	}

	file.Close()
}

func getA(alpacaAPI *alpaca.Client, avAPI *alphaVantage.Client) {

}

func main() { /////////////////////////////////////////	  MAIN	 ///////////////////////////////////////////////////////

	// free API key: MHL1PVXKA24TUHYG
	// prem API key: B5NM7SCV8LFLME8Y
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

	AlpClient := alpaca.NewClient(common.Credentials())
	AvClient := alphaVantage.New("B5NM7SCV8LFLME8Y")

	testing := false

	if testing {
		print("testing\n")

		d := time.Now()
		fileName := time.Now().Format("2006-01-02")
		getIndData(AlpClient, AvClient, d, fileName, "weekly", "14", "close")

		return
	}

	timeT := time.Now()

	fileName := timeT.Format("2006-01-02")

	file, _ := os.Create("C:/Trading/" + fileName + "-tipsheet.txt") // i think i am ignoring an error here using the _

	assetsOut, _ := os.Create("C:/Trading/" + fileName + "-assets.txt")

	status := "active"
	assets, err := AlpClient.ListAssets(&status)
	if err != nil {
		panic(err)
	}

	max := 10 // price to keep under
	//min := 5 // nathan
	print("max: ", max, "\n")

	missed := 0

	print("see if price is cheaper than max(", max, "):\n")

	for _, asset := range assets { // check for price to see if affordable
		//note: Alpaca's api allows for 3 calls per second; this is much faster than running through the algorithm
		time.Sleep(300 * time.Millisecond) // this could be closer to 1/3 i think, so long as it is <= 1/3
		lastPrice, err := AlpClient.GetLastQuote(asset.Symbol)
		if err != nil {

			if err.Error() != "resource not found" {
				missed += 1
				continue
			}
			continue
		}

		if lastPrice.Last.AskPrice < float32(max) { // && lastPrice.Last.AskPrice > float32(min) { // nathan
			print(asset.Symbol, " ", lastPrice.Last.AskPrice, "\n")
			assetsOut.WriteString(asset.Symbol + " " + fmt.Sprintf("%f", lastPrice.Last.AskPrice) + "\n")
		}
	}
	print("missed: ", missed, "\n")
	assetsOut.Close()

	dat, err := ioutil.ReadFile("C:/Trading/" + fileName + "-assets.txt") // reads assets form file
	if err != nil {
		panic(err)
	}

	affordableAssets := strings.Split(strings.TrimSpace(string(dat)), "\n")

	print("ticker rsi ema(rsi) adx plusDI minusDI error\n ")
	for i := 0; i < len(affordableAssets); i++ { // this needs to be re-evaluated
		a := strings.Split(affordableAssets[i], " ")

		tip := trendFollowing.GetTrade(false, a[0], AvClient) // need to check position open
		print("\n")

		if tip.Action == "buy" && tip.Side == "long" { // write and read to files
			file.WriteString(a[0] + " " + a[1] + "\n")
		}
	}

	e := file.Close()
	if e != nil {
		panic(e)
	}

	getIndData(AlpClient, AvClient, timeT, fileName, "weekly", "14", "close")
}

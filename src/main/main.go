package main

import (
	"../alphaVantage"
	"../trendFollowing"
	"encoding/json"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"os"
	"sort"
	"strconv"
)

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

type TipSheet struct {
	buy  []string
	sell []string
}

func makeTipSheet(alpacaAPI *alpaca.Client, avAPI *alphaVantage.Client) []alpaca.Asset {
	status := "active"
	assets, err := alpacaAPI.ListAssets(&status)
	if err != nil {
		panic(err)
	}
	print(len(assets))

	//var result TipSheet

	//tips := make([]alpaca.Asset, 0)

	//for i := range assets {
	//	symbol := assets[i].Symbol
	//	//rsiArray := avAPI.GetRSI(symbol) // this needs to be of days
	//	//ema := trendFollowing.EMA(rsiArray, 100)
	//
	//	indicator, err := avAPI.IndicatorRSI(symbol, "daily", "14", "close")
	//	if err != nil {
	//		panic(err) // should change this to manage errors
	//	}
	//	if len(indicator.TechnicalAnalysis) == 0 {
	//		print(symbol, " ", 0, "\n")
	//		continue
	//	}
	//	date, latest := indicator.Latest()
	//	fmt.Printf("%s %s %f %d\n", symbol, date, latest.RSI, len(indicator.TechnicalAnalysis))
	//}

	return assets
}

func test(alpacaAPI *alpaca.Client, avAPI *alphaVantage.Client) {
	ind, err := avAPI.IndicatorRSI("NFLX", "60min", "14", "close")

	if err != nil {
		panic(err)
	}

	m := ind.TechnicalAnalysis

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		print(k, " ", m[k].RSI, "\n")
	}

	//PrettyPrint(ind)
}

func main() {

	// free API key: MHL1PVXKA24TUHYG
	// prem API key: B5NM7SCV8LFLME8Y
	AlpClient := Alpaca("PKXAF267QI7IJV5EUW3L", "p2dCv7ZWkykxY2L7Q3mK6EpLemlAiE5zPxxRd4PR")
	AvClient := alphaVantage.New("B5NM7SCV8LFLME8Y")

	//test(AlpClient, AvClient)

	file, _ := os.Create("stocks.txt")

	temp, _ := os.Create("assets.txt")

	status := "active"
	assets, err := AlpClient.ListAssets(&status)
	if err != nil {
		panic(err)
	}

	for asset := range assets {
		temp.WriteString(assets[asset].Symbol + "\n")
	}
	temp.Close()

	for i := 0; i < len(assets); i++ {
		a := assets[i]
		ind, err := AvClient.IndicatorRSI(a.Symbol, "60min", "14", "close")
		if err != nil { // write error to file
			print(err.Error(), "\n")
			continue
		}

		// check for over flow / running out of calls
		print(len(ind.TechnicalAnalysis), " ")
		if len(ind.TechnicalAnalysis) < 1 { // dashes are not friendly
			PrettyPrint(ind.TechnicalAnalysis)
			continue
		}

		latest, array := ind.GetRSI()
		tip, ema := trendFollowing.GetTrade(latest, array, false)
		s := fmt.Sprintf("%f", latest)
		e := fmt.Sprintf("%f", ema)

		if tip.Action == "buy" {
			file.WriteString(a.Symbol + " " + s + " " + e + "\n")

		}
		fmt.Println(a.Symbol + " " + s + " " + e + " " + strconv.Itoa(i))
	}

	e := file.Close()
	if e != nil {
		panic(e)
	}
}

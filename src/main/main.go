package main

import (
	"../alphaVantage"
	"encoding/json"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
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

func makeTipSheet(alpacaAPI *alpaca.Client, avAPI *alphaVantage.Client) TipSheet {
	status := "active"
	assets, err := alpacaAPI.ListAssets(&status)
	if err != nil {
		panic(err)
	}
	print(len(assets))

	var result TipSheet

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

	return result
}

func test(alpacaAPI *alpaca.Client, avAPI *alphaVantage.Client) {

}

func main() {

	ALP_client := Alpaca("PKXAF267QI7IJV5EUW3L", "p2dCv7ZWkykxY2L7Q3mK6EpLemlAiE5zPxxRd4PR")
	AV_client := alphaVantage.New("MHL1PVXKA24TUHYG")

	test(ALP_client, AV_client)

	//makeTipSheet(ALP_client, AV_client)

	//indicator, _ := c.IndicatorRSI("NYMT", "daily", "14", "open")
	//indicator, _ := c.IndicatorRSI("NYMT", "5min", "21", "open")
	//day, _ := c.IndicatorRSI("NYMT", "daily", "21", "open")
	//week, err := c.IndicatorRSI("MSFT", "weekly", "21", "open")

	//tipSheet := trendFollowing.GetTipSheet()
	//
	//for tip := range tipSheet {
	//	trade := trendFollowing.GetTrade(tip)
	//	if trade != nil {
	//		rsiArray := alphaVantage.GetRSI(tip.ticker)
	//		rsi :=
	//		api.PlaceOrder(rsiArray)
	//	}
	//}
}

package main

import (
	"../alphaVantage"
	"../trendFollowing"
	"encoding/json"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
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
	symbol := "F"
	args := []string{
		//"1min open",
		//"1min close",
		//"1min low",
		//"1min high",
		//
		//"30min open",
		//"30min close",
		//"30min low",
		//"30min high",
		//
		//"60min open",
		//"60min close",
		//"60min low",
		//"60min high",
		//
		//"daily open",
		//"daily close",
		//"daily low",
		//"daily high",
		//
		//"weekly open",
		//"weekly close",
		//"weekly low",
		//"weekly high",
		//
		//"monthly open",
		//"monthly close",
		//"monthly low",
		//"monthly high",

		"1min open",
		"30min open",
		"60min open",
		"daily open",
		"weekly open",
		"monthly open",

		"1min close",
		"30min close",
		"60min close",
		"daily close",
		"weekly close",
		"monthly close",

		"1min low",
		"30min low",
		"60min low",
		"daily low",
		"weekly low",
		"monthly low",

		"1min high",
		"30min high",
		"60min high",
		"daily high",
		"weekly high",
		"monthly high",
	}

	inds := make(map[string]*alphaVantage.IndicatorRSI)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	for i := range args { // get all of the indicators
		a := strings.Split(args[i], " ")
		print(a[0] + " " + a[1] + "\n")
		ind, err := avAPI.IndicatorRSI(symbol, a[0], "14", a[1])
		if err != nil {
			print(err.Error() + "\n")
			continue
		}
		inds[args[i]] = ind
	}

	for i := range args { // print all of the indicators as individuals
		date, rsi := inds[args[i]].Latest()
		a := strings.Split(args[i], " ")
		fmt.Fprintln(w, fmt.Sprintf("%f", rsi.RSI)+"\t"+a[0]+"\t"+a[1]+"\t"+date)
	}

	T := []string{
		"1min",
		"30min",
		"60min",
		"daily",
		"weekly",
		"monthly",
	}

	P := []string{
		" high",
		" low",
		" close",
	}

	for _, t := range T { // calculate lhc3
		total := 0.0
		for _, p := range P {
			param := t + p
			_, rsi := inds[param].Latest()
			total += rsi.RSI
		}
		total /= 3

		fmt.Fprintln(w, fmt.Sprintf("%f", total)+"\t"+t+"\t"+"hlc3")
	}

	P = []string{
		" open",
		" high",
		" low",
		" close",
	}

	for _, t := range T { // calculate lhc3
		total := 0.0
		for _, p := range P {
			param := t + p
			_, rsi := inds[param].Latest()
			total += rsi.RSI
		}
		total /= 4

		fmt.Fprintln(w, fmt.Sprintf("%f", total)+"\t"+t+"\t"+"ohlc4")
	}

	w.Flush()
}

func testADX(alpacaAPI *alpaca.Client, avAPI *alphaVantage.Client) {
	x, y := avAPI.IndicatorADX("IBM", "daily", "10")
	PrettyPrint(x)
	PrettyPrint(y)
	date, val := x.Latest()
	PrettyPrint(date)
	PrettyPrint(val)
}

func main() {

	// free API key: MHL1PVXKA24TUHYG
	// prem API key: B5NM7SCV8LFLME8Y
	API_KEY := "PKXAF267QI7IJV5EUW3L"
	API_SECRET := "p2dCv7ZWkykxY2L7Q3mK6EpLemlAiE5zPxxRd4PR"
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

	AccountPercentPerShare := 0.0001 // find way to normalize this or to stick it with a range

	testing := true

	if testing {
		print("testing\n")
		//params := polygon.HistoricTicksV2Params{
		//	Timestamp:      0,
		//	TimestampLimit: 0,
		//	Reverse:        false,
		//	Limit:          0,
		//}
		//x, y := polygon.Client{}.GetHistoricAggregatesV2("F", 1, "close", time.Now(), time.Now(), false)
		//if y != nil {
		//	print(y.Error(), "\n")
		//}
		//PrettyPrint(x)
		//test(AlpClient, AvClient)
		testADX(AlpClient, AvClient)

		//x, y := AlpClient.GetAccount()
		//PrettyPrint(x)
		//PrettyPrint(y)

		return
	}

	fileName := time.Now().Format("02-Jan-2006")

	file, _ := os.Create(fileName + ".txt")

	temp, _ := os.Create("assets.txt")

	status := "active"
	assets, err := AlpClient.ListAssets(&status)
	if err != nil {
		panic(err)
	}

	act, err := AlpClient.GetAccount()
	if err != nil {
		panic(err)
	}

	eqt, _ := act.Equity.Float64()
	max := eqt * AccountPercentPerShare

	for asset := range assets { // check for price to see if affordable
		// note: Alpaca's api allows for 3 calls per second
		print(max, "\n")
		temp.WriteString(assets[asset].Symbol + "\n")
	}
	temp.Close()

	for i := 0; i < len(assets); i++ {
		a := assets[i]
		ind, err := AvClient.IndicatorRSI(a.Symbol, "weekly", "14", "close")
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

		if tip.Action == "buy" && tip.Side == "long" {
			file.WriteString(a.Symbol + " " + s + " " + e + "\n")
			fmt.Println(a.Symbol + " " + s + " " + e + " " + strconv.Itoa(i))
		} else {
			fmt.Println(a.Symbol)
		}
	}

	e := file.Close()
	if e != nil {
		panic(e)
	}

}

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
			print(err.Error() + " RSI\n")
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

func checkError(e error) {
	if e != nil {
		PrettyPrint(e.Error())
	}
}

func getData(alpacaAPI *alpaca.Client, avAPI *alphaVantage.Client, date string, fileName string,
	interval string, timePeriod string, seriesType string) {
	// date in "YYYY-MM-DD" format or "2006-01-02"
	//interval := "weekly"
	//timePeriod := "14"
	//seriesType := "close"
	dat, err := ioutil.ReadFile(fileName) // reads assets form file // CHANGE TO BE DATE UNIQUE
	if err != nil {
		panic(err)
	}
	affordableAssets := strings.Split(string(dat), "\r\n") // this is a bug that needs ot be squashed

	outputName := date + "_" + interval + "_" + timePeriod + "_" + seriesType + ".txt"

	file, e := os.Open("C:/Trading/" + outputName)
	if e != nil {
		file, _ = os.Create("C:/Trading/" + outputName)
	}

	for _, symbol := range affordableAssets {

		t, _ := time.Parse("2006-01-02", "2020-05-17")
		param := alpaca.ListBarParams{
			Timeframe: "day",
			StartDt:   &t,
			EndDt:     nil,
			Limit:     nil,
		}

		//o := float32(0.0)
		c := float32(0.0)
		ag, erro := alpacaAPI.GetSymbolBars(symbol, param)
		if erro == nil {
			for _, i := range ag {
				if time.Unix(i.Time, 0).Format("2006-01-02") == t.Format("2006-01-02") {
					c = i.Close
					//print(t.Format("2006-01-02"),"\n")
					t = t.Add(time.Hour * 24)
					break
				}
			}
			//for _,i := range ag {
			//	if time.Unix(i.Time,0).Format("2006-01-02") == t.Format("2006-01-02") {
			//		//print(t.Format("2006-01-02"),"\n")
			//		//PrettyPrint(i)
			//		o = i.Open
			//		break
			//	}
			//}
			// find the following open value
		}
		close := fmt.Sprintf("%f", c)
		//open := fmt.Sprintf("%f", o)

		INDrsi, err := avAPI.IndicatorRSI(symbol, interval, timePeriod, seriesType)
		checkError(err)
		INDadx, err := avAPI.IndicatorADX(symbol, interval, timePeriod)
		checkError(err)
		INDplus, err := avAPI.IndicatorPLUS_DI(symbol, interval, timePeriod)
		checkError(err)
		INDminus, err := avAPI.IndicatorMINUS_DI(symbol, interval, timePeriod)
		checkError(err)

		rsi := fmt.Sprintf("%f", INDrsi.TechnicalAnalysis[date].RSI)
		adx := fmt.Sprintf("%f", INDadx.TechnicalAnalysis[date].ADX)
		plus := fmt.Sprintf("%f", INDplus.TechnicalAnalysis[date].PLUS_DI)
		minus := fmt.Sprintf("%f", INDminus.TechnicalAnalysis[date].MINUS_DI)

		outString := "" + symbol + " " + rsi + " " + adx + " " + plus + " " + minus + " " + close + "\n"
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

	//x := polygon.NewClient(common.Credentials())
	//x.GetStockExchanges()

	//AccountPercentPerShare := 0.0001 // TODO: find way to normalize this or to stick it with a range

	testing := false

	if testing {
		print("testing\n")
		//intervals := []string{
		//	"weekly",
		//}
		//
		//tPs := []string{
		//	"10",
		//	"11",
		//	"12",
		//	"13",
		//	"14",
		//	"15",
		//	"16",
		//}
		//
		//sTs := []string{
		//	"open",
		//	"high",
		//	"low",
		//	"close",
		//}
		//
		//for _, interval := range intervals{
		//	for _, timePeriod := range tPs{
		//		for _, seriesType := range sTs {
		//			print(interval, " ", timePeriod, " ", seriesType, "\n")
		//			getData(AlpClient, AvClient,"2020-05-15", "17-May-2020.txt", interval, timePeriod, seriesType)
		//		}
		//	}
		//}

		os.Create("C:/Trading/hi_there.txt")

		return
	}

	fileName := time.Now().Format("02-Jan-2006")

	file, _ := os.Create("C:/Trading/" + fileName + ".txt") // i think i am ignoring an error here using the _

	temp, _ := os.Create("C:/Trading/" + fileName + "assets.txt")

	status := "active"
	assets, err := AlpClient.ListAssets(&status)
	if err != nil {
		panic(err)
	}

	//act, err := AlpClient.GetAccount() // commented out for Nathan's testing
	//if err != nil {
	//	panic(err)
	//}
	//
	//eqt, _ := act.Equity.Float64()
	//max := eqt * AccountPercentPerShare
	// TODO: better equation for this ^^^

	max := 10 // this is for testing purposes
	print("max: ", max, "\n")

	missed := 0 // this is for testing

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

		if lastPrice.Last.AskPrice < float32(max) {
			print(asset.Symbol, " ", lastPrice.Last.AskPrice, "\n")
			temp.WriteString(asset.Symbol + " " + fmt.Sprintf("%f", lastPrice.Last.AskPrice) + "\n")
		}

		//temp.WriteString(asset.Symbol + "\n")
	}
	print("missed: ", missed, "\n")
	temp.Close()

	dat, err := ioutil.ReadFile("C:/Trading/" + fileName + "assets.txt") // reads assets form file
	if err != nil {
		panic(err)
	}
	affordableAssets := strings.Split(string(dat), "\n")

	for i := range affordableAssets { // assets from file
		affordableAssets[i] = strings.Split(affordableAssets[i], " ")[0] // this should be made cleaner
	}

	// split

	print("ticker rsi ema(rsi) adx plusDI minusDI error\n ")
	for i := 0; i < len(affordableAssets); i++ { // this needs to be re-evaluated
		a := affordableAssets[i]
		//print(a, "\n")

		tip := trendFollowing.GetTrade(false, a, AvClient) // need to check position open
		print("\n")

		if tip.Action == "buy" && tip.Side == "long" { // write and read to files
			file.WriteString(a + "\n")
			//fmt.Println("\n" + a + " " + e + " " + strconv.Itoa(i))
		} else {
			//fmt.Print(" " + a)
		}
	}

	e := file.Close()
	if e != nil {
		panic(e)
	}

	d := time.Now().Format("2006-01-02")
	getData(AlpClient, AvClient, d, fileName, "weekly", "14", "close")
}

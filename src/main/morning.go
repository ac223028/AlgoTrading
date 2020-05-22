package main

import (
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

//func PrettyPrint(v interface{}) (err error) {
//	b, err := json.MarshalIndent(v, "", "  ")
//	if err == nil {
//		fmt.Println(string(b))
//	}
//	return
//}

func buyStocks(alpacaAPI *alpaca.Client, account *alpaca.Account, tips []string, times int) {
	print("Started Buying\n")
	for j := 0; j < times; j++ {
		for _, tip := range tips {
			ticker := strings.Split(tip, " ")[0]

			request := alpaca.PlaceOrderRequest{
				AccountID:   account.ID,
				AssetKey:    &ticker,
				Qty:         decimal.NewFromInt(int64(1)),
				Side:        "buy",
				Type:        "market",
				TimeInForce: "day",
			}
			time.Sleep(330 * time.Millisecond)
			_, err := alpacaAPI.PlaceOrder(request)

			if err != nil {
				PrettyPrint(err.Error())
				continue
			}

			print("Bought ", ticker, "\n")
		}
	}
}

func sumTipPrices(tips []string) float64 {
	result := 0.0
	for _, tip := range tips {
		p := strings.Split(tip, " ")[1]
		//PrettyPrint(p)
		price, err := strconv.ParseFloat(p, 64)
		if err != nil {
			panic(err)
		}
		result += price
	}
	return result
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

	AlpClient := alpaca.NewClient(common.Credentials())

	testing := false

	if testing {

		return
	}

	t := time.Now()
	fileHead := t.Format("2006-01-02")

	tipsheetFileName := fileHead + "-tipsheet"

	dat, err := ioutil.ReadFile("C:/Trading/" + tipsheetFileName + ".txt") // reads assets form file
	if err != nil {
		dat, err = ioutil.ReadFile(tipsheetFileName + ".txt") // reads assets form file
		if err != nil {
			panic(err)
		}
	}

	tips := strings.Split(strings.TrimSpace(string(dat)), "\n") // TODO: sort by prices low to high
	PrettyPrint(tips)

	sum := sumTipPrices(tips)
	print("Total cost of single set of stock: ", sum, "\n")

	account, err := AlpClient.GetAccount()
	if err != nil {
		panic(err)
	}

	AlpClient.CancelAllOrders()   // this frees up as much cap as pos
	AlpClient.CloseAllPositions() // this sells all stocks

	//print("Waiting for server to catch up\n")
	//time.Sleep(time.Second)

	//bp, _ := account.BuyingPower.Float64() // not sure why this doesnt work; uses leverage of which i have none
	bp, _ := account.Cash.Float64() // this assumes i use no leverage
	print("Buying Power: ", bp, "\n")

	mult := int(math.Floor(bp / sum))
	print("Will buy ", mult, " set(s) of stocks\n")

	//PrettyPrint(mult)
	//PrettyPrint(bp)

	print("Buying stocks\n")
	buyStocks(AlpClient, account, tips, mult)
}

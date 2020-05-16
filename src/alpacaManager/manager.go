package alpacaManager

import "github.com/alpacahq/alpaca-trade-api-go/alpaca"

type AlpManager struct {
	Client alpaca.Client
}

func main() {

}

// needs ways to determine how many shares to buy and if my account can afford it
//     split the account into quarters (when determining number of shares to buy, floor divide)
// needs a way to find if a share is sell-able or not
//     could be a function
// a queue for buying?
//     so if a spot opens up, a stock can still be bought
// should wait for data client and integrate them here

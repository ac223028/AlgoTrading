package trendFollowing

import (
	"../alphaVantage"
	"encoding/json"
	"fmt"
)

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

type Trade struct {
	Action string
	Side   string
}

func GetTrade(openPosition bool, ticker string, AvClient *alphaVantage.Client) (Trade, float32) { // need to incorporate time as an input
	indRSI, err := AvClient.IndicatorRSI(ticker, "weekly", "14", "close")
	rsi, rsiArray := indRSI.GetRSI()

	ema := EMA(rsiArray, 10) // this N does not have to be 10
	result := Trade{"", ""}

	if err != nil { // write error to file
		print(err.Error(), "\n")
		return result, -1 // there needs to be a better solution to this
	}

	// check for over flow / running out of calls
	print(len(indRSI.TechnicalAnalysis), " ")
	if len(indRSI.TechnicalAnalysis) < 1 { // dashes are not friendly
		PrettyPrint(indRSI.TechnicalAnalysis) // there needs to be a better check of this
		return result, -1
	}

	if ema < 40 || ema > 60 { // there are ways to figure out if trend is going up
		if rsi > ema { // also need to figure out how much to buy
			// run the check on momentum here
			if !openPosition {
				return Trade{"buy", "long"}, ema // open long position
			} else {
				return Trade{"sell", "short"}, ema //     close the short position
			}
		}
		if rsi < ema {
			// run the check on momentum here
			if !openPosition {
				return Trade{"buy", "short"}, ema // open a short position
			} else {
				return Trade{"sell", "long"}, ema // close the long position
			}
		}
	}
	return result, ema
}

func checkMomentum(AvClient *alphaVantage.Client, ticker string) float64 {
	// returns +float for upward, -float for downward, zero otherwise
	interval := ""
	timePeriod := ""

	indADX, err := AvClient.IndicatorADX(ticker, interval, timePeriod)
	if err != nil {
		panic(err)
	}

	_, latestADX := indADX.Latest()

	if latestADX.ADX >= 50 {
		indPLUS_DI, err := AvClient.IndicatorPLUS_DI(ticker, interval, timePeriod)
		if err != nil {
			panic(err)
		}
		indMINUS_DI, err := AvClient.IndicatorMINUS_DI(ticker, interval, timePeriod)
		if err != nil {
			panic(err)
		}

		_, latestPLUS := indPLUS_DI.Latest()
		_, latestMINUS := indMINUS_DI.Latest()
		pos := latestPLUS.PLUS_DI
		neg := latestMINUS.MINUS_DI

		if pos > neg {
			return pos - neg
		} else {
			return pos - neg
		}
	}

	return 0
}

func SMA(data []float32) float32 {
	result := float32(0)

	for i := range data {
		result += data[i]
	}

	return result / float32(len(data))
}

func EMA(data []float32, N int, kOptional ...float32) float32 { // may want to look at Wilder's moving average
	// this is one way of calculating it but requires the input to be entered in reverse [oldest:today]
	// might convert to queue
	// a loop is faster

	k := float32(0)
	if len(kOptional) > 0 {
		k = kOptional[0]
	} else {
		k = float32(2) / float32(N+1)
	}

	if len(data) <= N { // this is the base case
		result := SMA(data)
		return result
	}

	y := EMA(data[1:], N, k)
	result := k*(data[0]-y) + y // based on the excel sheet
	//result := (data[0] - y * k) + y // based on the research paper (way off...)
	//result := data[0] * k + y * (1-k) // based on investopedia (almost the same as excel)
	return result
}

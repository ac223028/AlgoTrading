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

func GetTrade(openPosition bool, ticker string, AvClient *alphaVantage.Client) Trade {
	interval := "weekly" // these params can be changed
	timePeriod := "14"
	seriesType := "close"
	indRSI, err := AvClient.IndicatorRSI(ticker, interval, timePeriod, seriesType)

	print(ticker)

	result := Trade{"", ""}

	if err != nil { // write error to file
		//print(" ", err.Error(), " RSI")
		return result // TODO: there needs to be a better solution to this
	}

	//PrettyPrint(indRSI)
	rsi, rsiArray := indRSI.GetRSI()
	print(" ", rsi)

	ema := EMA(rsiArray, 10) // TODO: this N does not have to be 10
	print(" ", ema)

	if len(indRSI.TechnicalAnalysis) < 1 { // dashes are not friendly
		//PrettyPrint(indRSI.TechnicalAnalysis) // TODO: there needs to be a better check of this
		return result
	}

	if ema < 40 || ema > 60 { // there are ways to figure out if trend is going up
		if rsi > ema { // also need to figure out how much to buy
			// run the check on momentum here: advancing
			if checkMomentum(AvClient, ticker, interval, timePeriod) > 0 { // think about error as a means of correction
				if !openPosition {
					return Trade{"buy", "long"} // open long position
				} else {
					return Trade{"sell", "short"} // close the short position
				}
			}
		}
		if rsi < ema {
			// run the check on momentum here: declining
			if checkMomentum(AvClient, ticker, interval, timePeriod) < 0 {
				if !openPosition {
					return Trade{"buy", "short"} // open a short position
				} else {
					return Trade{"sell", "long"} // close the long position
				}
			}
		}
	}
	return result
}

func checkMomentum(AvClient *alphaVantage.Client, ticker string, interval string, timePeriod string) float64 {
	// returns +float for upward, -float for downward, zero otherwise

	indADX, err := AvClient.IndicatorADX(ticker, interval, timePeriod)
	if err != nil {
		//print(" ", err.Error(), " ADX")
		return 0
	}

	_, latestADX := indADX.Latest()
	print(" ", latestADX.ADX)

	if latestADX.ADX >= 45 {
		indPLUS_DI, err := AvClient.IndicatorPLUS_DI(ticker, interval, timePeriod)
		if err != nil {
			//print(" ", err.Error(), " +DI")
			return 0 // there needs to be a better solution to this
		}
		indMINUS_DI, err := AvClient.IndicatorMINUS_DI(ticker, interval, timePeriod)
		if err != nil {
			//print(" ", err.Error(), " -DI")
			return 0 // there needs to be a better solution to this
		}

		_, latestPLUS := indPLUS_DI.Latest()
		_, latestMINUS := indMINUS_DI.Latest()
		pos := latestPLUS.PLUS_DI
		neg := latestMINUS.MINUS_DI

		print(" ", pos)
		print(" ", neg)

		return pos - neg
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

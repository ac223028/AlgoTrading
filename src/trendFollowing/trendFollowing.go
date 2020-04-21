package trendFollowing

import (
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
	action string
	side   string
}

func GetTrade(rsi float32, rsiArray []float32, openPosition bool) Trade { // need to incorporate time as an input
	ema := EMA(rsiArray, 100) // this N does not have to be 10
	result := Trade{"nil", "nil"}

	if ema < 40 || ema > 60 {
		if rsi > ema { // also need to figure out how much to buy
			if openPosition {
				return Trade{"buy", "long"} // open long position
			} else {
				return Trade{"sell", "short"} //     close the short position
			}
		}
		if rsi < ema {
			if !openPosition {
				return Trade{"buy", "short"} // open a short position
			} else {
				return Trade{"sell", "long"} // close the long position
			}
		}
	}
	return result
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

func GetTipSheet() int {
	return 0
}

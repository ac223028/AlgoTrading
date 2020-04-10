package trendFollowing

import (
	"time"
)

type trade struct {
	ticker string
	price  float32
	date   time.Time
}

func getTrade(rsi float32, rsiArray []float32) trade { // need to incorporate time as an input
	ema := EMA(rsiArray, 10) // this N does not have to be 10
	result := trade(nil)     // this makes sense for some reason

	if ema < 40 || ema > 60 {
		if rsi > ema {
			// if no position is open
			//     open a long position
			// else if a short position is open
			//     close the short position
			return result
		}
		if rsi < ema {
			// if no position is open
			//     open a short position
			// else if long position is open
			//     close the long position
			return result
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

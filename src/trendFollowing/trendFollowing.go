package trendFollowing

import (
)

func SMA(data []float32) float32{
	result := float32(0)

	for i := range data{
		result += data[i]
	}

	return result / float32(len(data))
}

func EMA(data []float32, N int, kOptional ...float32) float32{  // may want to look at Wilder's moving average
	// this is one way of calculating it but requires the input to be entered in reverse [oldest:today]
	// might convert to queue
	// a loop is faster

	k := float32(0)
	if len(kOptional) > 0 {
		k = kOptional[0]
	} else {
		k = float32(2) / float32(N + 1)
	}

	if len(data) <= N { // this is the base case
		result := SMA(data)
		return result
	}

	y := EMA(data[1:], N, k)
	result := k * (data[0] - y) + y  // based on the excel sheet
	//result := (data[0] - y * k) + y // based on the research paper (way off...)
	//result := data[0] * k + y * (1-k) // based on investopedia (almost the same as excel)
	return result
}

func avgGainLoss(data []float32) (float32, float32) {
	return 0.0, 0.0
}

func RSI(data []float32, N int, averages ...float32) float32 { // this may just be pulled from alpha vantage
	// Alpha vantage API key: MHL1PVXKA24TUHYG
	// the input should be a set of +/- values for bar to bar change

	//if len(data) <= N {  // this should return an error
	//	return float32(-1)
	//}

	var aGain float32
	var aLoss float32

	if len(averages) > 0 {
		aGain, aLoss = averages[0], averages[1]
	} else {
		aGain, aLoss = avgGainLoss(data[:N])
	}

	return aGain - aLoss
}
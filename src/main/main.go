package main

import (
	"../alphaVantage"
	"fmt"
)

func main() {
	//var data []float32
	//
	//data = []float32{ // may want to change this to floats
	//	22.17, 22.40, 23.10, 22.68, 23.33, 23.10, 23.19, 23.65,
	//	23.87, 23.82, 23.63, 23.95, 23.83, 23.75, 24.05, 23.36,
	//	22.61, 22.38, 22.39, 22.15, 22.29, 22.24, 22.43, 22.23,
	//	22.13, 22.18, 22.17, 22.08, 22.19, 22.27} // checked using the excel data form stockcharts.com
	//
	//start := time.Now()
	//x := trendFollowing.EMA(data, 10)
	//elapsed := time.Since(start)
	//print("EMA ", x, " in ", float32(elapsed), "\n")
	//
	//data = []float32{
	//	43.13, 42.66, 43.42, 44.57, 44.22, 44.18, 44.03, 45.35,
	//	45.78, 46.45, 45.71, 46.25, 46.21, 45.64, 46.22, 46.41,
	//	46.03, 46.00, 46.28, 46.28, 45.61, 46.03, 45.89, 46.08,
	//	45.84, 45.42, 45.10, 44.83, 44.33, 43.61, 44.15, 44.09,
	//	44.34,
	//}
	//
	//start = time.Now()
	//y := trendFollowing.RSI(data, 14)
	//elapsed = time.Since(start)
	//print("RSI ", y, " in ", float32(elapsed))

	c := alphaVantage.New("MHL1PVXKA24TUHYG")
	indicator, _ := c.IndicatorRSI("NYMT", "dail", "14", "close")

	//time, _ := time.Parse("2006-01-02", "2020-02-24")
	//fmt.Println(indicator)
	//fmt.Println(err)
	//fmt.Println(time)

	//date, RSI := indicator.Latest() // this can generate key errors if i do not formulate my key correctly

	//print(date, " ", RSI.RSI, "\n")
	_map := indicator.TechnicalAnalysis

	for key, value := range _map {
		fmt.Println(key, value)
	}
}

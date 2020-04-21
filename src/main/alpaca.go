package main

import (
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"os"
)

func Alpaca(key string, secretKey string) *alpaca.Client {
	os.Setenv(common.EnvApiKeyID, key)
	os.Setenv(common.EnvApiSecretKey, secretKey)

	//fmt.Printf("Running w/ credentials [%v %v]\n", common.Credentials().ID, common.Credentials().Secret)

	alpaca.SetBaseUrl("https://paper-api.alpaca.markets")

	api := alpaca.NewClient(common.Credentials())
	return api
}

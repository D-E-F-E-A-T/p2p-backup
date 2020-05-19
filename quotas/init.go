package quotas

import (
	"log"
	"net/http"
	"time"

	coingecko "github.com/superoo7/go-gecko/v3"

	"github.com/dmitriy-vas/p2p/handler/middleware"
)

var (
	toAbbreviation = map[string]string{
		"bitcoin":      "btc",
		"bitcoin-cash": "bch",
		"ethereum":     "eth",
		"litecoin":     "ltc",
		"prizm":        "pzm",
		"ripple":       "xrp",
		"waves":        "wvs",
	}
	currenciesID = []string{
		"bitcoin",
		"bitcoin-cash",
		"ethereum",
		"litecoin",
		"prizm",
		"ripple",
		"waves",
	}
	currenciesVS = []string{
		"bch",
		"btc",
		"eth",
		"eur",
		"gbp",
		"ltc",
		"rub",
		"usd",
	}
	fiatMapping = map[string]interface{}{
		"eur": nil,
		"gbp": nil,
		"rub": nil,
		"usd": nil,
	}
)

var Quotas map[string]map[string]float32

func init() {
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}
	client := coingecko.NewClient(httpClient)
	go UpdateQuotas(client)
}

func UpdateQuotas(client *coingecko.Client) {
	for {
		results, err := client.SimplePrice(currenciesID, currenciesVS)
		if err != nil {
			log.Printf("Error getting quotas: %v", err)
		}
		Quotas = ConvertMapToQuotas(*results)
		time.Sleep(time.Minute * 10)
	}
}

func ConvertMapToQuotas(m map[string]map[string]float32) map[string]map[string]float32 {
	quotas := make(map[string]map[string]float32)
	for currency, listing := range m {
		abbreviation := toAbbreviation[currency]
		if _, exists := listing[abbreviation]; exists {
			delete(listing, abbreviation)
		}
		quotas[abbreviation] = listing
	}
	for currency, listing := range quotas {
		for key, cost := range listing {
			if _, exists := fiatMapping[key]; exists {
				if _, exists = quotas[key]; !exists {
					quotas[key] = make(map[string]float32)
				}
				quotas[key][currency] = 1 / cost
			}
		}
		for currency2, listing2 := range quotas {
			if _, exists := fiatMapping[currency]; currency == currency2 || exists {
				continue
			}
			if _, exists := listing[currency2]; exists {
				continue
			}
			quotas[currency][currency2] = listing["usd"] / listing2["usd"]
		}
	}
	return quotas
}

func GetCostWithProfit(currency1, currency2 interface{}, profit int16) float32 {
	switch currency1.(type) {
	case uint8:
		currency1 = middleware.Currency(currency1.(uint8))
	}
	switch currency2.(type) {
	case uint8:
		currency2 = middleware.Currency(currency2.(uint8))
	}
	quota := Quotas[currency1.(string)][currency2.(string)]
	return quota + (quota/100.0) * float32(profit)
}

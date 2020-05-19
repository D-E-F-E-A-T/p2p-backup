package crypto

import (
	"github.com/spf13/viper"
)

func init() {
	bitcoinWrapper := NewDefaultWrapper(
		viper.GetString("crypto.bitcoin.host"),
		viper.GetString("crypto.bitcoin.token"),
		"btc",
	)
	Wrappers["btc"] = bitcoinWrapper
}

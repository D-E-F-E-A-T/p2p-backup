package crypto

import (
	"github.com/spf13/viper"
)

func init() {
	bitcoinCashWrapper := NewDefaultWrapper(
		viper.GetString("crypto.bcash.host"),
		viper.GetString("crypto.bcash.token"),
		"bch",
	)
	Wrappers["bch"] = bitcoinCashWrapper
}

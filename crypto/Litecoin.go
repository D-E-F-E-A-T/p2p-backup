package crypto

import (
	"github.com/spf13/viper"
)

func init() {
	litecoinWrapper := NewDefaultWrapper(
		viper.GetString("crypto.litecoin.host"),
		viper.GetString("crypto.litecoin.token"),
		"ltc",
	)
	Wrappers["ltc"] = litecoinWrapper
}

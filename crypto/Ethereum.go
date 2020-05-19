package crypto

import (
	"github.com/spf13/viper"
)

func init() {
	ethereumWrapper := NewDefaultWrapper(
		viper.GetString("crypto.ethereum.host"),
		viper.GetString("crypto.ethereum.token"),
		"eth",
	)
	Wrappers["eth"] = ethereumWrapper
}

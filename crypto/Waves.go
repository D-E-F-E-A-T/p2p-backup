package crypto

import (
	"github.com/spf13/viper"
)

func init() {
	wavesWrapper := NewDefaultWrapper(
		viper.GetString("crypto.waves.host"),
		viper.GetString("crypto.waves.token"),
		"wvs",
	)
	Wrappers["wvs"] = wavesWrapper
}

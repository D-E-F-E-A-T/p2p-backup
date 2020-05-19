package wire

import (
	"math/big"

	"github.com/spf13/viper"
)

const (
	RippleAbbreviation = "xrp"
)

func init() {
	if viper.Get("ripple") == nil {
		return
	}
	//client, err := websockets.NewRemote("test")
	//if err != nil {
	//	log.Panicf("Error loading %s: %v", RippleAbbreviation, err)
	//}
	//subRes, err := client.Subscribe(false,
	//	true,
	//	false,
	//	false)
	//ICryptoMap[RippleAbbreviation] = Ripple{
	//
	//}
}

type Ripple struct {
}

func (w Ripple) GetBalance(account string) (balance *big.Float, err error) {
	panic("implement me")
}

func (w Ripple) SendTransaction(from string, to string, amount *big.Float) (id string, err error) {
	panic("implement me")
}

func (w Ripple) CreateAccount(account string) (address string, err error) {
	panic("implement me")
}

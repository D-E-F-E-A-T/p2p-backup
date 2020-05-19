package wire

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/Dmitriy-Vas/go-cryptoconv"
	"github.com/go-pg/pg"
	"github.com/mr-tron/base58"
	"github.com/spf13/viper"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"

	"github.com/dmitriy-vas/node/models"
	"github.com/dmitriy-vas/node/postgres"
)

const (
	WavesAbbreviation  = "wvs"
	TransactionFee     = 100000
	TransactionPayment = 4
	AddressSize        = 26
)

func init() {
	if viper.Get("waves") == nil {
		return
	}
	client, err := client.NewClient(client.Options{
		BaseUrl: viper.GetString("waves.host"),
		Client:  nil,
		ApiKey:  viper.GetString("waves.key"),
	})
	if err != nil {
		log.Panicf("Error loading Waves: %v", err)
	}
	ICryptoMap[WavesAbbreviation] = Waves{
		client,
		context.Background(),
	}
}

type Waves struct {
	*client.Client
	context.Context
}

func (w Waves) GetBalance(account string) (balance *big.Float, err error) {
	if raw, err := base58.Decode(account); err != nil || len(raw) < AddressSize {
		out, err := postgres.Database.SearchWavesAccount(account)
		if err != nil {
			return balance, err
		}
		account = out.Address
	}

	address, err := proto.NewAddressFromString(account)
	if err != nil {
		return balance, err
	}

	addressBalance, _, err := w.Addresses.Balance(w.Context, address)
	if err != nil {
		return balance, err
	}

	log.Printf("%+v", *addressBalance)

	balanceString := strconv.FormatUint(addressBalance.Balance, 10)
	balance = cryptoconv.From(balanceString, WavesAbbreviation)
	return balance, nil
}

func (w Waves) SendTransaction(from string, to string, amount *big.Float) (id string, err error) {
	type SignTransaction struct {
		Type      int     `json:"type"`
		Sender    string  `json:"sender"`
		Recipient string  `json:"recipient"`
		Amount    float64 `json:"amount"`
		Fee       int     `json:"fee"`
	}

	amountWavelets := cryptoconv.To(amount.String(), WavesAbbreviation)
	amountFloat, _ := amountWavelets.Float64()

	raw, err := json.Marshal(SignTransaction{
		Type:      TransactionPayment,
		Sender:    from,
		Recipient: to,
		Amount:    amountFloat,
		Fee:       TransactionFee,
	})
	if err != nil {
		return id, err
	}

	request, _ := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/transactions/sign", w.GetOptions().BaseUrl),
		bytes.NewReader(raw))
	request.Header.Set("X-API-Key", w.GetOptions().ApiKey)

	var buf bufio.ReadWriter
	if _, err := w.Client.Do(w.Context, request, &buf); err != nil {
		return id, err
	}

	request2, _ := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/transactions/broadcast", w.GetOptions().BaseUrl),
		buf)
	_, err = w.Client.Do(w.Context, request2, nil)
	return id, nil
}

func (w Waves) CreateAccount(account string) (address string, err error) {
	out, err := postgres.Database.SearchWavesAccount(account)
	if err != nil && err != pg.ErrNoRows {
		return address, err
	} else if err == nil {
		return out.Address, nil
	}

	request, _ := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/addresses", w.GetOptions().BaseUrl),
		nil)
	request.Header.Set("X-API-Key", w.GetOptions().ApiKey)

	type CreateAccountResponse struct {
		Address string `json:"address"`
	}
	var response CreateAccountResponse
	if _, err := w.Client.Do(w.Context, request, &response); err != nil {
		return address, err
	}
	address = response.Address

	if err := postgres.Database.AddWavesAccount(&models.Waves{
		Account: account,
		Address: address,
	}); err != nil {
		return address, err
	}

	return address, nil
}

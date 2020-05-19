package wire

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/spf13/viper"

	"github.com/dmitriy-vas/node/models"
	"github.com/dmitriy-vas/node/postgres"
)

const (
	BitcoinAbbreviation = "btc"
	BitcoinFee          = 0.000000007
)

func init() {
	if viper.Get("bitcoin") == nil {
		return
	}
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         viper.GetString("bitcoin.host"),
		User:         viper.GetString("bitcoin.user"),
		Pass:         viper.GetString("bitcoin.pass"),
		DisableTLS:   true,
		HTTPPostMode: true,
	}, nil)
	if err != nil {
		log.Panicf("Error loading %s: %v", BitcoinAbbreviation, err)
	}
	ICryptoMap[BitcoinAbbreviation] = Bitcoin{
		client,
		&chaincfg.MainNetParams,
	}
}

type Bitcoin struct {
	*rpcclient.Client
	*chaincfg.Params
}

func (b Bitcoin) CreateAccount(account string) (address string, err error) {
	addresses, err := b.GetAddressByLabel(account)
	if err != nil {
		log.Printf(`Addresses by label "%s" not found: %v`, account, err)
		a, err := b.GetNewAddress(account)
		if err != nil {
			return address, err
		}
		address = a.String()
		if err := postgres.Database.AddBitcoinAccount(&models.Bitcoin{
			Account: account,
			Address: address,
		}); err != nil {
			log.Printf("Error adding new %s account: %v", BitcoinAbbreviation, err)
		}
	} else {
		address = addresses[0].String()
		postgres.Database.AddBitcoinAccount(&models.Bitcoin{
			Account: account,
			Address: address,
		})
	}
	return address, nil
}

func (b Bitcoin) GetBalance(account string) (balance *big.Float, err error) {
	log.Printf("Decoding %s to address", account)
	address, err := btcutil.DecodeAddress(account, b.Params)
	var addresses []btcutil.Address
	if err != nil {
		log.Printf("Decoding is not succeed: %v", err)
		addresses, err = b.GetAddressByLabel(account)
		if err != nil {
			return balance, err
		}
		log.Printf("Successfully get addresses by label: %v", addresses)
	} else {
		log.Printf("Decoding is succeed: %v", address)
		addresses = []btcutil.Address{address}
	}
	log.Printf("Trying to get slice of unspent outputs")
	unspent, err := b.GetUnspent(addresses)
	if err != nil {
		return balance, err
	}
	balance = big.NewFloat(0)
	for _, u := range unspent {
		balance = big.NewFloat(0).Add(balance, big.NewFloat(u.Amount))
	}
	return balance, nil
}

//func (b Bitcoin) GetBalance(account string) (balance *big.Float, err error) {
//	out, err := postgres.Database.SearchBitcoinAccount(account)
//	if err != nil {
//		return balance, err
//	}
//	response, err := b.RawRequest("getreceivedbylabel", []json.RawMessage{
//		[]byte(fmt.Sprintf("\"%s\"", account))},
//	)
//	if err != nil {
//		return balance, err
//	}
//	log.Printf("Received [%s]: %s", account, response)
//	var received float32
//	if err := json.Unmarshal(response, &received); err != nil {
//		return balance, err
//	}
//	return big.NewFloat(float64(received - out.Spent)), nil
//}

func (b Bitcoin) GetAddressByLabel(account string) (addresses []btcutil.Address, err error) {
	rawAccount, _ := json.Marshal(account)
	raw, err := b.RawRequest("getaddressesbylabel", []json.RawMessage{rawAccount})
	if err != nil {
		return addresses, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return addresses, err
	}
	log.Printf("Get %d addresses by label %s", len(result), account)
	for address := range result {
		a, _ := btcutil.DecodeAddress(address, b.Params)
		addresses = append(addresses, a)
	}
	return addresses, nil
}

func (b Bitcoin) GetUnspent(addresses []btcutil.Address) (unspent []btcjson.ListUnspentResult, err error) {
	return b.ListUnspentMinMaxAddresses(
		1,
		math.MaxInt32,
		addresses,
	)
}

func (b Bitcoin) SendTransaction(from string, to string, amount *big.Float) (id string, err error) {
	if _, err := btcutil.DecodeAddress(from, b.Params); err == nil {
		out, err := postgres.Database.SearchBitcoinAddress(from)
		if err != nil {
			return id, err
		}
		from = out.Account
	}

	if _, err := btcutil.DecodeAddress(to, b.Params); err != nil {
		addresses, err := b.GetAddressByLabel(to)
		if err != nil {
			return id, err
		}
		to = addresses[0].String()
	}

	addresses, err := b.GetAddressByLabel(from)
	if err != nil {
		return id, err
	}
	unspentList, err := b.GetUnspent(addresses)
	if err != nil {
		return id, err
	}

	amountFloat, _ := amount.Float64()
	total := 0.0
	totalSize := 10.0
	txInputs := make([]btcjson.TransactionInput, 0)
	for _, unspent := range unspentList {
		if total > amountFloat {
			break
		}
		totalSize += 181
		total += unspent.Amount
		txInputs = append(txInputs, btcjson.TransactionInput{
			Txid: unspent.TxID,
			Vout: unspent.Vout,
		})
	}

	if total < amountFloat {
		return id, fmt.Errorf("insufficient funds")
	}
	returnAmount, err := btcutil.NewAmount(total - amountFloat)
	if err != nil {
		return id, err
	}
	sendAmount, err := btcutil.NewAmount(amountFloat - (totalSize+68)*BitcoinFee)
	if err != nil {
		return id, err
	}

	toAddress, _ := btcutil.DecodeAddress(to, b.Params)
	amounts := map[btcutil.Address]btcutil.Amount{
		addresses[0]: returnAmount,
		toAddress:    sendAmount,
	}

	rawTX, err := b.CreateRawTransaction(txInputs, amounts, nil)
	signedRawTX, _, err := b.SignRawTransaction(rawTX)
	if err != nil {
		return id, err
	}
	hash, err := b.SendRawTransaction(signedRawTX, false)
	if err != nil {
		return id, err
	}
	return hash.String(), nil
}

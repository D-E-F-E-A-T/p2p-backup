package crypto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/dmitriy-vas/p2p/handler/middleware"
)

var Wrappers = map[string]ICrypto{}

type ICrypto interface {
	GetBalance(account string) (balance *big.Float, err error)
	SendTransaction(request SendTransactionRequest) (id string, err error)
	CreateAccount(account string) (address string, err error)
}

type DefaultWrapper struct {
	*http.Client
	Address      string
	Token        string
	Abbreviation string
}

type ErrorResponse struct {
	Content string `json:"error"`
}

func (er ErrorResponse) Error() string {
	return er.Content
}

const (
	WavesFee    = 0.001
	EthereumFee = 0.000105
)

func CostWithFee(currency interface{}, amount *big.Float) *big.Float {
	switch currency.(type) {
	case uint8:
		currency = middleware.Currency(currency.(uint8))
	}
	switch currency.(string) {
	case "eth":
		return big.NewFloat(0).Add(amount, big.NewFloat(EthereumFee))
	case "wvs":
		return big.NewFloat(0).Add(amount, big.NewFloat(WavesFee))
	}
	return amount
}

func NewDefaultWrapper(address, token, abbreviation string) DefaultWrapper {
	return NewDefaultWrapperWithClient(address, token, abbreviation, &http.Client{
		Timeout: time.Second * 10,
	})
}

func NewDefaultWrapperWithClient(address, token, abbreviation string, client *http.Client) DefaultWrapper {
	return DefaultWrapper{
		Client:       client,
		Address:      address,
		Token:        token,
		Abbreviation: abbreviation,
	}
}

type GetBalanceResponse struct {
	Account string     `json:"account"`
	Balance *big.Float `json:"balance"`
}

func (d DefaultWrapper) GetBalance(address string) (balance *big.Float, err error) {
	url := fmt.Sprintf("%s/%s/GetBalance", d.Address, d.Abbreviation)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return balance, err
	}
	request.Header.Set("Authorization", d.Token)
	request.Header.Set("Content-Type", "application/json")

	query := request.URL.Query()
	query.Set("account", address)
	request.URL.RawQuery = query.Encode()

	response, err := d.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		err = json.Unmarshal(raw, &errResp)
		if err == nil {
			err = errResp
		}
	} else {
		var response GetBalanceResponse
		err = json.Unmarshal(raw, &response)
		balance = response.Balance
	}
	return
}

type SendTransactionRequest struct {
	From   string     `json:"from"`
	To     string     `json:"to"`
	Amount *big.Float `json:"amount"`
}

type SendTransactionResponse struct {
	ID string `json:"id"`
}

func (d DefaultWrapper) SendTransaction(transaction SendTransactionRequest) (id string, err error) {
	raw, err := json.Marshal(transaction)
	if err != nil {
		return id, err
	}

	url := fmt.Sprintf("%s/%s/SendTransaction", d.Address, d.Abbreviation)
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(raw))
	if err != nil {
		return id, err
	}
	defer request.Body.Close()

	//q := request.URL.Query()
	//q.Set("from", transaction.From)
	//q.Set("to", transaction.To)
	//q.Set("amount", transaction.Amount.String())
	//request.URL.RawQuery = q.Encode()

	request.Header.Set("Authorization", d.Token)
	request.Header.Set("Content-Type", "application/json")

	response, err := d.Client.Do(request)
	if err != nil {
		return id, err
	}
	defer response.Body.Close()
	raw, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return id, err
	}

	if response.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		err = json.Unmarshal(raw, &errResp)
		if err == nil {
			err = errResp
		}
	} else {
		var response SendTransactionResponse
		err = json.Unmarshal(raw, &response)
		id = response.ID
	}
	return
}

type CreateAccountResponse struct {
	Account string `json:"account"`
	Address string `json:"address"`
}

func (d DefaultWrapper) CreateAccount(account string) (address string, err error) {
	url := fmt.Sprintf("%s/%s/CreateAccount", d.Address, d.Abbreviation)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return address, err
	}
	request.Header.Set("Authorization", d.Token)
	request.Header.Set("Content-Type", "application/json")

	query := request.URL.Query()
	query.Set("account", account)
	request.URL.RawQuery = query.Encode()

	response, err := d.Client.Do(request)
	if err != nil {
		return address, err
	}
	defer response.Body.Close()
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return address, err
	}
	if response.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		err = json.Unmarshal(raw, &errResp)
		if err == nil {
			err = errResp
		}
	} else {
		var response CreateAccountResponse
		err = json.Unmarshal(raw, &response)
		address = response.Address
	}
	return
}

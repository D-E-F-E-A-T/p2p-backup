package wire

import (
	"context"
	"log"
	"math/big"

	"github.com/Dmitriy-Vas/go-cryptoconv"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-pg/pg"
	"github.com/spf13/viper"

	"github.com/dmitriy-vas/node/models"
	"github.com/dmitriy-vas/node/postgres"
)

const (
	EthereumAbbreviation = "eth"
)

func init() {
	if viper.Get("ethereum") == nil {
		return
	}
	rpcClient, err := rpc.Dial(viper.GetString("ethereum.host"))
	if err != nil {
		log.Panicf("Error loading %s: %v", EthereumAbbreviation, err)
	}
	ICryptoMap[EthereumAbbreviation] = Ethereum{
		Client: ethclient.NewClient(rpcClient),
		RPC:    rpcClient,
		KeyStore: keystore.NewKeyStore(viper.GetString("ethereum.keystore"),
			keystore.StandardScryptN,
			keystore.StandardScryptP),
		Context: context.Background(),
	}
}

type Ethereum struct {
	*ethclient.Client
	RPC      *rpc.Client
	KeyStore *keystore.KeyStore
	Context  context.Context
}

func (e Ethereum) CreateAccount(account string) (address string, err error) {
	out, err := postgres.Database.SearchEthereumAccount(account)
	if err != nil && err != pg.ErrNoRows {
		return address, err
	} else if err == nil {
		return out.Address, nil
	}

	acc, err := e.KeyStore.NewAccount("")
	if err != nil {
		return address, err
	}
	if err := postgres.Database.AddEthereumAccount(&models.Ethereum{
		Account: account,
		Address: acc.Address.String(),
	}); err != nil {
		return address, err
	}
	return acc.Address.String(), nil
}

func (e Ethereum) GetBalance(account string) (balance *big.Float, err error) {
	if !common.IsHexAddress(account) {
		out, err := postgres.Database.SearchEthereumAccount(account)
		if err != nil {
			return balance, err
		}
		account = out.Address
	}
	address := common.HexToAddress(account)
	balanceInt, err := e.Client.BalanceAt(e.Context, address, nil)
	if err != nil {
		return balance, err
	}

	return cryptoconv.From(balanceInt.String(), "eth"), nil
}

func (e Ethereum) SendTransaction(from string, to string, amount *big.Float) (id string, err error) {
	if !common.IsHexAddress(from) {
		out, err := postgres.Database.SearchEthereumAccount(from)
		if err != nil {
			return id, err
		}
		from = out.Address
	}

	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)
	fromNonce, err := e.NonceAt(e.Context, fromAddress, nil)
	if err != nil {
		return id, err
	}

	value, _ := amount.Int(nil)
	transaction := types.NewTransaction(
		fromNonce,
		toAddress,
		new(big.Int).Mul(value, big.NewInt(params.Wei)), // amount in wei
		21000,                                           // default transaction gas limit
		big.NewInt(5000000000),                          // 5 gwei
		nil,
	)

	fromAccount, _ := e.KeyStore.Find(accounts.Account{
		Address: fromAddress,
	})
	e.KeyStore.Unlock(fromAccount, "")
	chainID, _ := e.ChainID(e.Context)
	signedTransaction, err := e.KeyStore.SignTx(fromAccount, transaction, chainID)
	if err != nil {
		return id, err
	}

	return signedTransaction.Hash().String(),
		e.Client.SendTransaction(e.Context, signedTransaction)
}

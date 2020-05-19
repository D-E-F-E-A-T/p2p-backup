package Client

import (
	"encoding/base64"
	"log"
	"math/big"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"

	"github.com/dmitriy-vas/p2p/crypto"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
	"github.com/dmitriy-vas/p2p/quotas"
)

type WalletInfo struct {
	Address  string
	Currency string
	Balance  *big.Float
	Volume   float32
	Volumes  map[string]float32
}

type WalletsRequest struct {
	Currency string `form:"currency" binding:"omitempty,fiat"`
}

func Wallets(c *gin.Context) {
	var request WalletsRequest
	if err := c.Bind(&request); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	if request.Currency == "" {
		request.Currency = "usd"
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	data := gin.H{
		"NumOffers": c.GetInt("NumOffers"),
		"NumDeals":  c.GetInt("NumDeals"),
		"IsLogged":  c.GetBool("IsLogged"),
		"Currency":  request.Currency,
	}

	var totalVolume float32
	for _, currency := range middleware.CryptoCurrencies {
		if wrapper, ok := crypto.Wrappers[currency]; ok {
			balance, err := wrapper.GetBalance(user.Login)
			if err != nil {
				log.Printf("Error getting %s balance: %v", currency, err)
				continue
			}
			balanceFloat, _ := balance.Float32()
			cost := quotas.Quotas[currency][request.Currency]

			walletInfo := WalletInfo{
				Currency: currency,
				Balance:  balance,
				Volume:   cost * balanceFloat,
			}
			totalVolume += walletInfo.Volume

			data[currency] = walletInfo
		}
	}
	data["TotalVolume"] = totalVolume

	c.HTML(http.StatusOK, "wallets.html", data)
}

type WalletRequest struct {
	Currency string `form:"currency" binding:"crypto"`
}

func Wallet(c *gin.Context) {
	var request WalletRequest
	if err := c.Bind(&request); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	data := gin.H{
		"NumOffers":           c.GetInt("NumOffers"),
		"NumDeals":            c.GetInt("NumDeals"),
		"IsLogged":            c.GetBool("IsLogged"),
		"NotificationsAmount": c.GetInt("NotificationsAmount"),
		"Notifications":       c.MustGet("Notifications"),
	}

	var wallets []WalletInfo
	for _, currency := range middleware.CryptoCurrencies {
		if wrapper, ok := crypto.Wrappers[currency]; ok {
			walletInfo := WalletInfo{
				Currency: currency,
			}

			if currency == request.Currency {
				address, err := wrapper.CreateAccount(user.Login)
				if err != nil {
					log.Printf("Error creating %s address: %v", currency, err)
				}
				walletInfo.Address = address
			}

			balance, err := wrapper.GetBalance(user.Login)
			if err != nil {
				log.Printf("Error getting %s balance: %v", currency, err)
				continue
			}
			walletInfo.Balance = balance

			if currency == request.Currency {
				codeBytes, _ := qrcode.Encode(walletInfo.Address, qrcode.Highest, 200)
				data["Image"] = base64.StdEncoding.EncodeToString(codeBytes)

				balanceFloat, _ := balance.Float32()
				walletInfo.Volumes = quotas.Quotas[currency]
				walletInfo.Volume = walletInfo.Volumes["usd"] * balanceFloat
				data["Wallet"] = walletInfo
			} else {
				wallets = append(wallets, walletInfo)
			}
		}
	}
	data["Wallets"] = wallets

	c.HTML(http.StatusOK, "wallet.html", data)
}

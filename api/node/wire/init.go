package wire

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type ICrypto interface {
	GetBalance(account string) (balance *big.Float, err error)
	SendTransaction(from string, to string, amount *big.Float) (id string, err error)
	CreateAccount(account string) (address string, err error)
}

var ICryptoMap = make(map[string]ICrypto)
var Mutex = new(sync.Mutex)

func GetICrypto(currency string) ICrypto {
	Mutex.Lock()
	defer Mutex.Unlock()
	return ICryptoMap[currency]
}

func Serve() {
	server := gin.Default()
	server.Use(Middleware)

	file, err := os.OpenFile("output.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0655)
	if err != nil {
		log.Panicf("Error creating or opening log output: %v", err)
	}
	defer file.Close()
	log.SetOutput(file)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("node", validateNode)
	}

	server.Use(gin.LoggerWithWriter(log.Writer()))

	server.GET("/:node/GetBalance", GetBalance)
	server.PUT("/:node/SendTransaction", SendTransaction)
	server.GET("/:node/CreateAccount", CreateAccount)

	server.Run(fmt.Sprintf("%s:%d",
		viper.GetString("server.host"),
		viper.GetInt("server.port"),
	))
}

var validNodes = []string{
	"btc",
	"bch",
	"eth",
	"ltc",
	"xrp",
	"wvs",
}

func validateNode(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	for _, node := range validNodes {
		if v == node {
			return true
		}
	}
	return false
}

type MiddlewareStruct struct {
	Node string `uri:"node" binding:"node"`
}

func Middleware(c *gin.Context) {
	var middleware MiddlewareStruct
	if err := c.BindUri(&middleware); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Set("Node", middleware.Node)
	c.Next()
}

type GetBalanceRequest struct {
	Account string `form:"account" binding:"required"`
}

func GetBalance(c *gin.Context) {
	var request GetBalanceRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	wrapper := GetICrypto(c.GetString("Node"))
	balance, err := wrapper.GetBalance(request.Account)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"account": request.Account,
		"balance": balance,
	})
}

type SendTransactionRequest struct {
	From   string     `form:"from" binding:"required"`
	To     string     `form:"to" binding:"required"`
	Amount *big.Float `form:"amount" binding:"required"`
}

func SendTransaction(c *gin.Context) {
	var request SendTransactionRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	wrapper := GetICrypto(c.GetString("Node"))
	id, err := wrapper.SendTransaction(
		request.From,
		request.To,
		request.Amount,
	)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

type CreateAccountRequest struct {
	Account string `form:"account" binding:"required"`
}

func CreateAccount(c *gin.Context) {
	var request CreateAccountRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	wrapper := GetICrypto(c.GetString("Node"))
	address, err := wrapper.CreateAccount(request.Account)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"account": request.Account,
		"address": address,
	})
}

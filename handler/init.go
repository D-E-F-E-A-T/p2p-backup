package handler

import (
	"encoding/gob"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/olahol/melody.v1"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/Account"
	"github.com/dmitriy-vas/p2p/handler/Client"
	"github.com/dmitriy-vas/p2p/handler/Requests"
	"github.com/dmitriy-vas/p2p/handler/Socket"
	"github.com/dmitriy-vas/p2p/handler/SocketBackup"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/l18n"
	"github.com/dmitriy-vas/p2p/models"
)

func init() {
	gin.SetMode(gin.DebugMode)
	server := gin.New()
	mrouter := melody.New()

	logOutput, err := os.OpenFile("./output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Panicf("Error opening log file: %v", err)
	}
	log.SetOutput(logOutput)

	server.Use(gin.RecoveryWithWriter(log.Writer()))
	server.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: log.Writer(),
	}))
	authKey, _ := hex.DecodeString(viper.GetString("api.keys.auth"))
	encryptKey, _ := hex.DecodeString(viper.GetString("api.keys.encrypt"))
	cookieStore := cookie.NewStore(
		authKey,
		encryptKey,
	)
	gob.Register((*models.User)(nil))
	server.Use(sessions.Sessions(middleware.SessionName, cookieStore))
	server.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    cors.DefaultConfig().AllowMethods,
		AllowHeaders:    cors.DefaultConfig().AllowHeaders,
		MaxAge:          cors.DefaultConfig().MaxAge,
		AllowWebSockets: true,
	}))
	server.StaticFS("/js", http.Dir("./public/templates/js"))
	server.StaticFS("/css", http.Dir("./public/templates/css"))
	server.StaticFS("/img", http.Dir("./public/templates/img"))
	server.StaticFS("/data", http.Dir("./public/templates/data"))
	server.StaticFS("/files", http.Dir(viper.GetString("api.files")))
	server.Use(middleware.TokenChecker)
	server.Use(middleware.Client)
	server.GET("/Socket", func(c *gin.Context) {
		sess := sessions.Default(c)
		userInterface := sess.Get("User")
		user := userInterface.(*models.User)

		if err := mrouter.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{
			"User": user,
		})
			err != nil {
			log.Printf("Error: %v", err)
		}
	})
	mrouter.HandleConnect(Socket.SocketHandler.Connect)
	mrouter.HandleDisconnect(Socket.SocketHandler.Disconnect)
	mrouter.HandleMessage(func(session *melody.Session, bytes []byte) {
		Socket.SocketHandler.Handle(mrouter, session, bytes)
	})

	//server.Any("/Socket/*any", func(context *gin.Context) {
	//	SocketBackup.Server.ServeHTTP(context)
	//})

	account := server.Group("/Account")
	{
		account.GET("/Activate", Account.Activate)
		account.GET("/Key", Account.Key)
		account.POST("/Login", Account.Login)
		account.PUT("/Register", Account.Register)
		account.GET("/Restore", Account.Restore)
		account.POST("/Recovery", Account.Recovery)
		account.GET("/Logout", Account.Logout)
	}
	server.GET("/GetOffers", Requests.GetOffers)
	server.GET("/GetUserOffers", Requests.GetUserOffers)
	server.GET("/GetOffersOld", Requests.GetOffersOld)
	server.GET("/GetUserDeals", Requests.GetUserDeals)
	server.GET("/GetNotifications", Requests.GetNotifications)
	server.GET("/GetComments", Requests.GetComments)
	server.GET("/GetEntries", Requests.GetEntries)
	server.GET("/GetTransactions", Requests.GetTransactions)
	server.GET("/GetQuotas", Requests.GetQuotas)
	server.GET("/GetCountries", Requests.GetCountries)
	server.GET("/GetAuthCode", Requests.GetAuthCode)
	server.GET("/GetChats", Requests.GetChats)
	server.GET("/GetMessages", Requests.GetMessages)
	server.GET("/GetUsersList", Requests.GetUsersList)
	server.GET("/GetReserveList", Requests.GetReserveList)
	server.GET("/IsAuthorized", middleware.IsAuthorized)

	server.PUT("/NewOffer", Requests.NewOffer)
	server.PUT("/NewDeal", Requests.NewDeal)
	server.PUT("/NewComment", Requests.NewComment)

	server.POST("/UploadFile", Requests.UploadFile)
	server.POST("/EditOffer", Requests.EditOffer)
	server.POST("/DeleteOffer", Requests.DeleteOffer)
	//server.POST("/ToggleOffer", Requests.ToggleOffer)
	server.POST("/AcceptDeal", Requests.AcceptDeal)
	server.POST("/CancelDeal", Requests.CancelDeal)
	server.POST("/ApproveDeal", Requests.ApproveDeal)
	server.POST("/FinishDeal", Requests.FinishDeal)
	server.POST("/SendTransaction", Requests.SendTransaction)
	server.POST("/NewArgue", Requests.NewArgue)
	server.POST("/RevokeDeal", Requests.RevokeDeal)

	server.HTMLRender = middleware.LoadRenderer(viper.GetString("api.templates"), template.FuncMap{
		"l18n":        l18n.T,
		"formatTime":  middleware.FormatTime,
		"formatDate":  middleware.FormatDate,
		"formatFloat": middleware.FormatFloat,
		"isCrypto":    middleware.IsCryptoUint,
		"isFiat":      middleware.IsFiatUint,
		"currency":    middleware.Currency,
		"provider":    middleware.Provider,
		"isActiveDay": middleware.IsActiveDay,
		"country":     middleware.Country,
		"sign":        middleware.Sign,
		"isOnline":    middleware.IsOnline,
		"isGroup":     middleware.IsGroup,
	})
	server.GET("/", Client.Index)
	server.GET("/deal", Client.Deal)
	server.GET("/deals", Client.Deals)
	server.GET("/error", Client.Error)
	server.GET("/new_offer", Client.NewOffer)
	server.GET("/offer", Client.Offer)
	server.GET("/edit_offer", Client.EditOffer)
	server.GET("/offers", Client.Offers)
	server.GET("/verified", Client.Verified)
	server.GET("/cabinet", Client.Cabinet)
	server.GET("/login", Client.Login)
	server.GET("/recovery_login", Client.RecoveryLogin)
	server.GET("/wallets", Client.Wallets)
	server.GET("/wallet", Client.Wallet)
	server.GET("/user", Client.User)
	server.GET("/chat", Client.Chat)
	server.GET("/notifications", Client.Notifications)
	server.GET("/test", Test)

	go SocketBackup.Server.Serve()
	go server.Run(viper.GetString("api.address"))
}

func Test(c *gin.Context) {
	offer, _ := postgres.Postgres.GetOffer(7)

	c.HTML(http.StatusOK, "test.html", gin.H{
		"Offer": offer,
	})
}

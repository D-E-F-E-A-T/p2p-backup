package Client

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/crypto"
	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
	"github.com/dmitriy-vas/p2p/quotas"
)

type DealRequest struct {
	ID uint64 `form:"id" binding:"required"`
}

func Deal(c *gin.Context) {
	var request DealRequest
	if err := c.Bind(&request); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	data, err := postgres.Postgres.GetDeal(request.ID)
	if err != nil {
		if err == pg.ErrNoRows {
			c.Redirect(http.StatusTemporaryRedirect, "/")
		} else {
			c.HTML(http.StatusOK, "error.html", gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	if user.ID != data.UserID && user.ID != data.Offer.UserID {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	haveCurrency := middleware.Currency(data.Offer.HaveCurrency)
	wantCurrency := middleware.Currency(data.Offer.WantCurrency)

	out := gin.H{
		"NumOffers":           c.GetInt("NumOffers"),
		"NumDeals":            c.GetInt("NumDeals"),
		"IsLogged":            c.GetBool("IsLogged"),
		"NotificationsAmount": c.GetInt("NotificationsAmount"),
		"Notifications":       c.MustGet("Notifications"),
		"Data":                data,
		"IsOwner":             user.ID == data.Offer.UserID,
		"Quotas":              quotas.Quotas,
	}

	out["CryptoHolder"] = (out["IsOwner"].(bool) && middleware.IsCryptoUint(data.Offer.HaveCurrency)) ||
		(!out["IsOwner"].(bool) && middleware.IsCryptoUint(data.Offer.WantCurrency))

	if out["CryptoHolder"].(bool) && data.Deal.Accepted {
		settings, err := postgres.Postgres.GetUserSettings(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if settings.Settings.AuthenticatorKey != "" {
			out["TwoFA"] = true
		}

		var wrapper crypto.ICrypto
		if out["IsOwner"].(bool) {
			wrapper = crypto.Wrappers[haveCurrency]
		} else {
			wrapper = crypto.Wrappers[wantCurrency]
		}
		if (out["IsOwner"].(bool) && !data.Deal.FromApproved) ||
			(!out["IsOwner"].(bool) && !data.Deal.ToApproved) {
			if wrapper == nil {
				c.Status(http.StatusForbidden)
				return
			}
			address, err := wrapper.CreateAccount(user.Login)
			if err != nil {
				c.JSON(http.StatusConflict, gin.H{
					"error": err.Error(),
				})
				return
			}
			out["Address"] = address
		}
	}

	var userInfo *postgres.User
	if out["IsOwner"].(bool) {
		userInfo, err = postgres.Postgres.GetUserInfo(data.UserID)
		if err != nil {
			c.HTML(http.StatusOK, "error.html", gin.H{
				"error": err.Error(),
			})
			return
		}
		userInfo.User.ID = data.UserID
	} else {
		userInfo, err = postgres.Postgres.GetUserInfo(data.Offer.UserID)
		if err != nil {
			c.HTML(http.StatusOK, "error.html", gin.H{
				"error": err.Error(),
			})
			return
		}
		userInfo.User.ID = data.Offer.UserID
	}
	if data.Offer.WithDynamic {
		quota := quotas.GetCostWithProfit(wantCurrency, haveCurrency, 0)
		out["Cost"] = data.Deal.Amount * quota
	} else {
		out["Cost"] = data.Offer.Cost * data.Deal.Amount
	}
	out["User"] = userInfo

	c.HTML(http.StatusOK, "deal.html", out)
}

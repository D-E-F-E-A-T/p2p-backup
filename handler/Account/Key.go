package Account

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

type KeyRequest struct {
	Token string `json:"token" form:"token" binding:"required"`
}

func Key(c *gin.Context) {
	var request KeyRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")

	if user, ok := userInterface.(*models.User); ok {
		settings, err := postgres.Postgres.SearchUserSettings(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if checkTelegram(request.Token, user.Telegram) ||
			CheckAuthenticator(request.Token, settings.AuthenticatorKey) {
			SaveToken(c, user.ID)

			if err := sess.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
		} else {
			c.Status(http.StatusForbidden)
		}
	} else {
		c.Status(http.StatusUnauthorized)
	}
}

func checkTelegram(token string, telegram int64) bool {
	if telegram == 0 {
		return false
	}
	if token == "" {
		return false
	}
	// TODO fix telegram authentication
	return true
}

func CheckAuthenticator(token string, secretKey string) bool {
	if secretKey == "" {
		return false
	}
	if token == "" {
		return false
	}
	return token == getTOTP(secretKey)
}

func getTOTP(input string) string {
	key, _ := base32.StdEncoding.DecodeString(input)
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs,
		uint64(time.Now().Unix()/30),
	)

	hash := hmac.New(sha1.New, key)
	hash.Write(bs)
	h := hash.Sum(nil)

	o := h[19] & 15
	r := bytes.NewReader(h[o : o+4])
	var header uint32
	binary.Read(r, binary.BigEndian, &header)

	h12 := (int(header) & 0x7fffffff) % 1000000
	return strconv.Itoa(h12)
}

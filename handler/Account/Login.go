package Account

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

type LoginRequest struct {
	Credentials string `json:"credentials" form:"credentials" binding:"min=5,max=32|email"`
	Password    string `json:"password" form:"password" binding:"min=8,max=48"`
	Key         string `json:"key" form:"key" binding:"omitempty,min=8,max=48"`
}

func Login(c *gin.Context) {
	var request LoginRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := postgres.Postgres.SearchUser(request.Credentials)
	if err != nil {
		if err == pg.ErrNoRows {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	entry := &models.Entry{
		ID:          user.ID,
		Address:     c.ClientIP(),
		Application: c.GetHeader("User-Agent"),
		Timestamp:   time.Now(),
	}
	defer func() {
		postgres.Postgres.AddNewEntry(entry)
	}()

	hash := md5.New()
	hash.Write([]byte(request.Password))
	encodedPassword := hex.EncodeToString(hash.Sum(nil))

	if strings.Compare(encodedPassword, user.Password) != 0 {
		c.Status(http.StatusForbidden)
		return
	}

	statistic, err := postgres.Postgres.SearchUserStatistic(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !statistic.EmailVerified {
		c.Status(http.StatusUnauthorized)
		return
	} else {
		entry.Status = true
	}

	settings, err := postgres.Postgres.SearchUserSettings(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(settings.AvailableAddresses) != 0 {
		if !IsExist(settings.AvailableAddresses, c.ClientIP()) {
			c.Status(http.StatusForbidden)
			return
		}
	}

	needAuth := false
	if settings.TelegramAuthentication || settings.AuthenticatorKey != "" {
		needAuth = true
	}

	var sess sessions.Session
	if !needAuth {
		sess = SaveToken(c, user.ID)
	}
	if sess == nil {
		sess = sessions.Default(c)
	}

	sess.Options(sessions.Options{
		MaxAge: 86400 * 7,
		Path:   "/",
	})
	sess.Set("User", *user)

	if err := sess.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !needAuth {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusAccepted)
	}
}

func SaveToken(c *gin.Context, id uint64) sessions.Session {
	sess := sessions.Default(c)

	bs := make([]byte, 64)
	rand.Read(bs)
	token := hex.EncodeToString(bs)
	sess.Set("Token", token)

	if err := postgres.Postgres.AddNewUserSession(&models.Session{
		ID:      id,
		Expires: time.Now().AddDate(0, 0, 7),
		Token:   token,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	return sess
}

func IsExist(slice []string, value string) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		}
	}
	return false
}

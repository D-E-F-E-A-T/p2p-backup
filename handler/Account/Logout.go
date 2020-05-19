package Account

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

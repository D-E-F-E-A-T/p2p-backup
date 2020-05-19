package Client

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "recovery_login.html", nil)
}

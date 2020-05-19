package Client

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Verified(c *gin.Context) {
	c.HTML(http.StatusOK, "verified.html", nil)
}

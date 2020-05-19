package Requests

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/quotas"
)

func GetQuotas(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, quotas.Quotas)
}

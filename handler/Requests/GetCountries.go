package Requests

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/database/postgres"
)

type GetCountriesRequest struct {
	Language string `form:"language" binding:"language"`
}

func GetCountries(c *gin.Context) {
	var request GetCountriesRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	countries, err := postgres.Postgres.GetCountries(request.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"countries": countries,
	})
}

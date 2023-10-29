package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func (t HealthController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "Ok")
}

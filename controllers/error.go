package controllers

import (
	"log"
	"platform_api/models"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, statusCode int, msg string, err error) {
	log.Print(msg, err.Error())
	c.JSON(statusCode, models.HTTPError{
		Code:    statusCode,
		Message: msg,
		Error:   err.Error(),
	})
}

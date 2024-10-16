package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Error_ struct {
	Message string `json:"message"`
}

func isNotFoundInDB(err string) bool {
	return strings.Contains(err, "sql: no rows in result set")
}

func NewErrorResponse(c *gin.Context, statusCode int, message string, logMessage string) {
	if statusCode == http.StatusInternalServerError && len(message) == 0 {
		message = "something is wrong"
	}
	if len(logMessage) == 0 {
		logMessage = message
	}
	logrus.Error(logMessage)
	if isNotFoundInDB(logMessage) {
		c.AbortWithStatusJSON(404, Error_{Message: "Not found"})
	} else {
		c.AbortWithStatusJSON(statusCode, Error_{Message: message})
	}
}

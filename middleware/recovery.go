package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LogFailedRequests Logs failed requests where ever a panic is created. and logs the incoming request .
func LogFailedRequests(c *gin.Context, recovered interface{}) {

	if err, ok := recovered.(string); ok {
		FailedRequestLogger(c)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}

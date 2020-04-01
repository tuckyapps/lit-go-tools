package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tuckyapps/lit-go-tools/api"
)

// ExtractAppVersion is the handler used to retrieve the application version
// from HTTP header
func ExtractAppVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		if headerValue := c.Request.Header.Get(api.HTTPHeaderAppVersion); headerValue != "" {
			c.Set(api.HandlerKeyAppVersion, headerValue)
		}
		c.Next()
	}
}

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tuckyapps/lit-go-tools/api"
)

// ExtractPlatform is the handler used to retrieve the device platoform
// from HTTP header
func ExtractPlatform() gin.HandlerFunc {
	return func(c *gin.Context) {
		if headerValue := c.Request.Header.Get(api.HTTPHeaderPlatform); headerValue != "" {
			c.Set(api.HandlerKeyPlatform, headerValue)
		}
		c.Next()
	}
}

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tuckyapps/lit-go-tools/api"
)

// ExtractLanguage is the handler used to retrieve language header
func ExtractLanguage() gin.HandlerFunc {
	return func(c *gin.Context) {
		if headerValue := c.Request.Header.Get(api.HTTPHeaderAcceptLanguage); headerValue != "" {
			c.Set(api.HandlerKeyLanguage, headerValue)
		}
		c.Next()
	}
}

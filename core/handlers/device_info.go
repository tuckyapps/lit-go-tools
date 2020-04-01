package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tuckyapps/lit-go-tools/api"
)

// ExtractDeviceInfo is the handler used to retrieve the device info
func ExtractDeviceInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		if headerValue := c.Request.Header.Get(api.HTTPHeaderDeviceInfo); headerValue != "" {
			c.Set(api.HandlerKeyDeviceInfo, headerValue)
		}

		c.Next()
	}
}

// ExtractDeviceID is the handler used to retrieve the device id
func ExtractDeviceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if headerValue := c.Request.Header.Get(api.HTTPHeaderDeviceID); headerValue != "" {
			c.Set(api.HandlerKeyDeviceID, headerValue)
		}
		c.Next()
	}
}

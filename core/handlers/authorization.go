package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuckyapps/lit-go-tools/api"
)

// AuthorizeAccessToResource is in charge of checking if the requested resource
// can be accesed by the sent token
//
// For example:
// Let's say you want to get the lists for a Venue ID 12345 and you have
// a valid token generated for your (user ID 9988); so you send a request to
// something like GET->/venue/12345/lists.
//
// If you are not the owner of the venue, system should return a
//
func AuthorizeAccessToResource(resource api.ResourceType) gin.HandlerFunc {
	return func(c *gin.Context) {

		switch resource {

		case api.ResourceUser:
			resourceID := c.Param(api.ParamUserID)
			requestorID, _ := c.Get(api.HandlerKeyTokenID)

			if validateAccessToUser(resourceID, requestorID.(string)) {
				c.Next()
			} else {
				// 403 -> should we return 404?
				c.Writer.WriteHeader(http.StatusForbidden)
				c.Abort()
			}

		default:
			// 403 -> should we return 404?
			c.Writer.WriteHeader(http.StatusForbidden)
		}

	}
}

// Validates if requestorID is allowed to access userID.
// At this implementation, the user accesing the resource must have the same ID.
func validateAccessToUser(userID string, requestorID string) bool {
	return userID == requestorID
}

package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tuckyapps/lit-go-tools/api"
	"github.com/tuckyapps/lit-go-tools/core/authentication"
)

// RequireAccessToken is the function used to handle all the request
// that requires a token as part of the request.
func RequireAccessToken(realm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract token from request
		token, err := ExtractAuthorizationToken(c.Request, api.AuthorizationMethodBearer)
		if err != nil {

			if err == api.ErrAuthHeaderNotFound {
				// if api.HTTPHeaderAuthorization header was not found, send response without details
				resp := api.BuildEmptyUnauthorizedResponse(realm, api.AuthorizationMethodBearer)
				resp.Send(c.Writer)
			} else {
				msg := api.GetMessage(api.ErrorInvalidToken, api.GetLanguage(c))
				resp := api.BuildUnauthorizedResponse(api.ErrorInvalidToken, msg, realm, api.AuthorizationMethodBearer)
				resp.Send(c.Writer)
			}

			// prevent executing other handlers
			c.Abort()

		} else {

			// validate extracted token
			id, errToken := authentication.ValidateToken(token)
			if errToken != nil {

				switch errToken {

				case api.ErrInvalidToken:

					// why this? avoid token validity check if the operation is Logout; I don't like it, but
					// this will prevent 401 during logout if the token ha expired.

					if c.Request.URL.Path != "/v1/auth/revoke" {
						api.BuildUnauthorizedResponse(api.ErrorInvalidToken, errToken.Error(), realm, api.AuthorizationMethodBearer).Send(c.Writer)
					} else {

						c.Set(api.HandlerKeyTokenID, id)
						c.Next()
					}

				case api.ErrAuthService:
					api.BuildInternalErrorResponse().Send(c.Writer)

				default:
					api.BuildUnauthorizedResponse(api.ErrorInvalidToken, errToken.Error(), realm, api.AuthorizationMethodBearer).Send(c.Writer)

				}

				// prevent executing other handlers
				c.Abort()

			} else {
				// Before executing the next stage in the pipeline, add the ID extracted from the token
				// as a Header in the request, so it's available to the rest of the pipeline
				c.Set(api.HandlerKeyTokenID, id)
				c.Next()
			}
		}
	}
}

// RequireBasicAuthorization is the funcion used to require Basic authentication
// header in a request
func RequireBasicAuthorization(realm string, mandatory bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		// extract token from request
		token, err := ExtractAuthorizationToken(c.Request, api.AuthorizationMethodBasic)
		if err != nil {

			if mandatory {
				if err == api.ErrAuthHeaderNotFound {
					// if api.HTTPHeaderAuthorization header was not found, send response without details
					resp := api.BuildEmptyUnauthorizedResponse(realm, api.AuthorizationMethodBasic)
					resp.Send(c.Writer)
				} else {
					msg := api.GetMessage(api.ErrorInvalidToken, api.GetLanguage(c))
					resp := api.BuildUnauthorizedResponse(api.ErrorInvalidToken, msg, realm, api.AuthorizationMethodBasic)
					resp.Send(c.Writer)
				}

				// prevent executing other handlers
				c.Abort()
			} else {

				// if presence of Basic auth is not mandatory, the handler will not fail and the
				// rest of the pipeline will be executed
				c.Next()

			}

		} else {
			// OK -> continue

			basicData, errDecode := base64.StdEncoding.DecodeString(token)
			if errDecode != nil {
				apiResponse := new(api.Response)
				apiResponse.Status = http.StatusBadRequest
				apiResponse.ErrCode = api.ErrorUnkownUser
				apiResponse.Send(c.Writer)
				return
			}

			// extract data and save it in the pipeline
			clientID, clientSecret, errParse := parseBasicAuth(string(basicData))
			if errParse == nil {
				c.Set(api.HandlerKeyClientID, clientID)
				c.Set(api.HandlerKeyClientSecret, clientSecret)
				c.Next()
			} else {
				apiResponse := new(api.Response)
				apiResponse.Status = http.StatusBadRequest
				apiResponse.ErrCode = api.ErrorInvalidToken
				apiResponse.Send(c.Writer)
			}
		}
	}
}

// ExtractAuthorizationToken parses the token from the 'Authorization' header
func ExtractAuthorizationToken(req *http.Request, method string) (token string, err error) {
	if headerValue := req.Header.Get(api.HTTPHeaderAuthorization); headerValue != "" {
		// supported methods are Bearer and Basic
		switch method {
		case api.AuthorizationMethodBearer:
			token, err = stripBearerPrefixFromTokenString(headerValue)
		case api.AuthorizationMethodBasic:
			token, err = stripBasicPrefixFromTokenString(headerValue)
		default:
			return "", fmt.Errorf("Authorizaion method '%s' is not supported by the platform", method)
		}
	} else {
		// header not found
		err = api.ErrAuthHeaderNotFound
	}
	return
}

// strips 'Bearer ' prefix from bearer token string
func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

// strips 'Basic ' prefix from basic auth string
func stripBasicPrefixFromTokenString(tok string) (string, error) {
	// should be a bearer token
	if len(tok) > 5 && strings.ToUpper(tok[0:6]) == "BASIC " {
		return tok[6:], nil
	}
	return tok, nil
}

// parses the basic auth token extacting the client id and secret
func parseBasicAuth(basicToken string) (clientID, clientSecret string, err error) {
	split := strings.Split(basicToken, ":")
	if len(split) == 2 {
		clientID = split[0]
		clientSecret = split[1]
	} else {
		err = api.ErrUserNotFound
	}

	return
}

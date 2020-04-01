package authentication

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/parnurzeal/gorequest"
	"github.com/pquerna/ffjson/ffjson"

	"github.com/tuckyapps/lit-go-tools/api"
	"github.com/tuckyapps/lit-go-tools/logger"
)

type authServiceError struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

// AuthService is the standard auth service implementation
type AuthService struct {
	authClient     string
	authServiceURL string
}

var (
	authImpl AuthService
)

// Init initializes the auth service.
func Init(authClient string, authServiceURL string) {
	authImpl.authClient = authClient
	authImpl.authServiceURL = authServiceURL
	return
}

// ValidateToken sends a request to the auth-service and validates that the
// token is geniune: not in black list and it hasn't expired.
func ValidateToken(token string) (id string, err error) {

	log := logger.GetLogger()
	requestURL := fmt.Sprintf("%s/v1/token/validate", authImpl.authServiceURL)
	requestBody := fmt.Sprintf(`{"token":"%s"}`, token)
	decoder := ffjson.NewDecoder()

	// build request
	request := gorequest.New().Post(requestURL).Send(requestBody).Timeout(api.HTTPTimeout)
	request.Header.Set("Auth-Client", authImpl.authClient)

	if response, _, errPost := request.End(); errPost == nil {

		if response.StatusCode != http.StatusOK {

			// error response should match authServiceError
			var errBody authServiceError
			if errJSON := decoder.DecodeReader(response.Body, &errBody); errJSON == nil {
				if errBody.Error != api.ErrorInvalidToken {
					// other error from auth service
					err = errors.New(errBody.Description)
				} else {
					// invalid token
					err = api.ErrInvalidToken
				}
			} else {
				// could not understand response from auth service
				err = api.ErrBadRequest
			}
		} else {

			// retrieve ID since it was already already parsed by auth-service
			var claims api.Claim
			if errJSON := decoder.DecodeReader(response.Body, &claims); errJSON == nil {
				id = claims.Claims["id"].(string)
			}
		}
	} else {
		log.Errorf("Error calling auth-service: %v", errPost)
		err = api.ErrAuthService
	}

	return
}

// GenerateToken connects to the authorization service and retrieves a new token
func GenerateToken(id string) (token *api.Token, err error) {

	log := logger.GetLogger()
	requestURL := fmt.Sprintf("%s/v1/token/generate", authImpl.authServiceURL)
	requestBody := fmt.Sprintf(`{"claims":{"grant":"access_token","id":"%s"}}`, id)

	// build request
	request := gorequest.New().Post(requestURL).Send(requestBody).Timeout(api.HTTPTimeout)
	request.Header.Set("Auth-Client", authImpl.authClient)

	if response, _, errRequest := request.End(); errRequest == nil {

		decoder := ffjson.NewDecoder()
		if response.StatusCode != http.StatusOK {
			// get error details
			errBody := new(authServiceError)

			if errJSON := decoder.DecodeReader(response.Body, errBody); errJSON == nil {
				err = errors.New(errBody.Description)
			}
		} else {
			if errJSON := decoder.DecodeReader(response.Body, &token); errJSON != nil {
				err = errors.New("Token response from auth service is not in the expected format")
			}
		}
	} else {
		log.Errorf("Error calling auth-service: %v", errRequest)
		err = api.ErrAuthService
	}

	return
}

// Revoke marks the current token as invalid
func Revoke(token string) error {

	log := logger.GetLogger()
	requestURL := fmt.Sprintf("%s/v1/token/destroy", authImpl.authServiceURL)
	requestBody := fmt.Sprintf(`{"token":"%s"}`, token)

	// build request
	request := gorequest.New().Post(requestURL).Send(requestBody).Timeout(api.HTTPTimeout)
	request.Header.Set("Auth-Client", authImpl.authClient)

	if response, _, errRequest := request.End(); errRequest == nil {

		if response.StatusCode != http.StatusOK {
			// error response should match authServiceError
			errBody := authServiceError{}
			decoder := ffjson.NewDecoder()

			if errJSON := decoder.DecodeReader(response.Body, &errBody); errJSON == nil {
				if errBody.Error != api.ErrorInvalidToken {
					return errors.New(errBody.Description)
				}
			}

			return api.ErrBadRequest
		}
	} else {
		log.Errorf("Error calling auth-service: %v", errRequest)
	}

	return nil
}

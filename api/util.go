package api

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/pquerna/ffjson/ffjson"
)

// Response represents a HTTP Response
type Response struct {
	ContentType    string
	Status         int
	Header         map[string]string
	ErrCode        string
	ErrDescription string
	Fields         *[]string
	Payload        []byte
	Language       string
}

// Send writes the HTTP response in a http.ResponseWriter
func (response *Response) Send(w http.ResponseWriter) {

	// default content type is json
	if response.ContentType == "" {
		response.ContentType = "application/json"
	}
	w.Header().Set("Content-Type", response.ContentType)

	// write default headers and others
	if len(response.Header) > 0 {
		for k, v := range response.Header {
			w.Header().Set(k, v)
		}
	}

	// status code
	w.WriteHeader(response.Status)

	// create a default payload if it wasn't proportionated
	if response.Status >= http.StatusBadRequest {
		if response.Payload == nil && response.ErrCode != "" {

			// check if an error description was provided or look for a default message
			var errDesc string
			if response.ErrDescription != "" {
				errDesc = response.ErrDescription
			} else {
				errDesc = GetMessage(response.ErrCode, response.Language)
			}

			// create error data, convert to json and store it in the payload
			errData := ErrorData{
				Error:       response.ErrCode,
				Description: errDesc,
				Fields:      response.Fields,
			}

			response.Payload = errData.Marshall()
		}
	}

	w.Write(response.Payload)
}

// String returns a string representation of the error, suitable for
// including it in a HTTP header
func (response *Response) String() string {
	var buffer bytes.Buffer
	hasSomething := false

	// error
	if response.ErrCode != "" {
		buffer.WriteString(fmt.Sprintf("error=\"%s\"", response.ErrCode))
		hasSomething = true
	}

	// description
	if response.ErrDescription != "" {
		if hasSomething {
			buffer.WriteString(", ")
		}
		buffer.WriteString(fmt.Sprintf("error_description=\"%s\"", response.ErrDescription))
	}

	return buffer.String()
}

// ErrorData is used to represent API errors
type ErrorData struct {
	Error       string    `json:"error,omitempty"`
	Description string    `json:"description,omitempty"`
	Fields      *[]string `json:"fields,omitempty"`
}

// Message is used to send no error responses to the client
type Message struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Marshall encodes ErrorData content to a json byte array
func (errData *ErrorData) Marshall() []byte {
	if b, err := ffjson.Marshal(errData); err == nil {
		return b
	}
	return []byte("")
}

// BuildUnauthorizedResponse creates a tipical response for a 401 status code
func BuildUnauthorizedResponse(code, description, realm, method string) *Response {
	resp := new(Response)
	resp.Status = http.StatusUnauthorized
	resp.ErrCode = code
	resp.ErrDescription = description

	// headers
	headers := make(map[string]string)
	headers["WWW-Authenticate"] = fmt.Sprintf("%s realm=\"%s\", %s", method, realm, resp.String())
	resp.Header = headers

	return resp
}

// BuildEmptyUnauthorizedResponse creates a 401 response, without any error details.
// It should be used when an authentication method was not supported.
func BuildEmptyUnauthorizedResponse(realm, method string) *Response {
	resp := new(Response)
	resp.Status = http.StatusUnauthorized

	// headers
	headers := make(map[string]string)
	headers["WWW-Authenticate"] = fmt.Sprintf("%s realm=\"%s\"", method, realm)
	resp.Header = headers

	return resp
}

// BuildInternalErrorResponse creates a tipical response for a 500 status code
func BuildInternalErrorResponse() *Response {
	resp := new(Response)
	resp.Status = http.StatusInternalServerError
	resp.ErrCode = ErrorInternalError

	return resp
}

// BuildBadRequestResponse creates a tipical response for a 400 status code
func BuildBadRequestResponse() *Response {
	resp := new(Response)
	resp.Status = http.StatusBadRequest
	resp.ErrCode = ErrorInvalidRequest

	return resp
}

// BuildForbiddenResponse creates a 403 response, without any error details
func BuildForbiddenResponse() *Response {
	resp := new(Response)
	resp.Status = http.StatusForbidden
	resp.ErrCode = ErrorForbiddenRequest

	return resp
}

// GetLanguage returns the specified of default language
func GetLanguage(c *gin.Context) string {
	var lang string
	if l, exists := c.Get(HandlerKeyLanguage); exists {
		lang = strings.ToLower(l.(string))
	} else {
		lang = DefaultLanguage
	}

	// if not a supported language, use the default
	if !strings.Contains(SupportedLanguages, lang) {
		lang = DefaultLanguage
	}

	return lang
}

// CreateNewUUID creates a new UUID without dashes
func CreateNewUUID() string {
	uuid, _ := uuid.NewV4()
	return strings.Replace(uuid.String(), "-", "", -1)
}

// OnlyNumbers extract only numbers from str
func OnlyNumbers(str string) string {
	re := regexp.MustCompile("[0-9]+")
	s := re.FindAllString(str, -1)
	return strings.Join(s, "")
}

// ArrayFind searches for an element `val` within array `a`. If the value exists,
// it will be returned by the function.
func ArrayFind(a []interface{}, val interface{}) interface{} {
	for _, elem := range a {
		if elem == val {
			return elem
		}
	}
	return nil
}

// ValidateUserAge validate user age
func ValidateUserAge(birthDate *time.Time) (valid bool) {
	birthDate18 := birthDate.AddDate(18, 0, 0)
	if birthDate18.Before(time.Now()) {
		return true
	}
	return false
}

// ParseToInt64 Parses string to int64.
func ParseToInt64(ID string) (int64, error) {
	return strconv.ParseInt(ID, 10, 64)
}

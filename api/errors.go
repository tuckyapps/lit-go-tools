package api

import "errors"

// API error codes
const (
	ErrorInvalidToken            = "invalid_token"
	ErrorAuthService             = "err_auth_service"
	ErrorInternalError           = "internal_error"
	ErrorMissingParameters       = "missing_parameters"
	ErrorUserExists              = "user_exists"
	ErrorCountryNotSupported     = "invalid_country"
	ErrorInvalidRequest          = "invalid_request"
	ErrorUnkownUser              = "unknown_user"
	ErrorUnkownVenue             = "unknown_venue"
	ErrorAccountBlocked          = "account_blocked"
	ErrorInvalidParameters       = "invalid_parameters"
	ErrorFileTooBig              = "file_too_big"
	ErrorMissingInvite           = "missing_invite"
	ErrorExpectedBoolNotRecieved = "expected_bool_but_recieved_other"
	ErrorImageFormatNotSupported = "image_format_not_supported"
	ErrorImageSizeNotSupported   = "image_size_not_supported"
	ErrorForbiddenRequest        = "forbidden_request"
)

// Internal error types
var (
	ErrInternalError           = errors.New("Internal error")
	ErrAuthService             = errors.New("Error sending request to the Authorization Service")
	ErrAuthHeaderNotFound      = errors.New("Missing Authentication header")
	ErrBadRequest              = errors.New("Bad request")
	ErrUserNotFound            = errors.New("Invalid user")
	ErrInvalidToken            = errors.New("Invalid token")
	ErrAccountBlocked          = errors.New("Account blocked")
	ErrMissingRequiredFields   = errors.New("Missing required fields")
	ErrAwsNotEnabled           = errors.New("AWS not enabled")
	ErrAwsS3Timeout            = errors.New("Upload s3 timeout")
	ErrAwsUploadS3             = errors.New("Error upload s3 object")
	ErrUnregisterdUser         = errors.New("UNREGISTERED")
	ErrImageFormatNotSupported = errors.New("Image format not supported")
	ErrImageSizeNotSupported   = errors.New("Image size not supported")
)

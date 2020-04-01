package api

import "time"

// General constants
const (
	RealmUsers        = "users"
	RealmVenues       = "venues"
	HTTPTimeout       = time.Second * 10
	TimestampFormat   = "2006-01-02T15:04:05.000Z"
	AppInitVersioning = "versioning"
	AppInitURLs       = "urls"
	PlatformAndroid   = "android"
	PlatformIOS       = "ios"
	AppIDAndroid      = "lit.android"
	AppIDIOS          = "lit.ios"
	MaxFileUploadSize = int64(10240000) // 10MB
)

// Supported authorization methods
const (
	AuthorizationMethodBearer = "Bearer"
	AuthorizationMethodBasic  = "Basic"
)

// ResourceType is used to represent internal system resources like User
type ResourceType string

// API resources defined in the API, required to handle the authorization rules
const (
	ResourceUser ResourceType = "USER"
)

// HTTP headers used in the API
const (
	HTTPHeaderAuthorization  = "Authorization"
	HTTPHeaderAppName        = "AppName"
	HTTPHeaderPlatform       = "Platform"
	HTTPHeaderAppVersion     = "AppVersion"
	HTTPHeaderDeviceInfo     = "Device-Info"
	HTTPHeaderDeviceToken    = "Device-Token"
	HTTPHeaderAcceptLanguage = "Accept-Language"
	HTTPHeaderDeviceID       = "Device-ID"
)

// Key names used to store info in the HTTP handlers
const (
	HandlerKeyTokenID      = "Token-ID"
	HandlerKeyClientID     = "Client-ID"
	HandlerKeyClientSecret = "Client-Secret"
	HandlerKeyLanguage     = "Lang"
	HandlerKeyDeviceInfo   = "Device-Info"
	HandlerKeyPlatform     = "Platform"
	HandlerKeyAppVersion   = "App-Version"
	HandlerKeyDeviceID     = "Device-ID"
)

// Parameters defined in the API urls
const (
	ParamUserID   = "userID"
	ParamLanguage = "lang"
)

// Query parameters recognized by the API
const (
	QueryParameterPretty           = "pretty"
	QueryParameterCountry          = "country"
	QueryParameterLanguage         = "lang"
	QueryParameterUserData         = "user_data"
	QueryParameterFilter           = "filter"
	QueryParameterTimestamp        = "timestamp"
	QueryParameterSortOrder        = "sort_order"
	QueryParameterExternalUsername = "username"
)

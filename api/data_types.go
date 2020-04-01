package api

// Token is used to get responses from the auth service and forward it to the client
type Token struct {
	ID           string `json:"id,omitempty"`
	Grant        string `json:"grant,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Token        string `json:"token,omitempty"`
	Refresh      string `json:"refresh_token,omitempty"`
	Type         string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	Firebase     struct {
		Token string `json:"token,omitempty"`
	} `json:"firebase,omitempty"`
}

// Claim represents a request/response
type Claim struct {
	Claims map[string]interface{} `json:"claims"`
}

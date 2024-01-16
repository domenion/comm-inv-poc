package models

import "fmt"

type GetTokenRequest struct {
}

type GetTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (t *GetTokenResponse) CreateToken() string {
	return fmt.Sprintf("%s %s", t.TokenType, t.AccessToken)
}

package oauth2

import (
	"fmt"
	"net/http"
	"time"

	"dahbura.me/api/config"
)

type GrantSource interface {
	GrantResponse() (GrantResponse, error)
}

type GrantResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresAt    time.Time
}

func (resp *GrantResponse) HasExpired() bool {
	now := time.Now()
	expiresAtWithLeeway := resp.ExpiresAt.Add(-config.DefaultTokenLeeway)

	return now.After(expiresAtWithLeeway)
}

func ExpiresAt(expiresIn int) time.Time {
	now := time.Now()
	expiresAt := now.Add(time.Duration(expiresIn) * time.Second)

	return expiresAt
}

func SetAuthHeader(req *http.Request, gs GrantSource) error {
	gr, err := gs.GrantResponse()
	if err != nil {
		return err
	}

	val := fmt.Sprintf("Bearer %s", gr.AccessToken)
	req.Header.Set("Authorization", val)

	return nil
}

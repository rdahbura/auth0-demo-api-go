package clientcredentials

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"dahbura.me/api/config"
)

var httpClient = http.Client{
	Timeout: config.DefaultClientTimeout,
}

type AccessTokenSource struct {
	atreq *AccessTokenRequest
	atres *AccessTokenResponse
	mtx   sync.Mutex
}

type AccessTokenRequest struct {
	ClientId     string
	ClientSecret string
	Audience     string
	Url          string
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`

	ExpiresAt time.Time
}

func NewSource(atreq *AccessTokenRequest) *AccessTokenSource {
	atres := AccessTokenResponse{}
	atsrc := AccessTokenSource{
		atreq: atreq,
		atres: &atres,
		mtx:   sync.Mutex{},
	}

	return &atsrc
}

func (atsrc *AccessTokenSource) Token() (string, error) {
	atsrc.mtx.Lock()
	defer atsrc.mtx.Unlock()

	if atsrc.atres.hasExpired() {
		atres, err := atsrc.atreq.do()
		if err != nil {
			return "", err
		}

		atsrc.atres = atres
	}

	return atsrc.atres.AccessToken, nil
}

func NewRequest(tokenUrl string, clientId string, clientSecret string, audience string) (*AccessTokenRequest, error) {
	_, err := url.Parse(tokenUrl)
	if err != nil {
		return nil, err
	}

	if clientId == "" {
		return nil, errors.New("clientId required")
	}

	if clientSecret == "" {
		return nil, errors.New("clientSecret required")
	}

	if audience == "" {
		return nil, errors.New("audience required")
	}

	atreq := AccessTokenRequest{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Audience:     audience,
		Url:          tokenUrl,
	}

	return &atreq, nil
}

func (atreq *AccessTokenRequest) do() (*AccessTokenResponse, error) {
	atenc := atreq.encode()
	atreader := strings.NewReader(atenc)

	req, err := http.NewRequest(http.MethodPost, atreq.Url, atreader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", config.MimeApplicationXWwwFormUrlencoded)

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var atres AccessTokenResponse
	if err := json.Unmarshal(body, &atres); err != nil {
		return nil, err
	}

	now := time.Now()
	expiresIn := time.Duration(atres.ExpiresIn) * time.Second
	expiresAt := now.Add(expiresIn)

	atres.ExpiresAt = expiresAt

	return &atres, nil
}

func (atreq *AccessTokenRequest) encode() string {
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", atreq.ClientId)
	values.Set("client_secret", atreq.ClientSecret)
	values.Set("audience", atreq.Audience)

	return values.Encode()
}

func (atres *AccessTokenResponse) hasExpired() bool {
	now := time.Now()
	expiresAt := atres.ExpiresAt.Add(-config.DefaultTokenLeeway)
	hasExpired := now.After(expiresAt)

	return hasExpired
}

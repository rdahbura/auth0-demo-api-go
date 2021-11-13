package oauth2

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"dahbura.me/api/config"
	httppkg "dahbura.me/api/util/http"
)

type ClientCredentialsRequest struct {
	ClientId     string
	ClientSecret string
	Audience     string
	Scopes       []string
	TokenUrl     string
}

func (ccreq *ClientCredentialsRequest) encode() string {
	val := url.Values{}
	val.Set("grant_type", "client_credentials")
	val.Set("client_id", ccreq.ClientId)
	val.Set("client_secret", ccreq.ClientSecret)
	val.Set("audience", ccreq.Audience)
	if len(ccreq.Scopes) > 0 {
		val.Set("scope", strings.Join(ccreq.Scopes, " "))
	}

	return val.Encode()
}

type ClientCredentialsSource struct {
	ccreq ClientCredentialsRequest
	gres  GrantResponse
	mtx   sync.Mutex
}

func NewClientCredentialsSource(ccreq ClientCredentialsRequest) *ClientCredentialsSource {
	gres := GrantResponse{}
	src := &ClientCredentialsSource{
		ccreq: ccreq,
		gres:  gres,
	}

	return src
}

func (ccsrc *ClientCredentialsSource) GrantResponse() (GrantResponse, error) {
	ccsrc.mtx.Lock()
	defer ccsrc.mtx.Unlock()

	if !ccsrc.gres.HasExpired() {
		return ccsrc.gres, nil
	}

	gres, err := newClientCredentialsResponse(ccsrc.ccreq)
	if err != nil {
		return GrantResponse{}, err
	}

	ccsrc.gres = gres

	return ccsrc.gres, nil
}

func newClientCredentialsResponse(ccreq ClientCredentialsRequest) (GrantResponse, error) {
	url := ccreq.TokenUrl
	encoded := ccreq.encode()
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(encoded))
	if err != nil {
		return GrantResponse{}, err
	}

	req.Header.Set("Content-Type", config.MimeApplicationXWwwFormUrlencoded)

	body, err := httppkg.DoRequest(req)
	if err != nil {
		return GrantResponse{}, err
	}

	var gr GrantResponse
	if err := json.Unmarshal(body, &gr); err != nil {
		return GrantResponse{}, err
	}

	expiresAt := ExpiresAt(gr.ExpiresIn)

	gres := GrantResponse{
		AccessToken: gr.AccessToken,
		TokenType:   gr.TokenType,
		ExpiresIn:   gr.ExpiresIn,
		ExpiresAt:   expiresAt,
	}

	return gres, nil
}

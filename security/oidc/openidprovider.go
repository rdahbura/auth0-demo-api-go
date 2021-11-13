package oidc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dahbura.me/api/config"
)

var httpClient = http.Client{
	Timeout: config.DefaultClientTimeout,
}

type OpenIdProviderConfig struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	DeviceAuthorizationEndpoint       string   `json:"device_authorization_endpoint"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
	MfaChallengeEndpoint              string   `json:"mfa_challenge_endpoint"`
	JwksUri                           string   `json:"jwks_uri"`
	RegistrationEndpoint              string   `json:"registration_endpoint"`
	RevocationEndpoint                string   `json:"revocation_endpoint"`
	ScopesSupported                   []string `json:"scopes_supported"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
	ResponseModesSupported            []string `json:"response_modes_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	ClaimsSupported                   []string `json:"claims_supported"`
	RequestURIParameterSupported      bool     `json:"request_uri_parameter_supported"`
}

func ReadOpenIdProviderConfig(issuer string) (*OpenIdProviderConfig, error) {
	url := fmt.Sprintf("%s/.well-known/openid-configuration", issuer)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	decoder.DisallowUnknownFields()

	var config OpenIdProviderConfig
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

package jose

import (
	"encoding/json"
	"net/http"

	"dahbura.me/api/config"
)

var httpClient = http.Client{
	Timeout: config.DefaultClientTimeout,
}

func fetchEncodedDer(jwksUrl string, kid string) (string, error) {
	jwks, err := readJwkSet(jwksUrl)
	if err != nil {
		return "", err
	}

	var encodedDer string
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			encodedDer = key.X5C[0]
			break
		}
	}

	return encodedDer, nil
}

func readJwkSet(jwksUrl string) (*JwkSet, error) {
	req, err := http.NewRequest(http.MethodGet, jwksUrl, nil)
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

	var jwks JwkSet
	if err := decoder.Decode(&jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}

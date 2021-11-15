package jose

import (
	"encoding/json"
	"net/http"
	"sync"

	"dahbura.me/api/config"
)

var (
	cache      = map[string]string{}
	cacheMutex = sync.Mutex{}
)

var (
	httpClient = http.Client{
		Timeout: config.DefaultClientTimeout,
	}
)

func fetchEncodedDer(jwksUrl string, kid string) (string, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	encodedDer, ok := cache[jwksUrl]
	if ok {
		return encodedDer, nil
	}

	jwks, err := readJwkSet(jwksUrl)
	if err != nil {
		return "", err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			encodedDer = key.X5C[0]
			break
		}
	}

	cache[jwksUrl] = encodedDer

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

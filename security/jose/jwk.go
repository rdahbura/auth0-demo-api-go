package jose

import (
	"encoding/json"
	"net/http"
	"time"

	"dahbura.me/api/config"
	"dahbura.me/api/util/cache"
)

var memoryCache = cache.New()

var httpClient = http.Client{
	Timeout: config.DefaultClientTimeout,
}

func fetchEncodedDer(jwksUrl string, kid string) (string, error) {
	encodedDer, ok := memoryCache.Get(kid)
	if ok {
		return encodedDer.(string), nil
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

	item := cache.Item{
		Key:   kid,
		Value: encodedDer,
	}

	now := time.Now()
	itemPolicy := cache.ItemPolicy{
		AbsoluteExp: now.Add(time.Second * 10),
	}

	memoryCache.Set(item, itemPolicy)

	return encodedDer.(string), nil
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

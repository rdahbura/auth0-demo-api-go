package http

import (
	"net/http"
	"sync"

	"dahbura.me/api/config"
)

var (
	httpClient     *http.Client
	httpClientOnce sync.Once
)

func GetHttpClient() *http.Client {
	httpClientOnce.Do(func() {
		initHttpClient()
	})

	return httpClient
}

func initHttpClient() {
	httpClient = &http.Client{
		Timeout: config.DefaultClientTimeout,
	}
}

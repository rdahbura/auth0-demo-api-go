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
	httpClientOnce.Do(initHttpClient)

	return httpClient
}

func initHttpClient() {
	t := http.DefaultTransport.(*http.Transport).Clone()

	httpClient = &http.Client{
		Timeout:   config.DefaultClientTimeout,
		Transport: t,
	}
}

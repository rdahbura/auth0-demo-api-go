package security

import (
	"sync"

	"dahbura.me/api/config"
	"dahbura.me/api/security/oauth2/clientcredentials"
)

var (
	atsrc  *clientcredentials.AccessTokenSource
	atonce sync.Once
)

func GetClientCredSrc() *clientcredentials.AccessTokenSource {
	atonce.Do(func() {
		atreq := &clientcredentials.AccessTokenRequest{
			ClientId:     config.MgmtApiClientId,
			ClientSecret: config.MgmtApiClientSecret,
			Audience:     config.MgmtApiAudience,
			Url:          config.MgmtApiTokenUrl,
		}

		atsrc = clientcredentials.NewSource(atreq)
	})

	return atsrc
}

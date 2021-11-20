package security

import (
	"sync"

	"dahbura.me/api/config"
	"dahbura.me/api/security/oauth2/clientcredentials"
)

var (
	grant     *clientcredentials.ClientCredGrant
	grantOnce sync.Once
)

func GetClientCredGrant() *clientcredentials.ClientCredGrant {
	grantOnce.Do(func() {
		atreq := &clientcredentials.AccessTokenRequest{
			ClientId:     config.MgmtApiClientId,
			ClientSecret: config.MgmtApiClientSecret,
			Audience:     config.MgmtApiAudience,
			Url:          config.MgmtApiTokenUrl,
		}

		grant = clientcredentials.NewClientCredGrant(atreq)
	})

	return grant
}

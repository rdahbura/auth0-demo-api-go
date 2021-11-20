package security

import (
	"sync"

	"dahbura.me/api/config"
	clientCredentials "dahbura.me/api/security/oauth2/client_credentials"
)

var (
	grant     *clientCredentials.ClientCredGrant
	grantOnce sync.Once
)

func GetClientCredGrant() *clientCredentials.ClientCredGrant {
	grantOnce.Do(func() {
		initClientCredGrant()
	})

	return grant
}

func initClientCredGrant() {
	atreq := &clientCredentials.AccessTokenRequest{
		ClientId:     config.MgmtApiClientId,
		ClientSecret: config.MgmtApiClientSecret,
		Audience:     config.MgmtApiAudience,
		Url:          config.MgmtApiTokenUrl,
	}

	grant = clientCredentials.NewClientCredGrant(atreq)
}

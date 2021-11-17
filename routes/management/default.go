package management

import (
	"dahbura.me/api/config"
	"dahbura.me/api/security/oauth2"
)

var (
	clientCredSrc *oauth2.ClientCredentialsSource
)

func Startup() {
	clientCredReq := oauth2.ClientCredentialsRequest{
		ClientId:     config.MgmtApiClientId,
		ClientSecret: config.MgmtApiClientSecret,
		Audience:     config.MgmtApiAudience,
		Scopes:       []string{},
		TokenUrl:     config.MgmtApiTokenUrl,
	}

	clientCredSrc = oauth2.NewClientCredentialsSource(clientCredReq)
}

package management

import (
	"net/http"

	"dahbura.me/api/config"
	"dahbura.me/api/middleware"
	"dahbura.me/api/security/oauth2"

	"github.com/gin-gonic/gin"
)

var (
	clientCredSrc *oauth2.ClientCredentialsSource
)

func ApplyRoutes(rg *gin.RouterGroup) {
	clientCredReq := oauth2.ClientCredentialsRequest{
		ClientId:     config.MgmtApiClientId,
		ClientSecret: config.MgmtApiClientSecret,
		Audience:     config.MgmtApiAudience,
		Scopes:       []string{},
		TokenUrl:     config.MgmtApiTokenUrl,
	}

	clientCredSrc = oauth2.NewClientCredentialsSource(clientCredReq)

	checkJwtOptions := middleware.CheckJwtOptions{
		TokenAudience: config.TokenAudience,
		TokenIssuer:   config.TokenIssuer + "/",
	}

	checkScopeOptions := middleware.CheckScopeOptions{
		ScopesClaim: "permissions",
	}

	checkJwt := middleware.CheckJwt(checkJwtOptions)
	checkScope := middleware.CheckScope(checkScopeOptions)

	grp := rg.Group("/mgmt", checkJwt())
	grp.Handle(http.MethodGet, "clients", checkScope("read:clients"), getClients)
	grp.Handle(http.MethodGet, "clients/:id", checkScope("read:clients"), getClient)
	grp.Handle(http.MethodGet, "connections", checkScope("read:connections"), getConnections)
	grp.Handle(http.MethodGet, "connections/:id", checkScope("read:connections"), getConnection)
	grp.Handle(http.MethodGet, "users", checkScope("read:users"), getUsers)
	grp.Handle(http.MethodGet, "users/:id", checkScope("read:users"), getUser)
	grp.Handle(http.MethodPatch, "users/:id", checkScope("update:users"), patchUser)
}

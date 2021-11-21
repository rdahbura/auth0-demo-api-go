package routes

import (
	"net/http"

	"dahbura.me/api/config"
	"dahbura.me/api/middleware"
	"dahbura.me/api/routes/database"
	"dahbura.me/api/routes/management"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	checkJwtOptions := middleware.CheckJwtOptions{
		TokenAudience: config.TokenAudience,
		TokenIssuer:   config.TokenIssuer + "/",
	}
	checkJwt := middleware.CheckJwt(checkJwtOptions)

	checkScopeOptions := middleware.CheckScopeOptions{
		ScopesClaim: "permissions",
	}
	checkScope := middleware.CheckScope(checkScopeOptions)

	rg := router.Group("/")
	{
		rg.Handle(http.MethodGet, "/", rootHandler)
	}

	rgDb := rg.Group("/db", checkJwt())
	{
		rgDb.Handle(http.MethodPost, "logins", database.Logins)
		rgDb.Handle(http.MethodGet, "users", checkScope("read:users"), database.GetUsers)
		rgDb.Handle(http.MethodGet, "users/:id", checkScope("read:users"), database.GetUser)
		rgDb.Handle(http.MethodPost, "users", checkScope("create:users"), database.CreateUser)
		rgDb.Handle(http.MethodDelete, "users/:id", checkScope("delete:users"), database.DeleteUser)
		rgDb.Handle(http.MethodPatch, "users/:id", checkScope("update:users"), database.UpdateUser)
	}

	rgMgmt := rg.Group("/mgmt", checkJwt())
	{
		rgMgmt.Handle(http.MethodGet, "clients", checkScope("read:clients"), management.GetClients)
		rgMgmt.Handle(http.MethodGet, "clients/:id", checkScope("read:clients"), management.GetClient)
		rgMgmt.Handle(http.MethodGet, "connections", checkScope("read:connections"), management.GetConnections)
		rgMgmt.Handle(http.MethodGet, "connections/:id", checkScope("read:connections"), management.GetConnection)
		rgMgmt.Handle(http.MethodGet, "users", checkScope("read:users"), management.GetUsers)
		rgMgmt.Handle(http.MethodGet, "users/:id", checkScope("read:users"), management.GetUser)
		rgMgmt.Handle(http.MethodPatch, "users/:id", checkScope("update:users"), management.PatchUser)
	}
}

func rootHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

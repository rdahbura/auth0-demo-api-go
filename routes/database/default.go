package database

import (
	"net/http"

	"dahbura.me/api/config"
	"dahbura.me/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ApplyRoutes(rg *gin.RouterGroup) {
	validate = validator.New()

	checkJwtOptions := middleware.CheckJwtOptions{
		TokenAudience: config.TokenAudience,
		TokenIssuer:   config.TokenIssuer + "/",
	}

	checkScopeOptions := middleware.CheckScopeOptions{
		ScopesClaim: "permissions",
	}

	checkJwt := middleware.CheckJwt(checkJwtOptions)
	checkScope := middleware.CheckScope(checkScopeOptions)

	grp := rg.Group("/db", checkJwt())
	{
		grp.Handle(http.MethodPost, "logins", logins)
		grp.Handle(http.MethodGet, "users", checkScope("read:users"), getUsers)
		grp.Handle(http.MethodGet, "users/:id", checkScope("read:users"), getUser)
		grp.Handle(http.MethodPost, "users", checkScope("create:users"), createUser)
		grp.Handle(http.MethodDelete, "users/:id", checkScope("delete:users"), deleteUser)
		grp.Handle(http.MethodPatch, "users/:id", checkScope("update:users"), updateUser)
	}
}

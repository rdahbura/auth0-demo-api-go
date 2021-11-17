package routes

import (
	"net/http"

	"dahbura.me/api/config"
	"dahbura.me/api/middleware"
	"dahbura.me/api/routes/database"
	"dahbura.me/api/routes/management"

	"github.com/gin-gonic/gin"
)

var (
	checkJwt          func() gin.HandlerFunc
	checkJwtOptions   middleware.CheckJwtOptions
	checkScope        func(string) gin.HandlerFunc
	checkScopeOptions middleware.CheckScopeOptions
)

type Route struct {
	Method        string
	Path          string
	Handler       gin.HandlerFunc
	Authorization Authorization
	Routes        []Route
}

type Authorization struct {
	CheckJwt   func() func() gin.HandlerFunc
	CheckScope func() func(string) gin.HandlerFunc
	Scope      string
}

var routes = Route{
	Path: "/",
	Routes: []Route{
		{
			Path: "/db",
			Authorization: Authorization{
				CheckJwt: getCheckJwt,
			},
			Routes: []Route{
				{
					Method:  http.MethodPost,
					Path:    "/logins",
					Handler: database.Logins,
				},
				{
					Method:  http.MethodGet,
					Path:    "/users",
					Handler: database.GetUsers,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "read:users",
					},
				},
				{
					Method:  http.MethodGet,
					Path:    "/users/:id",
					Handler: database.GetUser,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "read:users",
					},
				},
				{
					Method:  http.MethodPost,
					Path:    "/users",
					Handler: database.CreateUser,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "create:users",
					},
				},
				{
					Method:  http.MethodDelete,
					Path:    "/users/:id",
					Handler: database.DeleteUser,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "delete:users",
					},
				},
				{
					Method:  http.MethodPatch,
					Path:    "/users/:id",
					Handler: database.UpdateUser,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "update:users",
					},
				},
			},
		},
		{
			Path: "/mgmt",
			Authorization: Authorization{
				CheckJwt: getCheckJwt,
			},
			Routes: []Route{
				{
					Method:  http.MethodGet,
					Path:    "/clients",
					Handler: management.GetClients,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "read:clients",
					},
				},
				{
					Method:  http.MethodGet,
					Path:    "/clients/:id",
					Handler: management.GetClient,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "read:clients",
					},
				},
				{
					Method:  http.MethodGet,
					Path:    "/connections",
					Handler: management.GetConnections,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "read:connections",
					},
				},
				{
					Method:  http.MethodGet,
					Path:    "/connections/:id",
					Handler: management.GetConnection,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "read:connections",
					},
				},
				{
					Method:  http.MethodGet,
					Path:    "/users",
					Handler: management.GetUsers,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "read:users",
					},
				},
				{
					Method:  http.MethodGet,
					Path:    "/users/:id",
					Handler: management.GetUser,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "read:users",
					},
				},
				{
					Method:  http.MethodPatch,
					Path:    "/users/:id",
					Handler: management.PatchUser,
					Authorization: Authorization{
						CheckScope: getCheckScope,
						Scope:      "update:users",
					},
				},
			},
		},
	},
}

func ApplyRoutes(eng *gin.Engine) {
	checkJwtOptions = middleware.CheckJwtOptions{
		TokenAudience: config.TokenAudience,
		TokenIssuer:   config.TokenIssuer + "/",
	}
	checkJwt = middleware.CheckJwt(checkJwtOptions)

	checkScopeOptions = middleware.CheckScopeOptions{
		ScopesClaim: "permissions",
	}
	checkScope = middleware.CheckScope(checkScopeOptions)

	setupRoutes(eng.Group("/"), routes)

	management.Setup()
}

func getCheckJwt() func() gin.HandlerFunc {
	return checkJwt
}

func getCheckScope() func(string) gin.HandlerFunc {
	return checkScope
}

func setupHandlers(route Route) []gin.HandlerFunc {
	handlers := []gin.HandlerFunc{}

	if route.Authorization.CheckJwt != nil {
		handlers = append(handlers, route.Authorization.CheckJwt()())
	}

	if route.Authorization.CheckScope != nil {
		scope := route.Authorization.Scope
		handlers = append(handlers, route.Authorization.CheckScope()(scope))
	}

	if route.Handler != nil {
		handlers = append(handlers, route.Handler)
	}

	return handlers
}

func setupRoutes(parentGrp *gin.RouterGroup, route Route) {
	if len(route.Routes) > 0 {
		handlers := setupHandlers(route)

		grp := parentGrp.Group(route.Path)
		grp.Use(handlers...)

		for _, r := range route.Routes {
			setupRoutes(grp, r)
		}

		return
	}

	handlers := setupHandlers(route)

	parentGrp.Handle(route.Method, route.Path, handlers...)
}

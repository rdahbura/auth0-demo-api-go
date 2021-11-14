package routes

import (
	"net/http"

	"dahbura.me/api/routes/database"
	"dahbura.me/api/routes/management"

	"github.com/gin-gonic/gin"
)

func ApplyRoutes(eng *gin.Engine) {
	grp := eng.Group("/")
	grp.Handle(http.MethodGet, "/", rootHandler)

	database.ApplyRoutes(grp)
	management.ApplyRoutes(grp)
}

func rootHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

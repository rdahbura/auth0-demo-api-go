package management

import (
	"fmt"
	"net/http"

	"dahbura.me/api/config"
	"dahbura.me/api/security/oauth2"
	httppkg "dahbura.me/api/util/http"

	"github.com/gin-gonic/gin"
)

func GetConnections(c *gin.Context) {
	url := fmt.Sprintf("%s/connections", config.MgmtApiBaseUrl)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if httppkg.HandleError(c, err) {
		return
	}

	err = oauth2.SetAuthHeader(req, clientCredSrc)
	if httppkg.HandleError(c, err) {
		return
	}

	body, err := httppkg.DoRequest(req)
	if httppkg.HandleError(c, err) {
		return
	}

	c.Header("Content-Type", config.MimeApplicationJson)
	c.String(http.StatusOK, string(body))
}

func GetConnection(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("%s/connections/%s", config.MgmtApiBaseUrl, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if httppkg.HandleError(c, err) {
		return
	}

	err = oauth2.SetAuthHeader(req, clientCredSrc)
	if httppkg.HandleError(c, err) {
		return
	}

	body, err := httppkg.DoRequest(req)
	if httppkg.HandleError(c, err) {
		return
	}

	c.Header("Content-Type", config.MimeApplicationJson)
	c.String(http.StatusOK, string(body))
}

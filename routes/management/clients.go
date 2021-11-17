package management

import (
	"fmt"
	"net/http"

	"dahbura.me/api/config"
	"dahbura.me/api/routes/management/security"
	"dahbura.me/api/security/oauth2"
	httppkg "dahbura.me/api/util/http"

	"github.com/gin-gonic/gin"
)

func GetClients(c *gin.Context) {
	url := fmt.Sprintf("%s/clients", config.MgmtApiBaseUrl)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if httppkg.HandleError(c, err) {
		return
	}

	err = oauth2.SetAuthHeader(req, security.GetClientCredSrc())
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

func GetClient(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("%s/clients/%s", config.MgmtApiBaseUrl, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if httppkg.HandleError(c, err) {
		return
	}

	err = oauth2.SetAuthHeader(req, security.GetClientCredSrc())
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

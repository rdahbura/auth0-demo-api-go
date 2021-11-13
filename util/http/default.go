package http

import (
	"io"
	"net/http"

	"dahbura.me/api/config"

	"github.com/gin-gonic/gin"
)

var httpClient = http.Client{
	Timeout: config.DefaultClientTimeout,
}

func DoRequest(req *http.Request) ([]byte, error) {
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func HandleError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return true
	}

	return false
}

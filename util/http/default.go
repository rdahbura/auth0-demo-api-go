package http

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

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

func HandleErrorMiddleware(c *gin.Context, err error) bool {
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
		c.Error(err)
		return true
	}

	return false
}

func TokenFromContext(c *gin.Context) (string, error) {
	tokenValue, exists := c.Get(config.ContextBearerToken)
	if !exists {
		return "", errors.New("unable to get bearer token")
	}

	token, ok := tokenValue.(string)
	if !ok {
		return "", errors.New("unable to convert token to string")
	}

	return token, nil
}

func TokenFromHeader(header string) (string, error) {
	if len(header) == 0 {
		return "", errors.New("authorization header is missing")
	}

	headerSegments := strings.Split(header, " ")
	if len(headerSegments) != 2 {
		return "", errors.New("authorization header segment count is incorrect")
	}

	schemeSegment := headerSegments[0]
	if !strings.EqualFold(schemeSegment, "Bearer") {
		return "", errors.New("authorization scheme is missing")
	}

	tokenSegment := headerSegments[1]

	return tokenSegment, nil
}

func VerifyScope(token string, scopesClaim string, scope string) error {
	encodedPayload := strings.Split(token, ".")[1]
	decodedPayload, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{}
	err = json.Unmarshal(decodedPayload, &payload)
	if err != nil {
		return err
	}

	err = hasScope(payload, scopesClaim, scope)
	if err != nil {
		return err
	}

	return nil
}

func hasScope(payload map[string]interface{}, scopesClaim string, scope string) error {
	scopes, exists := payload[scopesClaim]
	if !exists {
		return errors.New("scopes claim not found")
	}

	scopesArray, ok := scopes.([]interface{})
	if !ok {
		return errors.New("scopes claim improperly formatted")
	}

	_, ok = toMap(scopesArray)[scope]
	if !ok {
		return errors.New("scope not found")
	}

	return nil
}

func toMap(values []interface{}) map[string]int {
	m := map[string]int{}
	for i, v := range values {
		m[v.(string)] = i
	}

	return m
}

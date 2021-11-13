package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"dahbura.me/api/config"
	"dahbura.me/api/security/jose"

	"github.com/gin-gonic/gin"
)

type CheckJwtOptions struct {
	TokenAudience string
	TokenIssuer   string
}

func CheckJwt(options CheckJwtOptions) func() gin.HandlerFunc {
	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			header := c.GetHeader("Authorization")

			token, err := parseBearerToken(header)
			if handleError(c, err) {
				return
			}

			jwksUrl := fmt.Sprintf("%s/.well-known/jwks.json", strings.TrimSuffix(options.TokenIssuer, "/"))

			err = jose.VerifyCompact(jwksUrl, token, options.TokenIssuer, options.TokenAudience)
			if handleError(c, err) {
				return
			}

			c.Set(config.ContextBearerToken, token)
		}
	}
}

type CheckScopeOptions struct {
	ScopesClaim string
}

func CheckScope(options CheckScopeOptions) func(string) gin.HandlerFunc {
	return func(scope string) gin.HandlerFunc {
		return func(c *gin.Context) {
			token, err := getBearerToken(c)
			if handleError(c, err) {
				return
			}

			encodedPayload := strings.Split(token, ".")[1]
			decodedPayload, err := base64.RawURLEncoding.DecodeString(encodedPayload)
			if handleError(c, err) {
				return
			}

			payload := map[string]interface{}{}
			err = json.Unmarshal(decodedPayload, &payload)
			if handleError(c, err) {
				return
			}

			_, err = hasScope(payload, options.ScopesClaim, scope)
			if handleError(c, err) {
				return
			}
		}
	}
}

func getBearerToken(c *gin.Context) (string, error) {
	v, exists := c.Get(config.ContextBearerToken)
	if !exists {
		return "", errors.New("unable to get bearer token")
	}

	s, ok := v.(string)
	if !ok {
		return "", errors.New("unable to convert token to string")
	}

	return s, nil
}

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
		c.Error(err)
		return true
	}

	return false
}

func hasScope(payload map[string]interface{}, scopesClaim string, scope string) (bool, error) {
	v, exists := payload[scopesClaim]
	if !exists {
		return false, errors.New("claim is missing")
	}

	i, ok := v.([]interface{})
	if !ok {
		return false, errors.New("claim is improperly formatted")
	}

	_, ok = toMap(i)[scope]
	if !ok {
		return false, errors.New("claim does not contain scope")
	}

	return true, nil
}

func parseBearerToken(header string) (string, error) {
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

func toMap(values []interface{}) map[string]int {
	m := map[string]int{}
	for i, v := range values {
		m[v.(string)] = i
	}

	return m
}

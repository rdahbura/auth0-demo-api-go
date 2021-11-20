package middleware

import (
	"dahbura.me/api/config"
	"dahbura.me/api/security/jose"
	httppkg "dahbura.me/api/util/http"

	"github.com/gin-gonic/gin"
)

type CheckJwtOptions struct {
	TokenAudience string
	TokenIssuer   string
}

func CheckJwt(options CheckJwtOptions) func() gin.HandlerFunc {
	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			token, err := httppkg.TokenFromHeader(c)
			if httppkg.HandleErrorMiddleware(c, err) {
				return
			}

			err = jose.VerifyCompact(token, options.TokenIssuer, options.TokenAudience)
			if httppkg.HandleErrorMiddleware(c, err) {
				return
			}

			c.Set(config.ContextBearerToken, token)
		}
	}
}

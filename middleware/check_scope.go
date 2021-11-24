package middleware

import (
	httppkg "dahbura.me/api/util/http"

	"github.com/gin-gonic/gin"
)

type CheckScopeOpts struct {
	ScopesClaim string
}

func CheckScope(opts CheckScopeOpts) func(string) gin.HandlerFunc {
	return func(scope string) gin.HandlerFunc {
		return func(c *gin.Context) {
			token, err := httppkg.TokenFromContext(c)
			if httppkg.HandleErrorMiddleware(c, err) {
				return
			}

			err = httppkg.VerifyScope(token, opts.ScopesClaim, scope)
			if httppkg.HandleErrorMiddleware(c, err) {
				return
			}
		}
	}
}

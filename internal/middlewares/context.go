package middlewares

import (
	"os"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/poligonoio/vega-core/pkg/env"
)

func SetVarsToContext() gin.HandlerFunc {
	return func(c *gin.Context) {

		enableAuth := env.GetBoolEnv("ENABLE_AUTHENTICATION")
		authType := os.Getenv("AUTHENTICATION_TYPE")

		if enableAuth && authType == "oauth2" {
			claims := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			customClaims := claims.CustomClaims.(*CustomClaims)

			sub := strings.Split(claims.RegisteredClaims.Subject, "|")[1]

			if customClaims.OrganizationId != "" {
				c.Set("owner_id", customClaims.OrganizationId)
				c.Set("sub", sub)
			} else {
				c.Set("owner_id", sub)
				c.Set("sub", sub)
			}
		} else {
			c.Set("owner_id", "poligono")
			c.Set("sub", "poligono")
		}

		c.Next()
	}
}

package middlewares

import (
	"os"

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

			if customClaims.OrganizationId != "" {
				c.Set("owner_id", customClaims.OrganizationId)
				c.Set("sub", claims.RegisteredClaims.Subject)
			} else {
				c.Set("owner_id", claims.RegisteredClaims.Subject)
				c.Set("sub", claims.RegisteredClaims.Subject)
			}
		} else {
			c.Set("owner_id", "poligono")
			c.Set("sub", "poligono")
		}

		c.Next()
	}
}

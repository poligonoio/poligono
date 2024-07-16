package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/poligonoio/vega-core/pkg/logger"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope          string   `json:"scope"`
	OrganizationId string   `json:"org_id"`
	Permissions    []string `json:"permissions"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func EnsureValidToken() gin.HandlerFunc {
	issuerURL, err := url.Parse(os.Getenv("OAUTH2_ISSUER"))
	logger.Info.Println(issuerURL.String())

	if err != nil {
		logger.Error.Printf("Failed to parse the issuer url: %v\n", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("OAUTH2_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)

	if err != nil {
		logger.Error.Printf("Failed to set up the jwt validator: %v\n", err)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Warning.Printf("Encountered error while validating JWT: %v\n", err)
		fmt.Printf("Encountered error while validating JWT: %v\n", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte(`{"error": "unauthorized", "description":"Failed to validate JWT"}`))
		if err != nil {
			logger.Error.Fatalln(fmt.Errorf("Error writing http: %v", err))
		}
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(c *gin.Context) {
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Next()
		}

		middleware.CheckJWT(handler).ServeHTTP(c.Writer, c.Request)
	}
}

// HasScope checks whether our claims have a specific scope.
func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

func (c CustomClaims) HasPermission(expectedPermission string) bool {
	for i := range c.Permissions {
		if c.Permissions[i] == expectedPermission {
			return true
		}
	}

	return false
}

func (c CustomClaims) HasPermissions(sub []string) bool {
	subLen := len(sub)
	if subLen == 0 {
		return true
	}

	permissionsLen := len(c.Permissions)

	for i := 0; i <= permissionsLen-subLen; i++ {
		j := 0
		for ; j < subLen; j++ {
			if c.Permissions[i+j] != sub[j] {
				break
			}
		}
		if j == subLen {
			return true
		}
	}
	return false
}

func EnsureValidRole() gin.HandlerFunc {
	n := make(map[string][]string)

	n["/v1/datasources/all GET"] = append(n["/v1/datasources/all"], "read:datasources")
	n["/v1/datasources/:name GET"] = append(n["/v1/datasources/:name"], "read:datasources")
	n["/v1/datasources POST"] = append(n["/v1/datasources/:name"], "write:datasources")
	n["/v1/datasources/:name PUT"] = append(n["/v1/datasources/:name"], "put:datasources")
	n["/v1/datasources/:name DELETE"] = append(n["/v1/datasources/:name"], "delete:datasources")

	n["/v1/prompt POST"] = append(n["/v1/prompt"], "read:datasources")
	n["/v1/prompt POST"] = append(n["/v1/prompt"], "read:rawdata")

	return func(c *gin.Context) {
		claims := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		customClaims := claims.CustomClaims.(*CustomClaims)

		request := fmt.Sprintf("%s %s", c.FullPath(), c.Request.Method)

		if !customClaims.HasPermissions(n[request]) {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				map[string]string{
					"error":       "forbidden",
					"description": fmt.Sprintf("Your current permissions do not allow you to perform this action. If you need access, please request the '%s' permission from your administrator.", n[request]),
				},
			)

			return
		}

		c.Next()
	}
}

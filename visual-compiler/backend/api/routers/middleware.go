package routers

import (
	"net/http"
	"os"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Name: Auth0Middleware
//
// Parameters: none
//
// Return: Function
//
// Authorizes middleware with the users auth token
func Auth0MiddleWare() gin.HandlerFunc {
	auth_domain := os.Getenv("AUTH0_DOMAIN")

	jwksURL := auth_domain + "/.well-known/jwks.json"
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		panic("Creating JWKS from Auth0 failed: " + err.Error())
	}

	return func(c *gin.Context) {
		header_auth := c.GetHeader("Authorization")
		if header_auth == "" || !strings.HasPrefix(header_auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorzation"})
			c.Abort()
			return
		}

		token_string := strings.TrimPrefix(header_auth, "Bearer ")

		token, err := jwt.Parse(token_string, jwks.Keyfunc)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		jwt_claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims invalid"})
			c.Abort()
			return
		}

		aud_claim, ok := jwt_claims["aud"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No audience in token"})
			c.Abort()
			return
		}

		aud_valid := os.Getenv("CUSTOM_API")
		aud_is_valid := false

		switch audience := aud_claim.(type) {
		case string:
			if strings.Contains(audience, aud_valid) {
				aud_is_valid = true
			}
		case []interface{}:
			for _, a := range audience {
				if contained, ok := a.(string); ok && strings.Contains(contained, aud_valid) {
					aud_is_valid = true
					break
				}
			}
		}

		if !aud_is_valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid audience in token"})
			c.Abort()
			return
		}

		users_sub, ok := jwt_claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No 'sub' claim in token"})
			c.Abort()
			return
		}

		c.Set("auth0_id", users_sub)

		c.Next()
	}
}

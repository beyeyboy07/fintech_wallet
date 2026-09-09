package middleware

import (
	"net/http"
	"strings"

	"fintech-wallet/backend/internal/auth"
	"fintech-wallet/backend/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Error(
				c,
				http.StatusUnauthorized,
				"Authorization header is required",
				"UNAUTHORIZED",
			)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(
				c,
				http.StatusUnauthorized,
				"Invalid authorization header",
				"UNAUTHORIZED",
			)
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwtManager.ValidateToken(tokenString)

		if err != nil {
			response.Error(
				c,
				http.StatusUnauthorized,
				"Invalid or expired token",
				"UNAUTHORIZED",
			)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			response.Error(
				c,
				http.StatusUnauthorized,
				"Invalid token claims",
				"UNAUTHORIZED",
			)
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

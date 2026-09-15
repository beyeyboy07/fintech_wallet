package middleware

import (
	"net/http"

	"fintech-wallet/backend/internal/response"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleValue, exists := c.Get("role")

		if !exists {
			response.Error(
				c,
				http.StatusForbidden,
				"User role is not available",
				"FORBIDDEN",
			)
			c.Abort()
			return
		}

		role, ok := roleValue.(string)

		if !ok {
			response.Error(
				c,
				http.StatusForbidden,
				"Invalid user role",
				"FORBIDDEN",
			)
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		response.Error(
			c,
			http.StatusForbidden,
			"Insufficient permissions",
			"FORBIDDEN",
		)

		c.Abort()
	}
}

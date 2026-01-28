package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"thank-you-card-backend/internal/models"
	"thank-you-card-backend/internal/repositories"
	"thank-you-card-backend/internal/services"
)

// Context keys
const (
	ContextCurrentUser     = "currentUser"
	ContextCurrentEmployee = "currentEmployee"
)

// AuthMiddleware validates JWT token and attaches user to context
func AuthMiddleware(authService services.AuthService, userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		tokenString := parts[1]
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Load user from database
		user, err := userRepo.GetByID(c.Request.Context(), claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		c.Set(ContextCurrentUser, user)
		if user.Employee != nil {
			c.Set(ContextCurrentEmployee, user.Employee)
		} else if user.EmployeeID != nil {
			// Load employee if not preloaded
			dbVal, exists := c.Get("db")
			if exists {
				if db, ok := dbVal.(*gorm.DB); ok {
					var emp models.Employee
					if err := db.WithContext(c.Request.Context()).First(&emp, *user.EmployeeID).Error; err == nil {
						c.Set(ContextCurrentEmployee, &emp)
					} else {
						// Log error but don't fail the request - some users might not have employees
						// This is acceptable for admin users or special cases
					}
				}
			}
		}
		c.Next()
	}
}

// RequireHRAdmin checks if user has HR or ADMIN role
func RequireHRAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(ContextCurrentUser)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		user, ok := val.(*models.User)
		if !ok || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Check if user has HR or ADMIN role
		if user.Role != "HR" && user.Role != "ADMIN" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: hr admin only"})
			return
		}

		c.Next()
	}
}

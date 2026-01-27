package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	cardApp "github.com/company/thank-you-card/internal/application/card"
	"github.com/company/thank-you-card/pkg/logger"
)

// AuthMiddleware provides authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For development/testing, we'll use a simple mock authentication
		// In production, this would validate JWT tokens or integrate with Azure AD
		
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// For development, allow requests without auth and set a default user
			logger.Debug("No authorization header, using default user for development")
			c.Set("user_id", "emp001")
			c.Set("user_role", "Employee")
			c.Next()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Error("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, cardApp.NewErrorResponse(
				"INVALID_AUTH_HEADER",
				"Invalid authorization header format",
				c.GetHeader("X-Request-ID"),
				nil,
			))
			c.Abort()
			return
		}

		token := parts[1]
		
		// Mock token validation - in production, this would validate JWT
		userID, role, err := validateToken(token)
		if err != nil {
			logger.Error("Token validation failed:", err)
			c.JSON(http.StatusUnauthorized, cardApp.NewErrorResponse(
				"INVALID_TOKEN",
				"Invalid or expired token",
				c.GetHeader("X-Request-ID"),
				nil,
			))
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", userID)
		c.Set("user_role", role)
		
		logger.Debug("User authenticated:", userID, "role:", role)
		c.Next()
	}
}

// validateToken validates a JWT token (mock implementation)
func validateToken(token string) (userID, role string, err error) {
	// Mock implementation - in production, this would:
	// 1. Parse and validate JWT signature
	// 2. Check token expiration
	// 3. Extract user claims
	// 4. Validate against user directory
	
	// Simple mock tokens for development
	mockTokens := map[string]struct {
		UserID string
		Role   string
	}{
		"emp001-token": {"emp001", "Employee"},
		"emp002-token": {"emp002", "Employee"},
		"emp003-token": {"emp003", "Employee"},
		"emp004-token": {"emp004", "HR"},
		"emp005-token": {"emp005", "Employee"},
		"admin-token":  {"admin", "HR"},
	}

	if user, exists := mockTokens[token]; exists {
		return user.UserID, user.Role, nil
	}

	return "", "", fmt.Errorf("invalid token")
}

// RequireRole middleware ensures user has required role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			logger.Error("User role not found in context")
			c.JSON(http.StatusUnauthorized, cardApp.NewErrorResponse(
				"MISSING_USER_ROLE",
				"User role not found",
				c.GetHeader("X-Request-ID"),
				nil,
			))
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			logger.Error("Invalid user role type in context")
			c.JSON(http.StatusInternalServerError, cardApp.NewErrorResponse(
				"INVALID_USER_ROLE",
				"Invalid user role",
				c.GetHeader("X-Request-ID"),
				nil,
			))
			c.Abort()
			return
		}

		// Check if user has required role
		if !hasRequiredRole(role, requiredRole) {
			logger.Error("Insufficient permissions, required:", requiredRole, "has:", role)
			c.JSON(http.StatusForbidden, cardApp.NewErrorResponse(
				"INSUFFICIENT_PERMISSIONS",
				"Insufficient permissions for this operation",
				c.GetHeader("X-Request-ID"),
				nil,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasRequiredRole checks if user role satisfies the requirement
func hasRequiredRole(userRole, requiredRole string) bool {
	// Role hierarchy: HR > Employee
	switch requiredRole {
	case "Employee":
		return userRole == "Employee" || userRole == "HR"
	case "HR":
		return userRole == "HR"
	default:
		return false
	}
}
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens and extracts user information
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Authorization header format must be Bearer {token}",
				},
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired token",
					"details": err.Error(),
				},
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Token is not valid",
				},
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_CLAIMS",
					"message": "Failed to parse token claims",
				},
			})
			c.Abort()
			return
		}

		// Extract user ID (adjust the claim name based on your JWT structure)
		var userID string
		if sub, ok := claims["sub"].(string); ok {
			userID = sub
		} else if uid, ok := claims["user_id"].(string); ok {
			userID = uid
		} else if uid, ok := claims["userId"].(string); ok {
			userID = uid
		}

		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "MISSING_USER_ID",
					"message": "User ID not found in token",
				},
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("userId", userID)
		c.Set("claims", claims)

		// Extract user email if available
		if email, ok := claims["email"].(string); ok {
			c.Set("userEmail", email)
		}

		// Extract user roles if available
		if roles, ok := claims["roles"].([]interface{}); ok {
			c.Set("userRoles", roles)
		}

		c.Next()
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("userId")
	if !exists {
		return "", fmt.Errorf("user ID not found in context")
	}

	str, ok := userID.(string)
	if !ok {
		return "", fmt.Errorf("user ID is not a string")
	}

	return str, nil
}

// GetUserEmail retrieves the user email from the context
func GetUserEmail(c *gin.Context) string {
	email, exists := c.Get("userEmail")
	if !exists {
		return ""
	}

	str, _ := email.(string)
	return str
}

// HasRole checks if the user has a specific role
func HasRole(c *gin.Context, roleName string) bool {
	roles, exists := c.Get("userRoles")
	if !exists {
		return false
	}

	roleList, ok := roles.([]interface{})
	if !ok {
		return false
	}

	for _, role := range roleList {
		if roleStr, ok := role.(string); ok && roleStr == roleName {
			return true
		}
	}

	return false
}

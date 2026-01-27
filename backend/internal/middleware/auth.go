package middleware

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"thank-you-card-backend/internal/models"
)

// Context keys
const (
	ContextCurrentEmployee = "currentEmployee"
)

// MockUserMiddleware attaches a mock current employee to the request context based on env MOCK_USER_ID.
// In MVP we assume the employee row already exists.
func MockUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := os.Getenv("MOCK_USER_ID")
		if idStr == "" {
			// default to 1 for local dev
			idStr = "1"
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil || id == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid mock user id"})
			return
		}

		// We don't hit DB here to keep middleware simple; assume caller seeded this user.
		// Handlers/services that need full employee details can look up from DB if needed.
		emp := &models.Employee{
			ID: uint(id),
		}
		c.Set(ContextCurrentEmployee, emp)
		c.Next()
	}
}

// RequireHRAdmin checks the attached employee's IsHRAdmin flag.
// It loads the full employee record from DB to verify IsHRAdmin status.
func RequireHRAdmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(ContextCurrentEmployee)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		emp, ok := val.(*models.Employee)
		if !ok || emp == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Load full employee record to check IsHRAdmin
		var fullEmployee models.Employee
		if err := db.First(&fullEmployee, emp.ID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "employee not found"})
			return
		}

		if !fullEmployee.IsHRAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: hr admin only"})
			return
		}

		// Update context with full employee record
		c.Set(ContextCurrentEmployee, &fullEmployee)
		c.Next()
	}
}

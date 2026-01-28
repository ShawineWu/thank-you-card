package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"thank-you-card-backend/internal/models"
	"thank-you-card-backend/internal/repositories"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
}

type LoginResult struct {
	Token    string
	User     *models.User
	Employee *models.Employee
}

type TokenClaims struct {
	UserID   uint
	Username string
	Role     string
	jwt.RegisteredClaims
}

type authService struct {
	userRepo     repositories.UserRepository
	employeeRepo repositories.EmployeeRepository
	jwtSecret    []byte
}

func NewAuthService(userRepo repositories.UserRepository, employeeRepo repositories.EmployeeRepository) AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "thank-you-card-secret-key-change-in-production" // Default for dev
	}
	return &authService{
		userRepo:     userRepo,
		employeeRepo: employeeRepo,
		jwtSecret:    []byte(secret),
	}
}

func (s *authService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	claims := &TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Load employee if linked
	var employee *models.Employee
	if user.EmployeeID != nil {
		emp, err := s.employeeRepo.GetByID(ctx, *user.EmployeeID)
		if err == nil {
			employee = emp
		}
	}

	return &LoginResult{
		Token:    tokenString,
		User:     user,
		Employee: employee,
	}, nil
}

func (s *authService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

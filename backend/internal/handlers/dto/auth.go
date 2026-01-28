package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID       uint             `json:"id"`
	Username string           `json:"username"`
	Role     string           `json:"role"`
	Employee *EmployeeSummary `json:"employee,omitempty"`
}

type EmployeeSummary struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
}

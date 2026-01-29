package models

// Employee represents an employee from the organization directory
type Employee struct {
	ID         string
	AADID      string // Azure AD Object ID, required for @mentions
	Name       string
	Email      string
	Department string
}

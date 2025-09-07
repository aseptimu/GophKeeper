package models

// User represents a user in the system.
type User struct {
	ID             string // Unique identifier for the user
	Login          string // User's login/username
	HashedPassword string // Hashed password for authentication
}

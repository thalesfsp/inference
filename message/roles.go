package message

// Role of the message's author.
type Role = string

const (
	// User role.
	User Role = "user"

	// System role.
	System Role = "system"
)

package domain

// UserRepository defines the contract for user persistence operations.
type UserRepository interface {
	SaveUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	DeleteUser(id string) error
	UpdateUser(user *User) error
	ListUsers(filter UserFilter, page int, limit int) ([]*User, int, error)
}

// UserFilter represents filter criteria for listing users.
type UserFilter struct {
	Name  string
	Email string
}

// PostRepository defines the contract for post persistence operations.
type PostRepository interface {
	SavePost(post *Post) error
}

// PasswordRecoveryRepository defines the contract for password recovery token persistence.
type PasswordRecoveryRepository interface {
	SavePasswordRecovery(recovery *PasswordRecovery) error
	GetPasswordRecoveryByToken(token string) (*PasswordRecovery, error)
	UpdatePasswordRecovery(recovery *PasswordRecovery) error
}

// EmailService defines the contract for email operations.
type EmailService interface {
	SendPasswordRecoveryEmail(email string, token string) error
}

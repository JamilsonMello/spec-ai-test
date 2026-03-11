package domain

type UserRepository interface {
	SaveUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	DeleteUser(id string) error
	UpdateUser(user *User) error
	ListUsers(filter UserFilter, page int, limit int) ([]*User, int, error)
}

type UserFilter struct {
	Name  string
	Email string
}

type PostRepository interface {
	SavePost(post *Post) error
}

type PasswordRecoveryRepository interface {
	SavePasswordRecovery(recovery *PasswordRecovery) error
	GetPasswordRecoveryByToken(token string) (*PasswordRecovery, error)
	UpdatePasswordRecovery(recovery *PasswordRecovery) error
}

type EmailService interface {
	SendPasswordRecoveryEmail(email string, token string) error
}

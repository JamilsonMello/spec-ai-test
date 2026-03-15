package domain

type UserRepository interface {
	SaveUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	FindUserByUuid(id string) (*User, error)
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
	GetPostByID(id string) (*Post, error)
	UpdatePost(post *Post) error
}

type PasswordRecoveryRepository interface {
	SavePasswordRecovery(recovery *PasswordRecovery) error
	GetPasswordRecoveryByToken(token string) (*PasswordRecovery, error)
	UpdatePasswordRecovery(recovery *PasswordRecovery) error
}

type EmailSender interface {
	SendPasswordRecoveryEmail(email string, token string) error
}

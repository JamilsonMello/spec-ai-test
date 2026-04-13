package domain

import "io"

type UserRepository interface {
	SaveUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	FindUserByUuid(id string) (*User, error)
	DeleteUser(id string) error
	UpdateUser(user *User) error
	UpdateProfilePictureURL(userID string, url string) error
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

type CommentRepository interface {
	SaveComment(comment *Comment) error
	GetCommentsByPostID(postID string) ([]*Comment, error)
	GetCommentByID(id string) (*Comment, error)
}

type ReactionRepository interface {
	GetReactionByCommentAndUser(commentID, userID, reactionType string) (*Reaction, error)
	SaveReaction(reaction *Reaction) error
	DeleteReaction(id string) error
}

type CommunityRepository interface {
	SaveCommunity(community *Community) error
	ExistsByName(name string) (bool, error)
	ExistsByID(id string) (bool, error)
}

type PasswordRecoveryRepository interface {
	SavePasswordRecovery(recovery *PasswordRecovery) error
	GetPasswordRecoveryByToken(token string) (*PasswordRecovery, error)
	UpdatePasswordRecovery(recovery *PasswordRecovery) error
}

type FileStorage interface {
	Save(file io.Reader, filename string) (string, error)
}

type EmailSender interface {
	SendPasswordRecoveryEmail(email string, token string) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword string, password string) error
}

type TransactionManager interface {
	ExecuteInTransaction(fn func(tx interface{}) error) error
}

type UserRepositoryTx interface {
	UpdateUserTx(tx interface{}, user *User) error
}

type PasswordRecoveryRepositoryTx interface {
	UpdatePasswordRecoveryTx(tx interface{}, recovery *PasswordRecovery) error
}

type ProductRepository interface {
	SaveProduct(product *Product) error
	GetProductByID(id string) (*Product, error)
	UpdateProduct(product *Product) error
}

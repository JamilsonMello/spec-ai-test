package domain

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

type PasswordRecoveryRepository interface {
	SavePasswordRecovery(recovery *PasswordRecovery) error
	GetPasswordRecoveryByToken(token string) (*PasswordRecovery, error)
	UpdatePasswordRecovery(recovery *PasswordRecovery) error
}

type EmailSender interface {
	SendPasswordRecoveryEmail(email string, token string) error
}

package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

const insertPostSQL = `
INSERT INTO posts (id, content, author_id, community_id, created_at)
VALUES ($1, $2, $3, $4, $5)
`

const updatePostContentSQL = `
UPDATE posts
SET content = $1, updated_at = $2
WHERE id = $3
`

const selectPostByIDSQL = `
SELECT id, content, author_id, created_at, updated_at
FROM posts
WHERE id = $1
`

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) SavePost(post *domain.Post) error {
	var communityID interface{}
	if post.CommunityID != "" {
		communityID = post.CommunityID
	}
	_, err := r.db.Exec(insertPostSQL, post.ID, post.Content, post.AuthorID, communityID, post.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save post: %w", err)
	}
	return nil
}

func (r *PostRepository) UpdatePost(post *domain.Post) error {
	result, err := r.db.Exec(updatePostContentSQL, post.Content, post.UpdatedAt, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrPostNotFound
	}
	return nil
}

func (r *PostRepository) GetPostByID(id string) (*domain.Post, error) {
	post := &domain.Post{}
	err := r.db.QueryRow(selectPostByIDSQL, id).Scan(
		&post.ID, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrPostNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post by id: %w", err)
	}
	return post, nil
}

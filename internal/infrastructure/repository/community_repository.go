package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type CommunityRepository struct {
	db *sql.DB
}

func NewCommunityRepository(db *sql.DB) *CommunityRepository {
	return &CommunityRepository{db: db}
}

func (r *CommunityRepository) SaveCommunity(community *domain.Community) error {
	query := `INSERT INTO communities (id, name, description, owner_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, community.ID, community.Name, community.Description, community.OwnerID, community.CreatedAt, community.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save community: %w", err)
	}
	return nil
}

func (r *CommunityRepository) ExistsByID(id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM communities WHERE id = $1)`
	var exists bool
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check community existence by id: %w", err)
	}
	return exists, nil
}

func (r *CommunityRepository) ExistsByName(name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM communities WHERE name = $1)`
	var exists bool
	err := r.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check community existence: %w", err)
	}
	return exists, nil
}

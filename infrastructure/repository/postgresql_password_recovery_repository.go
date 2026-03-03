package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/cadastro-de-usuarios/domain"
)

type PostgreSQLPasswordRecoveryRepository struct {
	db *sql.DB
}

func NewPostgreSQLPasswordRecoveryRepository(db *sql.DB) *PostgreSQLPasswordRecoveryRepository {
	return &PostgreSQLPasswordRecoveryRepository{db: db}
}

func (r *PostgreSQLPasswordRecoveryRepository) SavePasswordRecovery(recovery *domain.PasswordRecovery) error {
	query := `INSERT INTO password_recoveries (id, token, user_id, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, recovery.ID, recovery.Token, recovery.UserID, recovery.ExpiresAt, recovery.Used, recovery.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save password recovery: %w", err)
	}
	return nil
}

func (r *PostgreSQLPasswordRecoveryRepository) GetPasswordRecoveryByToken(token string) (*domain.PasswordRecovery, error) {
	query := `SELECT id, token, user_id, expires_at, used, created_at
		FROM password_recoveries WHERE token = $1`
	recovery := &domain.PasswordRecovery{}
	err := r.db.QueryRow(query, token).Scan(
		&recovery.ID, &recovery.Token, &recovery.UserID,
		&recovery.ExpiresAt, &recovery.Used, &recovery.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrRecoveryTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get password recovery by token: %w", err)
	}
	return recovery, nil
}

func (r *PostgreSQLPasswordRecoveryRepository) UpdatePasswordRecovery(recovery *domain.PasswordRecovery) error {
	query := `UPDATE password_recoveries SET used=$1 WHERE token=$2`
	result, err := r.db.Exec(query, recovery.Used, recovery.Token)
	if err != nil {
		return fmt.Errorf("failed to update password recovery: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrRecoveryTokenNotFound
	}
	return nil
}

func (r *PostgreSQLPasswordRecoveryRepository) InvalidateAllUserTokens(userID string) error {
	query := `UPDATE password_recoveries SET used = true WHERE user_id = $1 AND used = false`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to invalidate user tokens: %w", err)
	}
	return nil
}

func (r *PostgreSQLPasswordRecoveryRepository) GetLatestPasswordRecoveryByUserID(userID string) (*domain.PasswordRecovery, error) {
	query := `SELECT id, token, user_id, expires_at, used, created_at
		FROM password_recoveries WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`
	recovery := &domain.PasswordRecovery{}
	err := r.db.QueryRow(query, userID).Scan(
		&recovery.ID, &recovery.Token, &recovery.UserID,
		&recovery.ExpiresAt, &recovery.Used, &recovery.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrRecoveryTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest password recovery: %w", err)
	}
	return recovery, nil
}

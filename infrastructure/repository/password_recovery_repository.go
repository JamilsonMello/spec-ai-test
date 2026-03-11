package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/cadastro-de-usuarios/domain"
)

type PasswordRecoveryRepository struct {
	db *sql.DB
}

func NewPasswordRecoveryRepository(db *sql.DB) *PasswordRecoveryRepository {
	return &PasswordRecoveryRepository{db: db}
}

func (passwordRecoveryRepository *PasswordRecoveryRepository) SavePasswordRecovery(recovery *domain.PasswordRecovery) error {
	query := `INSERT INTO password_recoveries (id, token, user_id, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := passwordRecoveryRepository.db.Exec(query, recovery.ID, recovery.Token, recovery.UserID, recovery.ExpiresAt, recovery.Used, recovery.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save password recovery: %w", err)
	}
	return nil
}

func (passwordRecoveryRepository *PasswordRecoveryRepository) GetPasswordRecoveryByToken(token string) (*domain.PasswordRecovery, error) {
	query := `SELECT id, token, user_id, expires_at, used, created_at
		FROM password_recoveries WHERE token = $1`
	recovery := &domain.PasswordRecovery{}
	err := passwordRecoveryRepository.db.QueryRow(query, token).Scan(
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

func (passwordRecoveryRepository *PasswordRecoveryRepository) UpdatePasswordRecovery(recovery *domain.PasswordRecovery) error {
	query := `UPDATE password_recoveries SET used=$1 WHERE token=$2`
	result, dbErr := passwordRecoveryRepository.db.Exec(query, recovery.Used, recovery.Token)
	if dbErr != nil {
		return fmt.Errorf("failed to update password recovery: %w", dbErr)
	}
	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		return fmt.Errorf("failed to check rows affected: %w", rowsErr)
	}
	if affected == 0 {
		return domain.ErrRecoveryTokenNotFound
	}
	return nil
}

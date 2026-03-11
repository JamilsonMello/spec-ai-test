package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/example/cadastro-de-usuarios/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (userRepository *UserRepository) SaveUser(user *domain.User) error {
	query := `INSERT INTO users (id, name, surname, email, birth_date, password, recovery_token, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, dbErr := userRepository.db.Exec(query, user.ID, user.Name, user.Surname, user.Email, user.BirthDate, user.Password, user.RecoveryToken, user.Role, user.CreatedAt)
	if dbErr != nil {
		if pqErr, ok := dbErr.(*pq.Error); ok && pqErr.Code == "23505" {
			return domain.ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to save user: %w", dbErr)
	}
	return nil
}

func (userRepository *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, name, surname, email, birth_date, password, recovery_token, role, created_at
		FROM users WHERE email = $1`
	user := &domain.User{}
	err := userRepository.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Surname, &user.Email,
		&user.BirthDate, &user.Password, &user.RecoveryToken, &user.Role, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (userRepository *UserRepository) GetUserByID(id string) (*domain.User, error) {
	query := `SELECT id, name, surname, email, birth_date, password, recovery_token, role, created_at
		FROM users WHERE id = $1`
	user := &domain.User{}
	err := userRepository.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Surname, &user.Email,
		&user.BirthDate, &user.Password, &user.RecoveryToken, &user.Role, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (userRepository *UserRepository) UpdateUser(user *domain.User) error {
	query := `UPDATE users SET name=$1, surname=$2, email=$3, birth_date=$4, password=$5, recovery_token=$6, role=$7
		WHERE id=$8`
	result, err := userRepository.db.Exec(query, user.Name, user.Surname, user.Email, user.BirthDate, user.Password, user.RecoveryToken, user.Role, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (userRepository *UserRepository) DeleteUser(id string) error {
	result, err := userRepository.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (userRepository *UserRepository) ListUsers(filter domain.UserFilter, page int, limit int) ([]*domain.User, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1

	if filter.Name != "" {
		where += fmt.Sprintf(" AND LOWER(name) LIKE $%d", idx)
		args = append(args, "%"+strings.ToLower(filter.Name)+"%")
		idx++
	}
	if filter.Email != "" {
		where += fmt.Sprintf(" AND LOWER(email) LIKE $%d", idx)
		args = append(args, "%"+strings.ToLower(filter.Email)+"%")
		idx++
	}

	var total int
	if err := userRepository.db.QueryRow("SELECT COUNT(*) FROM users "+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	listArgs := append(args, limit, (page-1)*limit)
	listQuery := fmt.Sprintf(
		`SELECT id, name, surname, email, birth_date, password, recovery_token, role, created_at
		FROM users %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, idx, idx+1,
	)

	rows, err := userRepository.db.Query(listQuery, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(
			&user.ID, &user.Name, &user.Surname, &user.Email,
			&user.BirthDate, &user.Password, &user.RecoveryToken, &user.Role, &user.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate users: %w", err)
	}
	return users, total, nil
}

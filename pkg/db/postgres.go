package db

import (
	"database/sql"
	"fmt"

	"github.com/example/cadastro-de-usuarios/pkg/db/queries"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	schemas := []string{
		queries.CreateUsersTableSQL,
		queries.CreatePostsTableSQL,
		queries.CreatePasswordRecoveriesTableSQL,
	}

	if err := InitSchema(db, schemas); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

func InitSchema(db *sql.DB, schemaSQLs []string) error {
	for _, schemaSQL := range schemaSQLs {
		if _, err := db.Exec(schemaSQL); err != nil {
			return fmt.Errorf("failed to execute schema SQL: %w", err)
		}
	}
	return nil
}

package infrastructure

import "database/sql"

type SQLTransactionManager struct {
	db *sql.DB
}

func NewSQLTransactionManager(db *sql.DB) *SQLTransactionManager {
	return &SQLTransactionManager{db: db}
}

func (tm *SQLTransactionManager) ExecuteInTransaction(fn func(tx interface{}) error) error {
	sqlTx, err := tm.db.Begin()
	if err != nil {
		return err
	}

	if err := fn(sqlTx); err != nil {
		_ = sqlTx.Rollback()
		return err
	}

	return sqlTx.Commit()
}

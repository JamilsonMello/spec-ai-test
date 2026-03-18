package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var migrationFileRegex = regexp.MustCompile(`^\d{6}_[a-z0-9_]+\.(up|down)\.sql$`)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := RunMigrations(db, "migrations"); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB, migrationsDir string) error {
	upFiles, err := LoadMigrationFiles(migrationsDir)
	if err != nil {
		return err
	}

	for _, upFile := range upFiles {
		content, err := os.ReadFile(filepath.Join(migrationsDir, upFile))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", upFile, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			downFile := strings.Replace(upFile, ".up.sql", ".down.sql", 1)
			downContent, readErr := os.ReadFile(filepath.Join(migrationsDir, downFile))
			if readErr == nil {
				db.Exec(string(downContent))
			}
			return fmt.Errorf("failed to execute migration %s: %w", upFile, err)
		}
	}

	return nil
}

func LoadMigrationFiles(migrationsDir string) ([]string, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	upFiles := make([]string, 0)
	downFiles := make(map[string]bool)
	sequenceNumbers := make(map[int]string)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !migrationFileRegex.MatchString(name) {
			continue
		}

		if strings.HasSuffix(name, ".down.sql") {
			prefix := name[:6]
			downFiles[prefix] = true
		}

		if strings.HasSuffix(name, ".up.sql") {
			prefix := name[:6]
			seq, _ := strconv.Atoi(prefix)

			if existing, found := sequenceNumbers[seq]; found {
				return nil, fmt.Errorf("duplicate migration sequence number detected: %s and %s", existing, name)
			}
			sequenceNumbers[seq] = name
			upFiles = append(upFiles, name)
		}
	}

	sort.Strings(upFiles)

	for _, upFile := range upFiles {
		prefix := upFile[:6]
		if !downFiles[prefix] {
			seq, _ := strconv.Atoi(prefix)
			return nil, fmt.Errorf("migration pair missing for index [%d]", seq)
		}
	}

	if len(upFiles) > 0 {
		for i := 0; i < len(upFiles); i++ {
			expectedSeq := i + 1
			prefix := upFiles[i][:6]
			actualSeq, _ := strconv.Atoi(prefix)
			if actualSeq != expectedSeq {
				return nil, fmt.Errorf("migration sequence gap detected: expected %06d but found %06d", expectedSeq, actualSeq)
			}
		}
	}

	return upFiles, nil
}

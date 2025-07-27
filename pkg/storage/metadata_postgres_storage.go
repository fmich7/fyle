package storage

import (
	"database/sql"
	"fmt"

	"github.com/fmich7/fyle/pkg/config"
)

type PQMetadataStorage struct {
	Storage
	db      *sql.DB
	connStr string
}

func NewPQMetadataStorage(c config.PostgresCredentials) (*PQMetadataStorage, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DB_USER,
		c.DB_PASSWORD,
		c.DB_HOST,
		c.DB_PORT,
		c.DB_NAME,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PQMetadataStorage{
		db:      db,
		connStr: connStr,
	}, nil
}

func (s *PQMetadataStorage) Shutdown() error {
	return s.db.Close()
}

func (s *PQMetadataStorage) GetFileMetadata() error {
	// Implementation for getting file metadata
	return nil
}

func (s *PQMetadataStorage) CreateFileMetadata() error {
	// Implementation for creating file metadata
	return nil
}

func (s *PQMetadataStorage) RunMigrations(migrationsDir string) error {
	return RunMigrations(s.db, migrationsDir, []string{"postgres", "metadata"})
}

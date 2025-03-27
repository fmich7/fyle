package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/fmich7/fyle/pkg/config"

	_ "github.com/lib/pq"
)

// PQUserStorage represents a user storage.
type PQUserStorage struct {
	db      *sql.DB
	connStr string
}

// NewPQUserStorage creates a new user storage.
func NewPQUserStorage(c config.PostgresCredentials) (*PQUserStorage, error) {
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
		return nil, fmt.Errorf("opening database %v", err)
	}

	log.Println("Connected to PostgreSQL!")

	return &PQUserStorage{
		db:      db,
		connStr: connStr,
	}, nil
}

// CloseDatabase closes the database connection.
func (d *PQUserStorage) CloseDatabase() error {
	return d.db.Close()
}

// StoreUser stores a user in the database.
func (d *PQUserStorage) StoreUser(user *auth.User) error {
	existQuery := `
	SELECT EXISTS (
		SELECT 1 FROM users WHERE username = $1
	)
	`
	var exists bool
	err := d.db.QueryRow(existQuery, user.Username).Scan(&exists)

	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	if exists {
		return errors.New("user already exists")
	}

	// insert new user into db
	insertQuery := `
	INSERT INTO users (username, password, salt)
	VALUES ($1, $2, $3)
	`
	_, err = d.db.Exec(insertQuery, user.Username, user.Password, user.Salt)
	if err != nil {
		return fmt.Errorf("inserting user: %v", err)
	}

	return nil
}

// RetrieveUser retrieves a user from the database.
func (d *PQUserStorage) RetrieveUser(username string) (*auth.User, error) {
	user := new(auth.User)
	getUserQuery := `
	SELECT id, username, password, salt
	FROM users
	WHERE username = $1
	`
	err := d.db.QueryRow(getUserQuery, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Salt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("retrieving user: %v", err)
	}

	return user, nil
}

// RunMigrations runs all migrations in the given directory.
func (d *PQUserStorage) RunMigrations(migrationsDir string) error {
	log.Println("Starting migrations...")

	// Read migrations files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	// Get migration filenames
	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		filePath := filepath.Join(migrationsDir, file)

		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		sqlStatement := string(sqlBytes)

		log.Println("Applying migration:", file)
		_, err = d.db.Exec(sqlStatement)
		if err != nil {
			return err
		}
	}

	return nil
}

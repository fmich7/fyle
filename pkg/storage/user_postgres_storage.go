package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/fmich7/fyle/pkg/config"

	_ "github.com/lib/pq"
)

// PQUserStorage represents a user storage.
type PQUserStorage struct {
	Storage
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

	return &PQUserStorage{
		db:      db,
		connStr: connStr,
	}, nil
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
	return RunMigrations(d.db, migrationsDir, []string{"postgres", "users"})
}

func (d *PQUserStorage) Shutdown() error {
	return d.db.Close()
}

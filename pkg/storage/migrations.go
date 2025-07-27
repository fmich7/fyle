package storage

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// function that runs migrations on the database.
func RunMigrations(db *sql.DB, migrationsDir string, targets []string) error {
	log.Println("Starting migrations...", targets)

	// Read migrations files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	// Get migration filenames
	var migrationFiles []string
	for _, file := range files {
		name := file.Name()

		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		// Check if the file is in the target list
		for _, target := range targets {
			if strings.Contains(name, target) {
				migrationFiles = append(migrationFiles, name)
				break
			}
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
		_, err = db.Exec(sqlStatement)
		if err != nil {
			return err
		}
	}

	return nil

}

// TODO: Implement rollback mechanism

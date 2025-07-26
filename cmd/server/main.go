package main

import (
	"flag"
	"log"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/server"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := new(config.Config)
	cfg.LoadConfig(*configPath)

	// Setup user storage
	userStorage, err := storage.NewPQUserStorage(cfg.PostgresCredentials)
	if err != nil {
		log.Fatalf("error connecting to user db: %v\n", err)
	}
	defer userStorage.CloseDatabase()
	if err := userStorage.RunMigrations(cfg.MigrationPath); err != nil {
		log.Fatalf("error running migrations: %v\n", err)
	}

	// Setup file storage
	fileStorage, err := storage.NewDiskFileStorage(cfg.UploadsLocation, afero.NewOsFs())
	if err != nil {
		log.Fatalf("error creating disk storage: %v\n", err)
	}

	serverStorage := storage.NewServerStorage(fileStorage, userStorage)
	server := server.NewServer(cfg, serverStorage)

	log.Println("Server listening on port", cfg.ServerPort)
	log.Fatal(server.Start())
}

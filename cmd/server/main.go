package main

import (
	"flag"
	"log"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/server"
	"github.com/fmich7/fyle/pkg/storage"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := new(config.Config)
	cfg.LoadConfig(*configPath)

	storages := &storage.ServerStorage{
		FileStorage:     nil,
		UserStorage:     nil,
		MetadataStorage: nil,
	}
	defer storages.Shutdown()
	if err := storages.SetUpStorages(cfg); err != nil {
		log.Fatalf("Failed to set up storages: %v", err)
	}

	if errs := storages.RunMigrations(cfg.MigrationPath); len(errs) > 0 {
		log.Fatalf("Failed to run migrations: %v", errs)
	}

	server := server.NewServer(cfg, storages)
	defer server.Shutdown()

	log.Println("Server listening on port", cfg.ServerPort)
	log.Fatal(server.Start())
}

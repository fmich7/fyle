package main

import (
	"fmt"
	"log"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/server"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
)

func main() {
	cfg := new(config.Config)
	cfg.LoadConfig("example_server.env")

	fileStorage, err := storage.NewDiskFileStorage(cfg.UploadsLocation, afero.NewOsFs())
	if err != nil {
		log.Fatalf("error creating disk storage: %v\n", err)
	}
	serverStorage := storage.NewServerStorage(fileStorage, nil)

	server := server.NewServer(cfg, serverStorage)

	fmt.Println("Server listening on port", cfg.ServerPort)
	log.Fatal(server.Start())
}

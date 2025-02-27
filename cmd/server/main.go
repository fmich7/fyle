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

	store, err := storage.NewDiskStorage(cfg.UploadsLocation, afero.NewOsFs())
	if err != nil {
		log.Fatalf("error creating disk storage: %v\n", err)
	}

	server := server.NewServer(cfg, store)

	fmt.Println("Server listening on port", cfg.ServerPort)
	log.Fatal(server.Start())
}

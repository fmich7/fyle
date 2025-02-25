package main

import (
	"fmt"
	"log"

	"github.com/fmich7/fyle/pkg/server"
	"github.com/fmich7/fyle/pkg/storage"
)

func main() {
	PORT := ":3000"
	store, err := storage.NewDiskStorage("uploads", nil)
	if err != nil {
		log.Fatalf("error creating disk storage: %v\n", err)
	}

	server := server.NewServer(PORT, store)

	fmt.Println("Server listening on port", PORT)
	log.Fatal(server.Start())
}

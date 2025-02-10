package main

import (
	"fmt"
	"log"

	"github.com/fmich7/fyle/pkg/api"
	"github.com/fmich7/fyle/pkg/storage"
)

func main() {
	PORT := ":3000"
	store := storage.NewDiskStorage("../../uploads")

	server := api.NewServer(PORT, store)

	fmt.Println("Server listening on port", PORT)
	log.Fatal(server.Start())
}

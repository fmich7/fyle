package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: fyle <command> [options]")
		return
	}
	command := os.Args[1]
	switch command {
	case "upload":
		fmt.Println("Uploading file...")
	case "download":
		fmt.Println("Downloading file...")
	default:
		fmt.Println("Unknown command:", command)
	}
}

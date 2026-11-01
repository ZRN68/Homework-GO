package main

import (
	"os"

	"homework6/internal/cli"
	"homework6/internal/storage"
)

func main() {
	store := storage.NewFile("inventory.txt")
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, store))
}

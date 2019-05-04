package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyaglow/dor"
	"github.com/ilyaglow/dor/web"
)

func main() {
	var port string
	if port = os.Getenv("DOR_PORT"); port == "" {
		port = "8080"
	}
	port = fmt.Sprintf(":%s", port)

	storageURL := os.Getenv("DOR_STORAGE_URL")

	d, err := dor.New(os.Getenv("DOR_STORAGE"), storageURL, false)
	if err != nil {
		panic(err)
	}

	log.Fatal(dorweb.Serve(port, d))
}

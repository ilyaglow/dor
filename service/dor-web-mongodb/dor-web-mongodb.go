package main

import (
	"fmt"
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

	mongoURL := os.Getenv("DOR_MONGO_URL")

	d, err := dor.New("mongodb", mongoURL, false)
	if err != nil {
		panic(err)
	}

	dorweb.Serve(port, d)
}

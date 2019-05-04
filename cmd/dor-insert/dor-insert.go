package main

import (
	"os"
	"time"

	"github.com/ilyaglow/dor"
)

const updatePeriod = time.Hour * 24

func main() {
	storageURL := os.Getenv("DOR_STORAGE_URL")

	d, err := dor.New(os.Getenv("DOR_STORAGE"), storageURL, false)
	if err != nil {
		panic(err)
	}

	if err := d.FillByTimer(updatePeriod); err != nil {
		panic(err)
	}
}

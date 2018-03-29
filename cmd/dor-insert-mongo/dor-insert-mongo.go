package main

import (
	"os"
	"time"

	"github.com/ilyaglow/dor"
)

const updatePeriod = time.Hour * 24

func main() {
	mongoURL := os.Getenv("DOR_MONGO_URL")

	d, err := dor.New("mongodb", mongoURL, false)
	if err != nil {
		panic(err)
	}

	if err := d.FillByTimer(updatePeriod); err != nil {
		panic(err)
	}
}

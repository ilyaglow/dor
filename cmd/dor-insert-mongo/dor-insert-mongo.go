package main

import (
	"flag"
	"time"

	"github.com/ilyaglow/dor"
)

const updatePeriod = time.Hour * 24

func main() {
	mongoURL := flag.String("mongo", "", "MongoDB URL")
	flag.Parse()

	d, err := dor.New("mongodb", *mongoURL, false)
	if err != nil {
		panic(err)
	}

	if err := d.FillByTimer(updatePeriod); err != nil {
		panic(err)
	}
}

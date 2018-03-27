package main

import (
	"time"

	"github.com/ilyaglow/dor"
)

const updatePeriod = time.Hour * 24

func main() {
	d, err := dor.New("mongodb", "", false)
	if err != nil {
		panic(err)
	}

	if err := d.FillByTimer(updatePeriod); err != nil {
		panic(err)
	}
}

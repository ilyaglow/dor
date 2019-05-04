package main

import (
	"flag"
	"os"
	"time"

	"github.com/ilyaglow/dor"
	"github.com/peterbourgon/ff"
)

const updatePeriod = time.Hour * 24

func main() {
	fs := flag.NewFlagSet("DOR", flag.ExitOnError)
	var (
		storage  = fs.String("storage", "clickhouse", "storage type")
		location = fs.String("storage-url", "", "url of the storage")
	)
	ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("DOR"))

	d, err := dor.New(*storage, *location, false)
	if err != nil {
		panic(err)
	}

	if err := d.FillByTimer(updatePeriod); err != nil {
		panic(err)
	}
}

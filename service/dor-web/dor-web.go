package main

import (
	"github.com/ilyaglow/dor"
	"github.com/ilyaglow/dor/web"
)

func main() {
	d, err := dor.New("mongodb", "", false)
	if err != nil {
		panic(err)
	}

	d.Ingesters = []dor.Ingester{
		dor.NewAlexa(),
		dor.NewMajestic(),
		dor.NewStatvoo(),
		dor.NewUmbrella(),
		dor.NewPageRank(),
	}

	dorweb.Serve(":9999", d)
}

package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ilyaglow/dor"
	"github.com/ilyaglow/dor/web"
	"github.com/kataras/iris"
)

func main() {
	bindAddr := flag.String("host", "127.0.0.1", "IP-address to bind")
	bindPort := flag.String("port", "8080", "Port to bind")
	flag.Parse()
	hp := fmt.Sprintf("%s:%s", *bindAddr, *bindPort)
	fmt.Println(hp)

	duration := time.Hour * 24

	d := &dor.DomainRank{}
	if err := d.FillAndUpdate(duration); err != nil {
		log.Fatal(err)
	}
	dorweb.Serve(hp, d)
}

func Serve(address string, d *dor.DomainRank) {
	app := iris.New()
	app.Get("/rank/{domain:string}", func(ctx iris.Context) {
		r := d.Find(ctx.Params().Get("domain"))
		ctx.JSON(r)
	})

	app.Run(iris.Addr(address), iris.WithCharset("UTF-8"))
}

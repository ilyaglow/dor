package dorweb

import (
	"log"

	"github.com/ilyaglow/dor"
	"github.com/kataras/iris"
)

// Serve represents a web interaction with the DomainRank
func Serve(address string, d *dor.App) {
	app := iris.New()
	app.Get("/rank/{domain:string}", func(ctx iris.Context) {
		r, err := d.Find(ctx.Params().Get("domain"))
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
		} else {
			ctx.JSON(r)
		}
	})

	log.Fatal(app.Run(iris.Addr(address), iris.WithCharset("UTF-8")))
}

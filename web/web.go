package dorweb

import (
	"github.com/ilyaglow/dor"
	"github.com/kataras/iris"
)

// Serve represents a web interaction with the DomainRank
func Serve(address string, d *dor.DomainRank) {
	app := iris.New()
	app.Get("/rank/{domain:string}", func(ctx iris.Context) {
		r := d.Find(ctx.Params().Get("domain"))
		ctx.JSON(r)
	})

	app.Run(iris.Addr(address), iris.WithCharset("UTF-8"))
}

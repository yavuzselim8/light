package main

import (
	"fmt"
	"github.com/yavuzselim8/light"
)

func main() {
	app := light.NewApp()

	router := light.NewRouter()

	router.Get("/route", func(ctx *light.Context) {
		ctx.SendJSON("{\"10\": 10}")
	})
	router.UseBefore(func(ctx *light.Context) {
		fmt.Println("I am little gangsta, you know!")
	})

	app.UseBefore(func(ctx *light.Context) {
		fmt.Println("Hell Yeah!")
	})
	app.Get("/route", func(ctx *light.Context) {
		ctx.SendJSON("{\"20\": 20}")
	})

	app.RegisterRouter("/some", router)

	app.Listen("localhost:8080")
}

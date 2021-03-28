package light

import "github.com/valyala/fasthttp"

type App struct {
	*Router
	middlewaresAlwaysAfter  []handler
	middlewaresAlwaysBefore []handler
}

func NewApp() *App{
	return &App{Router: NewRouter()}
}

func (app *App) RegisterRouter(baseRoute string, router *Router) {
	for method, routes := range router.routes {
		for route, handler := range routes {
			app.routes[method][baseRoute+route] = app.getMiddlewaredHandler(handler)
		}
	}
	router.routes = nil
}

func (app *App) UseBeforeAlways(middleware handler) {
	app.middlewaresAlwaysBefore = append(app.middlewaresAlwaysBefore, middleware)
}

func (app *App) UseAfterAlways(middleware handler) {
	app.middlewaresAlwaysAfter = append(app.middlewaresAlwaysAfter, middleware)
}

func (app *App) globalHandler(ctx *fasthttp.RequestCtx) {
	lc := &Context{ctx: ctx, Data: make(map[string]interface{})}
	for _, middleware := range app.middlewaresAlwaysBefore {
		middleware(lc)
	}
	handler, exist := app.routes[string(ctx.Method())][string(ctx.RequestURI())]
	if exist {
		handler(lc)
		return
	}
	for _, middleware := range app.middlewaresAlwaysAfter {
		middleware(lc)
	}
}

func (app *App) Listen(address string) {
	err := fasthttp.ListenAndServe(address, app.globalHandler)
	if err != nil {
		panic(err)
	}
}

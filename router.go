package light

import (
	"github.com/valyala/fasthttp"
)

type handlerFunction func(lc *Context)

type Light struct {
	middlewares []handlerFunction
	routes      map[string]map[string]handlerFunction
}

type Context struct {
	ctx *fasthttp.RequestCtx
}

func (lc *Context) SendJSON(jsonString string) {
	lc.ctx.SuccessString("application/json", jsonString)
}

func New() *Light {
	return &Light{
		routes: map[string]map[string]handlerFunction{"GET": {}, "POST": {}},
	}
}

func (light *Light) RegisterRouter(baseRoute string, router *Light) {
	for method, routes := range router.routes {
		for route, handler := range routes {
			light.routes[method][baseRoute+route] = light.getMiddlewaredHandler(handler)
		}
	}
	router.routes = nil
}

func (light *Light) Use(middleware handlerFunction) {
	light.middlewares = append(light.middlewares, middleware)
}

func (light *Light) Get(route string, handler handlerFunction) {
	newHandler := light.getMiddlewaredHandler(handler)
	light.routes["GET"][route] = newHandler
}

func (light *Light) Post(route string, handler handlerFunction) {
	newHandler := light.getMiddlewaredHandler(handler)
	light.routes["POST"][route] = newHandler
}

func (light *Light) getMiddlewaredHandler(handler handlerFunction) func(lc *Context) {
	return func(lc *Context) {
		for _, middleware := range light.middlewares {
			middleware(lc)
		}
		handler(lc)
	}
}

func (light *Light) globalHandler(ctx *fasthttp.RequestCtx) {
	if handler, exist := light.routes[string(ctx.Method())][string(ctx.RequestURI())]; exist {
		handler(&Context{ctx: ctx})
	}
}

func (light *Light) ListenAndServe(address string) {
	err := fasthttp.ListenAndServe(address, light.globalHandler)
	if err != nil {
		panic(err)
	}
}

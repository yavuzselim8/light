package light

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type handlerFunction func(lc *LightContext)

type Light struct {
	middlewares []handlerFunction
	routes      map[string]map[string]handlerFunction
}

type LightContext struct {
	ctx *fasthttp.RequestCtx
}

func (lc *LightContext) SendJSON(jsonString string) {
	lc.ctx.SuccessString("application/json", jsonString)
}

func New() *Light {
	light := &Light{
		routes: map[string]map[string]handlerFunction{"GET": {}, "POST": {}},
	}
	return light
}

func (light *Light) RegisterLight(baseRoute string, otherLight *Light) {
	for method, routes := range otherLight.routes {
		for route, handler := range routes {
			light.routes[method][baseRoute+route] = light.getMiddlewaredHandler(handler)
		}
	}
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

func (light *Light) getMiddlewaredHandler(handler handlerFunction) func(lc *LightContext) {
	return func(lc *LightContext) {
		for _, middleware := range light.middlewares {
			middleware(lc)
		}
		handler(lc)
	}
}

func (light *Light) handler(ctx *fasthttp.RequestCtx) {
	if handlerFunc, exist := light.routes[string(ctx.Method())][string(ctx.RequestURI())]; exist {
		handlerFunc(&LightContext{ctx: ctx})
	}
}

func (light *Light) ListenAndServe(address string) {
	err := fasthttp.ListenAndServe(address, light.handler)
	if err != nil {
		fmt.Println(err)
	}
}


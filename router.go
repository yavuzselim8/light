package light

type handler func(ctx *Context)

type Router struct {
	middlewaresBefore []handler
	routes            map[string]map[string]handler
	middlewaresAfter  []handler
}

func NewRouter() *Router {
	return &Router{
		routes: map[string]map[string]handler{"GET": {}, "POST": {}, "PUT": {}, "DELETE": {}},
	}
}

func (router *Router) UseBefore(middleware handler) {
	router.middlewaresBefore = append(router.middlewaresBefore, middleware)
}

func (router *Router) UseAfter(middleware handler) {
	router.middlewaresAfter = append(router.middlewaresAfter, middleware)
}

func (router *Router) Get(route string, handler handler) {
	newHandler := router.getMiddlewaredHandler(handler)
	router.routes["GET"][route] = newHandler
}

func (router *Router) Post(route string, handler handler) {
	newHandler := router.getMiddlewaredHandler(handler)
	router.routes["POST"][route] = newHandler
}

func (router *Router) Put(route string, handler handler) {
	newHandler := router.getMiddlewaredHandler(handler)
	router.routes["PUT"][route] = newHandler
}

func (router *Router) Delete(route string, handler handler) {
	newHandler := router.getMiddlewaredHandler(handler)
	router.routes["DELETE"][route] = newHandler
}

func (router *Router) getMiddlewaredHandler(handler handler) func(ctx *Context) {
	return func(ctx *Context) {
		for _, middleware := range router.middlewaresBefore {
			middleware(ctx)
		}
		handler(ctx)
		for _, middleware := range router.middlewaresAfter {
			middleware(ctx)
		}
	}
}
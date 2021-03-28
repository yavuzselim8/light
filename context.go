package light

import "github.com/valyala/fasthttp"

type Context struct {
	ctx  *fasthttp.RequestCtx
	Data map[string]interface{}
}

func (lc *Context) SendJSON(jsonString string) {
	lc.ctx.SuccessString("application/json", jsonString)
}

func (lc *Context) GetParam(key string) {
}
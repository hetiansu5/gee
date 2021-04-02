package gee

import (
	"net/http"
)

type Engine struct {
	*RouterGroup
	handlers map[string]HandlerFunc
	router   *Router
}

func New() *Engine {
	engine := &Engine{
		handlers: make(map[string]HandlerFunc),
		router:   newRouter(),
	}
	engine.RouterGroup = newRouterGroup(engine, nil)
	return engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.handle(c)
}

func (e *Engine) handle(c *Context) {
	node, params := e.router.getRouter(c.method, c.path)
	c.params = params
	if node == nil {
		c.handlers = append(c.handlers, func(ctx *Context) {
			ctx.Page404()
		})
	} else {
		c.handlers = append(c.handlers, e.router.getFuncHandlers(c.method, node.pattern)...)
	}
	c.Next()
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

package gee

import (
	"sync"
)

const (
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodDelete  = "DELETE"
	MethodPatch   = "PATCH"
	MethodHead    = "HEAD"
	MethodOptions = "OPTIONS"
)

var methods = []string{
	MethodGet, MethodPost, MethodPut, MethodDelete, MethodPatch, MethodHead, MethodOptions,
}

// IRoutes defines all router handle interface.
type IRoutes interface {
	Use(...HandlerFunc) IRoutes

	Handle(string, string, ...HandlerFunc) IRoutes
	Any(string, ...HandlerFunc) IRoutes
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
	PATCH(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	OPTIONS(string, ...HandlerFunc) IRoutes
	HEAD(string, ...HandlerFunc) IRoutes
}

type RouterGroup struct {
	engine         *Engine
	parent         *RouterGroup
	prefix         string
	middlewares    []HandlerFunc
	middlewareInit bool
	lock           sync.Mutex
}

func newRouterGroup(engine *Engine, parent *RouterGroup) *RouterGroup {
	return &RouterGroup{
		engine: engine,
		parent: parent,
	}
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) *RouterGroup {
	g.middlewares = append(g.middlewares, middlewares...)
	return g
}

func (g *RouterGroup) insert(parts []string, height int) *RouterGroup {
	if len(parts) == height {
		return g
	}

	part := parts[height]
	child := newRouterGroup(g.engine, g)
	child.prefix = g.prefix + "/" + part
	return child.insert(parts, height+1)
}

func (g *RouterGroup) Get(path string, handlers ...HandlerFunc) *RouterGroup {
	return g.addRouter("GET", g.genFullPath(path), handlers)
}

func (g *RouterGroup) POST(path string, handlers ...HandlerFunc) *RouterGroup {
	return g.addRouter("POST", g.genFullPath(path), handlers)
}

func (g *RouterGroup) PUT(path string, handlers ...HandlerFunc) *RouterGroup {
	return g.addRouter("PUT", g.genFullPath(path), handlers)
}

func (g *RouterGroup) DELETE(path string, handlers ...HandlerFunc) *RouterGroup {
	return g.addRouter("DELETE", g.genFullPath(path), handlers)
}

func (g *RouterGroup) PATCH(path string, handlers ...HandlerFunc) *RouterGroup {
	return g.addRouter("PATCH", g.genFullPath(path), handlers)
}

func (g *RouterGroup) HEAD(path string, handlers ...HandlerFunc) *RouterGroup {
	return g.addRouter("HEAD", g.genFullPath(path), handlers)
}

func (g *RouterGroup) OPTIONS(path string, handlers ...HandlerFunc) *RouterGroup {
	return g.addRouter("OPTIONS", g.genFullPath(path), handlers)
}

func (g *RouterGroup) ANY(path string, handlers ...HandlerFunc) *RouterGroup {
	for _, method := range methods {
		g.addRouter(method, path, handlers)
	}
	return g
}

func (g *RouterGroup) addRouter(method string, pattern string, handlers HandlersChain) *RouterGroup {
	groupHandler := GroupHandler{
		routerGroup: g,
		handlers:    handlers,
	}
	g.engine.router.addRouter(method, pattern, groupHandler)
	return g
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	parts := parsePattern(prefix)
	return g.insert(parts, 0)
}

func (g *RouterGroup) genFullPath(path string) string {
	if g.prefix == "" {
		return path
	}
	return g.prefix + "/" + path
}

func (g *RouterGroup) getNestMiddlewares() []HandlerFunc {
	if g.middlewareInit || g.parent == nil {
		return g.middlewares
	}

	g.lock.Lock()
	defer g.lock.Unlock()
	if g.middlewareInit {
		return g.middlewares
	}

	g.middlewareInit = true
	g.middlewares = append(g.parent.getNestMiddlewares(), g.middlewares...)
	return g.middlewares
}

func (g *RouterGroup) returnObj() *RouterGroup {
	if g.parent == nil {
		return g.engine.RouterGroup
	}
	return g
}

package gee

import (
	"sync"
)

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

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
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

func (g *RouterGroup) Get(path string, handler HandlerFunc) {
	g.addRouter("GET", g.genFullPath(path), handler)
}

func (g *RouterGroup) POST(path string, handler HandlerFunc) {
	g.addRouter("POST", g.genFullPath(path), handler)
}

func (g *RouterGroup) PUT(path string, handler HandlerFunc) {
	g.addRouter("PUT", g.genFullPath(path), handler)
}

func (g *RouterGroup) DELETE(path string, handler HandlerFunc) {
	g.addRouter("DELETE", g.genFullPath(path), handler)
}

func (g *RouterGroup) addRouter(method string, pattern string, handler HandlerFunc) {
	groupHandler := GroupHandler{
		routerGroup: g,
		handler: handler,
	}
	g.engine.router.addRouter(method, pattern, groupHandler)
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

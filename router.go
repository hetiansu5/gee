package gee

import "strings"

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

type GroupHandler struct {
	routerGroup *RouterGroup
	handlers    HandlersChain
}

type Router struct {
	roots    map[string]*node
	handlers map[string]GroupHandler
}

func newRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		handlers: make(map[string]GroupHandler),
	}
}

func (r *Router) addRouter(method string, pattern string, handler GroupHandler) {
	method = formatMethod(method)
	parts := parsePattern(pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	key := formKey(method, pattern)
	r.handlers[key] = handler
}

func (r *Router) getRouter(method string, path string) (*node, map[string]string) {
	method = formatMethod(method)
	parts := parsePattern(path)
	_, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	nod := r.roots[method].search(parts, 0)
	if nod == nil {
		return nil, nil
	}

	return nod, nod.parseParams(parts)
}

func (r *Router) getFuncHandlers(method string, pattern string) []HandlerFunc {
	method = formatMethod(method)
	key := formKey(method, pattern)
	groupHandler := r.handlers[key]
	middlewares := groupHandler.routerGroup.getNestMiddlewares()
	return append(middlewares, groupHandler.handlers...)
}

func formatMethod(method string) string {
	return strings.ToUpper(method)
}

func formKey(method string, pattern string) string {
	return method + "-" + pattern
}

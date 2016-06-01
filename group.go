package xrouter

import (
	"path/filepath"

	"github.com/rs/xhandler"
)

func newGroup(prefix string, chain *xhandler.Chain, r RouterGroup) RouterGroup {
	return &routerGroup{prefix, &xhandler.Chain{}, r}
}

func wrapper(chain *xhandler.Chain, f xhandler.HandlerFuncC) xhandler.HandlerFuncC {
	return chain.HandlerCF(f).ServeHTTPC
}

type routerGroup struct {
	prefix string
	chain  *xhandler.Chain
	router RouterGroup
}

// Use adds middleware to the router.
func (r *routerGroup) Use(f func(next xhandler.HandlerC) xhandler.HandlerC) {
	r.chain.UseC(f)
}

// GET adds a GET handler at the given path.
func (r *routerGroup) GET(path string, handler xhandler.HandlerFuncC) {
	r.router.GET(r.prefix+path, wrapper(r.chain, handler))
}

// POST adds a POST handler at the given path.
func (r *routerGroup) POST(path string, handler xhandler.HandlerFuncC) {
	r.router.POST(r.prefix+path, wrapper(r.chain, handler))
}

// PUT adds a PUT handler at the given path.
func (r *routerGroup) PUT(path string, handler xhandler.HandlerFuncC) {
	r.router.PUT(r.prefix+path, wrapper(r.chain, handler))
}

// OPTIONS adds a OPTIONS handler at the given path.
func (r *routerGroup) OPTIONS(path string, handler xhandler.HandlerFuncC) {
	r.router.OPTIONS(r.prefix+path, wrapper(r.chain, handler))
}

// HEAD adds a HEAD handler at the given path.
func (r *routerGroup) HEAD(path string, handler xhandler.HandlerFuncC) {
	r.router.HEAD(r.prefix+path, wrapper(r.chain, handler))
}

// PATCH adds a PATCH handler at the given path.
func (r *routerGroup) PATCH(path string, handler xhandler.HandlerFuncC) {
	r.router.PATCH(r.prefix+path, wrapper(r.chain, handler))
}

// DELETE adds a DELETE handler at the given path.
func (r *routerGroup) DELETE(path string, handler xhandler.HandlerFuncC) {
	r.router.DELETE(r.prefix+path, wrapper(r.chain, handler))
}

// Group returns a new router which strips the given path before the request is handled. All the middleware from the router is transferred.
func (r *routerGroup) Group(path string) RouterGroup {
	return newGroup(path, r.chain, r)
}

func (r *routerGroup) Path() string {
	return filepath.Clean(r.router.Path() + r.prefix)
}

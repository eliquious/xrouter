package xrouter

import (
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func newGroup(prefix string, chain alice.Chain, r *httprouter.Router, evtHandler func(evt Event)) *routerGroup {
	return &routerGroup{prefix, chain, r, evtHandler}
}

func wrapper(chain alice.Chain, f Route) http.Handler {
	return chain.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		f(r.Context(), w, r)
	})
}

type routerGroup struct {
	prefix     string
	chain      alice.Chain
	router     *httprouter.Router
	evtHandler func(evt Event)
}

// Use adds middleware to the router.
func (r *routerGroup) Use(f func(next http.Handler) http.Handler) {
	r.chain = r.chain.Append(f)
}

// Chain gets the middleware chain.
func (r *routerGroup) Chain() alice.Chain {
	return r.chain
}

// GET adds a GET handler at the given path.
func (r *routerGroup) GET(path string, handler Route) {
	r.router.GET(r.prefix+path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"GET", r.prefix + path})
}

// POST adds a POST handler at the given path.
func (r *routerGroup) POST(path string, handler Route) {
	r.router.POST(r.prefix+path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"POST", r.prefix + path})
}

// PUT adds a PUT handler at the given path.
func (r *routerGroup) PUT(path string, handler Route) {
	r.router.PUT(r.prefix+path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"PUT", r.prefix + path})
}

// OPTIONS adds a OPTIONS handler at the given path.
func (r *routerGroup) OPTIONS(path string, handler Route) {
	r.router.OPTIONS(r.prefix+path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"OPTIONS", r.prefix + path})
}

// HEAD adds a HEAD handler at the given path.
func (r *routerGroup) HEAD(path string, handler Route) {
	r.router.HEAD(r.prefix+path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"HEAD", r.prefix + path})
}

// PATCH adds a PATCH handler at the given path.
func (r *routerGroup) PATCH(path string, handler Route) {
	r.router.PATCH(r.prefix+path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"PATCH", r.prefix + path})
}

// DELETE adds a DELETE handler at the given path.
func (r *routerGroup) DELETE(path string, handler Route) {
	r.router.DELETE(r.prefix+path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"DELETE", r.prefix + path})
}

// Group returns a new router which strips the given path before the request is handled. All the middleware from the router is transferred.
func (r *routerGroup) Group(path string) RouterGroup {
	return newGroup(r.prefix+path, r.chain.Append(), r.router, r.evtHandler)
}

func (r *routerGroup) Path() string {
	return filepath.Clean(r.prefix)
}

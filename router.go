package xrouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/xhandler"
)

type RouterGroup interface {

	// Use adds middleware to the router.
	Use(f func(next xhandler.HandlerC) xhandler.HandlerC)

	// Group returns a new router which strips the given path before the request is handled. All middleware is transferred to the child group.
	Group(path string) RouterGroup

	// GET adds a GET handler at the given path.
	GET(path string, handler xhandler.HandlerFuncC)

	// POST adds a POST handler at the given path.
	POST(path string, handler xhandler.HandlerFuncC)

	// PUT adds a PUT handler at the given path.
	PUT(path string, handler xhandler.HandlerFuncC)

	// OPTIONS adds a OPTIONS handler at the given path.
	OPTIONS(path string, handler xhandler.HandlerFuncC)

	// HEAD adds a HEAD handler at the given path.
	HEAD(path string, handler xhandler.HandlerFuncC)

	// PATCH adds a PATCH handler at the given path.
	PATCH(path string, handler xhandler.HandlerFuncC)

	// DELETE adds a DELETE handler at the given path.
	DELETE(path string, handler xhandler.HandlerFuncC)
}

type Router interface {
	RouterGroup

	// Static adds a directory of static content to serve at root.
	StaticRoot(fs http.Handler)

	// StaticFiles adds a directory of static content to a specific path.
	StaticFiles(path string, fs http.Handler)

	// Handler returns an http.Handler
	Handler() http.Handler
}

// New creates a router which wraps an httprouter.
func New() Router {
	return &router{&xhandler.Chain{}, httprouter.New()}
}

// Router is a simple abstraction on top of httprouter which allows for simpler use of the http.Handler interface from the standard library.
type router struct {
	chain  *xhandler.Chain
	router *httprouter.Router
}

// Use adds middleware to the router.
func (r *router) Use(f func(next xhandler.HandlerC) xhandler.HandlerC) {
	r.chain.UseC(f)
}

// GET adds a GET handler at the given path.
func (r *router) GET(path string, handler xhandler.HandlerFuncC) {
	r.router.GET(path, httpParamsHandler(r.chain, handler))
}

// POST adds a POST handler at the given path.
func (r *router) POST(path string, handler xhandler.HandlerFuncC) {
	r.router.POST(path, httpParamsHandler(r.chain, handler))
}

// PUT adds a PUT handler at the given path.
func (r *router) PUT(path string, handler xhandler.HandlerFuncC) {
	r.router.PUT(path, httpParamsHandler(r.chain, handler))
}

// OPTIONS adds a OPTIONS handler at the given path.
func (r *router) OPTIONS(path string, handler xhandler.HandlerFuncC) {
	r.router.OPTIONS(path, httpParamsHandler(r.chain, handler))
}

// HEAD adds a HEAD handler at the given path.
func (r *router) HEAD(path string, handler xhandler.HandlerFuncC) {
	r.router.HEAD(path, httpParamsHandler(r.chain, handler))
}

// PATCH adds a PATCH handler at the given path.
func (r *router) PATCH(path string, handler xhandler.HandlerFuncC) {
	r.router.PATCH(path, httpParamsHandler(r.chain, handler))
}

// DELETE adds a DELETE handler at the given path.
func (r *router) DELETE(path string, handler xhandler.HandlerFuncC) {
	r.router.DELETE(path, httpParamsHandler(r.chain, handler))
}

// Static adds a directory of static content to serve at root.
func (r *router) StaticRoot(fs http.Handler) {
	r.router.NotFound = HttpHandler(r.chain, fs)
}

// StaticFiles adds a directory of static content to a specific path.
func (r *router) StaticFiles(path string, fs http.Handler) {
	h := http.StripPrefix(path, fs)
	r.router.GET(path, func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, req)
	})
}

// Handler returns an http.Handler
func (r *router) Handler() http.Handler {
	return r.router
}

// Group returns a new router which strips the given path before the request is handled. All the middleware from the router is transferred.
func (r *router) Group(path string) RouterGroup {
	return newGroup(path, r.chain, r)
}

package xrouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/xhandler"
)

// RouterGroup allows for grouping routes with separate middleware.
type RouterGroup interface {

	// Use adds middleware to the router.
	Use(f func(next xhandler.HandlerC) xhandler.HandlerC)

	// Group returns a new router which strips the given path before the request is handled. All middleware is transferred to the child group.
	Group(path string) RouterGroup

	// Path returns the root path of the RouterGroup
	Path() string

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

// Router defines a root router for handling requests.
type Router interface {
	RouterGroup

	// Static adds a directory of static content to serve at root.
	StaticRoot(fs http.Handler)

	// StaticFiles adds a directory of static content to a specific path.
	StaticFiles(path string, fs http.Handler)

	// NotFound adds a handler for routes that don't exist.
	NotFound(http.Handler)

	// MethodNotAllowed handles requests in which the route exists but hte wrong method was used.
	MethodNotAllowed(http.Handler)

	// Handler returns an http.Handler
	Handler() http.Handler

	// EventHandler calls the given function for each handler as it is added to the router.
	EventHandler(func(evt Event))
}

// New creates a router which wraps an httprouter.
func New() Router {
	return &router{&xhandler.Chain{}, httprouter.New(), evtHandler}
}

// Default event handler
func evtHandler(Event) {}

// Router is a simple abstraction on top of httprouter which allows for simpler use of the http.Handler interface from the standard library.
type router struct {
	chain      *xhandler.Chain
	router     *httprouter.Router
	evtHandler func(Event)
}

// Use adds middleware to the router.
func (r *router) Use(f func(next xhandler.HandlerC) xhandler.HandlerC) {
	r.chain.UseC(f)
}

// NotFound adds a handler for unknown routes.
func (r *router) NotFound(h http.Handler) {
	r.router.NotFound = HTTPHandler(r.chain, h)
	r.evtHandler(NotFoundHandlerEvent{})
}

// MethodNotAllowed adds a handler for existing routes and unknown methods.
func (r *router) MethodNotAllowed(h http.Handler) {
	r.router.MethodNotAllowed = HTTPHandler(r.chain, h)
	r.evtHandler(MethodNotAllowedHandlerEvent{})
}

// GET adds a GET handler at the given path.
func (r *router) GET(path string, handler xhandler.HandlerFuncC) {
	r.router.GET(path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"GET", path})
}

// POST adds a POST handler at the given path.
func (r *router) POST(path string, handler xhandler.HandlerFuncC) {
	r.router.POST(path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"POST", path})
}

// PUT adds a PUT handler at the given path.
func (r *router) PUT(path string, handler xhandler.HandlerFuncC) {
	r.router.PUT(path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"PUT", path})
}

// OPTIONS adds a OPTIONS handler at the given path.
func (r *router) OPTIONS(path string, handler xhandler.HandlerFuncC) {
	r.router.OPTIONS(path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"OPTIONS", path})
}

// HEAD adds a HEAD handler at the given path.
func (r *router) HEAD(path string, handler xhandler.HandlerFuncC) {
	r.router.HEAD(path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"HEAD", path})
}

// PATCH adds a PATCH handler at the given path.
func (r *router) PATCH(path string, handler xhandler.HandlerFuncC) {
	r.router.PATCH(path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"PATCH", path})
}

// DELETE adds a DELETE handler at the given path.
func (r *router) DELETE(path string, handler xhandler.HandlerFuncC) {
	r.router.DELETE(path, httpParamsHandler(r.chain, handler))
	r.evtHandler(AddHandlerEvent{"DELETE", path})
}

// Static adds a directory of static content to serve at root. All requests not matched to a route will be handled here. It is an alias to the NotFound method.
func (r *router) StaticRoot(fs http.Handler) {
	r.router.NotFound = HTTPHandler(r.chain, fs)
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

// Path returns the root path for the router.
func (r *router) Path() string {
	return "/"
}

// Path returns the root path for the router.
func (r *router) EventHandler(h func(Event)) {
	r.evtHandler = h
}

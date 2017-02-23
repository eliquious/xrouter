package xrouter

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// RouterGroup allows for grouping routes with separate middleware.
type RouterGroup interface {

	// Use adds middleware to the router.
	Use(f func(next http.Handler) http.Handler)

	// Chain returns the middleware chain.
	Chain() alice.Chain
	
	// Group returns a new router which strips the given path before the request is handled. All middleware is transferred to the child group.
	Group(path string) RouterGroup

	// Path returns the root path of the RouterGroup
	Path() string

	// GET adds a GET handler at the given path.
	GET(path string, handler Route)

	// POST adds a POST handler at the given path.
	POST(path string, handler Route)

	// PUT adds a PUT handler at the given path.
	PUT(path string, handler Route)

	// OPTIONS adds a OPTIONS handler at the given path.
	OPTIONS(path string, handler Route)

	// HEAD adds a HEAD handler at the given path.
	HEAD(path string, handler Route)

	// PATCH adds a PATCH handler at the given path.
	PATCH(path string, handler Route)

	// DELETE adds a DELETE handler at the given path.
	DELETE(path string, handler Route)
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

// Route is a function with exposes the request context as an argument. For Go 1.7+, the request has an attached context.
// This provides for backward compatibility for the many handlers written before the Go 1.7.
type Route func(context.Context, http.ResponseWriter, *http.Request)

// New creates a router which wraps an httprouter.
func New() Router {
	c := alice.New()
	r := httprouter.New()
	return &router{r, newGroup("", c, r, evtHandler)}
}

// Default event handler
func evtHandler(Event) {}

// Router is a simple abstraction on top of httprouter which allows for simpler use of the http.Handler interface from the standard library.
type router struct {
	router *httprouter.Router
	group  *routerGroup
}

// Use adds middleware to the router.
func (r *router) Use(f func(next http.Handler) http.Handler) {
	r.group.Use(f)
}

// Chain returns the middleware chain.
func (r *router) Chain() alice.Chain {
	return r.group.Chain()
}

// NotFound adds a handler for unknown routes.
func (r *router) NotFound(h http.Handler) {
	r.router.NotFound = r.group.chain.Then(h)
	r.group.evtHandler(NotFoundHandlerEvent{})
}

// MethodNotAllowed adds a handler for existing routes and unknown methods.
func (r *router) MethodNotAllowed(h http.Handler) {
	r.router.MethodNotAllowed = r.group.chain.Then(h)
	r.group.evtHandler(MethodNotAllowedHandlerEvent{})
}

// GET adds a GET handler at the given path.
func (r *router) GET(path string, handler Route) {
	r.group.GET(path, handler)
}

// POST adds a POST handler at the given path.
func (r *router) POST(path string, handler Route) {
	r.group.POST(path, handler)
}

// PUT adds a PUT handler at the given path.
func (r *router) PUT(path string, handler Route) {
	r.group.PUT(path, handler)
}

// OPTIONS adds a OPTIONS handler at the given path.
func (r *router) OPTIONS(path string, handler Route) {
	r.group.OPTIONS(path, handler)
}

// HEAD adds a HEAD handler at the given path.
func (r *router) HEAD(path string, handler Route) {
	r.group.HEAD(path, handler)
}

// PATCH adds a PATCH handler at the given path.
func (r *router) PATCH(path string, handler Route) {
	r.group.PATCH(path, handler)
}

// DELETE adds a DELETE handler at the given path.
func (r *router) DELETE(path string, handler Route) {
	r.group.DELETE(path, handler)
}

// StaticRoot adds a directory of static content to serve at root. All requests not matched to a route will be handled here. It is an alias to the NotFound method.
func (r *router) StaticRoot(fs http.Handler) {
	r.router.NotFound = r.group.chain.Then(fs)
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
	return r.group.Group(path)
}

// EventHandler adds a function to handle route events.
func (r *router) EventHandler(h func(Event)) {
	r.group.evtHandler = h
}

func (r *router) Path() string {
	return "/"
}

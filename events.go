package xrouter

// Event represents router events which will fire as the router is being configured. This is used more as a debug tool than anything.
type Event interface{}

// AddHandlerEvent is fired whenever a generic handler is added to the router.
type AddHandlerEvent struct {
	Method string
	Path   string
}

// NotFoundHandlerEvent is fired when a NotFound handler is set for the router.
type NotFoundHandlerEvent struct {
}

// MethodNotAllowedHandlerEvent is fired when a MethodNotAllowed handler is set for the router.
type MethodNotAllowedHandlerEvent struct {
}

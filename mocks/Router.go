package mocks

import "github.com/eliquious/xrouter"
import "github.com/stretchr/testify/mock"

import "net/http"

type Router struct {
	mock.Mock
}

// StaticRoot provides a mock function with given fields: fs
func (_m *Router) StaticRoot(fs http.Handler) {
	_m.Called(fs)
}

// StaticFiles provides a mock function with given fields: path, fs
func (_m *Router) StaticFiles(path string, fs http.Handler) {
	_m.Called(path, fs)
}

// NotFound provides a mock function with given fields: _a0
func (_m *Router) NotFound(_a0 http.Handler) {
	_m.Called(_a0)
}

// MethodNotAllowed provides a mock function with given fields: _a0
func (_m *Router) MethodNotAllowed(_a0 http.Handler) {
	_m.Called(_a0)
}

// Handler provides a mock function with given fields:
func (_m *Router) Handler() http.Handler {
	ret := _m.Called()

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func() http.Handler); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(http.Handler)
	}

	return r0
}

// EventHandler provides a mock function with given fields: _a0
func (_m *Router) EventHandler(_a0 func(xrouter.Event)) {
	_m.Called(_a0)
}

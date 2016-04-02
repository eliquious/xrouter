package mocks

import "github.com/eliquious/xrouter"
import "github.com/stretchr/testify/mock"

import "github.com/rs/xhandler"

type RouterGroup struct {
	mock.Mock
}

// Use provides a mock function with given fields: f
func (_m *RouterGroup) Use(f func(xhandler.HandlerC) xhandler.HandlerC) {
	_m.Called(f)
}

// Group provides a mock function with given fields: path
func (_m *RouterGroup) Group(path string) xrouter.RouterGroup {
	ret := _m.Called(path)

	var r0 xrouter.RouterGroup
	if rf, ok := ret.Get(0).(func(string) xrouter.RouterGroup); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(xrouter.RouterGroup)
	}

	return r0
}

// Path provides a mock function with given fields:
func (_m *RouterGroup) Path() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GET provides a mock function with given fields: path, handler
func (_m *RouterGroup) GET(path string, handler xhandler.HandlerFuncC) {
	_m.Called(path, handler)
}

// POST provides a mock function with given fields: path, handler
func (_m *RouterGroup) POST(path string, handler xhandler.HandlerFuncC) {
	_m.Called(path, handler)
}

// PUT provides a mock function with given fields: path, handler
func (_m *RouterGroup) PUT(path string, handler xhandler.HandlerFuncC) {
	_m.Called(path, handler)
}

// OPTIONS provides a mock function with given fields: path, handler
func (_m *RouterGroup) OPTIONS(path string, handler xhandler.HandlerFuncC) {
	_m.Called(path, handler)
}

// HEAD provides a mock function with given fields: path, handler
func (_m *RouterGroup) HEAD(path string, handler xhandler.HandlerFuncC) {
	_m.Called(path, handler)
}

// PATCH provides a mock function with given fields: path, handler
func (_m *RouterGroup) PATCH(path string, handler xhandler.HandlerFuncC) {
	_m.Called(path, handler)
}

// DELETE provides a mock function with given fields: path, handler
func (_m *RouterGroup) DELETE(path string, handler xhandler.HandlerFuncC) {
	_m.Called(path, handler)
}

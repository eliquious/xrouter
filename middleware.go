package xrouter

import (
	"bufio"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/xlog"

	"golang.org/x/net/context"
)

// ParamsKey is the key for contexts which grant access to the url params.
const ParamsKey = "params"

// Param returns a URL parameter by name
func Param(ctx context.Context, key string) string {
	if params, ok := ctx.Value(ParamsKey).(httprouter.Params); ok {
		return params.ByName(key)
	}
	return ""
}

// httpParamsHandler is middleware which links the middleware and httprouter.
func httpParamsHandler(chain alice.Chain, handler Route) httprouter.Handle {
	h := wrapper(chain, handler)
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		req = req.WithContext(context.WithValue(req.Context(), ParamsKey, params))
		h.ServeHTTP(w, req)
	}
}

// HTTPHandler wraps a raw http.Handler in chain middleware.
func HTTPHandler(c alice.Chain, fs http.Handler) http.Handler {
	return c.Then(fs)
}

// LogHandler instantiates a new xlog HTTP handler using the given log.
func LogHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ptw := passThroughResponseWriter{200, w}
			start := time.Now()
			next.ServeHTTP(&ptw, r)
			xlog.FromContext(r.Context()).Info(http.StatusText(ptw.StatusCode), xlog.F{
				"duration": time.Now().Sub(start).String(),
				"status":   ptw.StatusCode,
			})
		})
	}
}

type passThroughResponseWriter struct {
	StatusCode     int
	ResponseWriter http.ResponseWriter
}

func (p *passThroughResponseWriter) WriteHeader(code int) {
	p.StatusCode = code
	p.ResponseWriter.WriteHeader(code)
}

func (p *passThroughResponseWriter) Header() http.Header {
	return p.ResponseWriter.Header()
}

func (p *passThroughResponseWriter) Write(data []byte) (int, error) {
	return p.ResponseWriter.Write(data)
}

func (p *passThroughResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return p.ResponseWriter.(http.Hijacker).Hijack()
}

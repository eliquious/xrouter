package xrouter

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/xhandler"
	"github.com/rs/xlog"

	"golang.org/x/net/context"
)

// ParamsKey is the key for xhandler contexts which grant access to the url params.
const ParamsKey = "params"

// Param returns a URL parameter by name
func Param(ctx context.Context, key string) string {
	if params, ok := ctx.Value(ParamsKey).(httprouter.Params); ok {
		return params.ByName(key)
	}
	return ""
}

// httpParamsHandler is middleware which links xhandler and httprouter.
func httpParamsHandler(chain *xhandler.Chain, handler xhandler.HandlerFuncC) httprouter.Handle {
	h := chain.HandlerC(xhandler.HandlerFuncC(handler))
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := context.WithValue(context.Background(), ParamsKey, params)
		h.ServeHTTPC(ctx, w, req)
	}
}

// HTTPHandler wraps a raw http.Handler in chain middleware.
func HTTPHandler(c *xhandler.Chain, fs http.Handler) http.Handler {
	return c.Handler(xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

// LogHandler instantiates a new xlog HTTP handler using the given log.
func LogHandler() func(xhandler.HandlerC) xhandler.HandlerC {
	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			ptw := passThroughResponseWriter{200, w}
			start := time.Now()
			next.ServeHTTPC(ctx, &ptw, r)
			xlog.FromContext(ctx).Info(http.StatusText(ptw.StatusCode), xlog.F{
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

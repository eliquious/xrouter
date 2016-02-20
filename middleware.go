package xrouter

import (
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/xhandler"
	"github.com/rs/xlog"

	"golang.org/x/net/context"
)

// ParamsKey is the key for xhandler contexts which grant access to the url params.
const ParamsKey = "params"

// Params returns a URL parameter by name
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

// httpParamsHandler is middleware which links xhandler and httprouter.
func stripPrefixHandler(prefix string, handler xhandler.HandlerFuncC) xhandler.HandlerFuncC {
	if prefix == "" {
		return handler
	}
	return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			r.URL.Path = p
			handler.ServeHTTPC(ctx, w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}

// HttpHandler wraps a raw http.Handler in chain middleware.
func HttpHandler(c *xhandler.Chain, fs http.Handler) http.Handler {
	return c.Handler(xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

// LogHandler instantiates a new xlog HTTP handler using the given log.
func LogHandler(l xlog.Logger) func(xhandler.HandlerC) xhandler.HandlerC {
	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			ptw := passThroughResponseWriter{200, w}
			ctx = xlog.NewContext(ctx, l)
			start := time.Now()
			next.ServeHTTPC(ctx, &ptw, r)
			l.Info(http.StatusText(ptw.StatusCode), xlog.F{
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

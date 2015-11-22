package xrouter

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/xhandler"

	"golang.org/x/net/context"
)

// ParamsKey is the key for xhandler contexts which grant access to the url params.
const ParamsKey = "params"

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

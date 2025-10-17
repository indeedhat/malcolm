package malcolm

import (
	"net/http"
	"path"
	"strings"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Router struct {
	mux        *http.ServeMux
	middleware []Middleware
	basePath   string
}

// NewRouter creates a new router instance with the provided middleware stack assigned
func NewDefaultRouter(middleware ...Middleware) Router {
	return Router{
		mux:        http.DefaultServeMux,
		middleware: middleware,
	}
}

// NewRouter creates a new router instance with the provided middleware stack assigned
func NewRouter(mux *http.ServeMux, middleware ...Middleware) Router {
	return Router{
		mux:        mux,
		middleware: middleware,
	}
}

// ServerMux returns the underlying http.ServeMux instance
func (r Router) ServerMux() *http.ServeMux {
	return r.mux
}

// HandleFunc wraps the base mux HanleFunc method and applies base routes + middleware
func (r Router) HandleFunc(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.mux.HandleFunc(r.prefix(path), r.wrap(handler, middleware...))
}

// Handlewraps the base mux Handle method and applies base routes + middleware
func (r Router) Handle(path string, handler http.Handler, middleware ...Middleware) {
	r.mux.HandleFunc(r.prefix(path), r.wrap(handler.ServeHTTP, middleware...))
}

// Group creates a sub router and assigns a base path and middleware to all routes assigned within it
func (r Router) Group(path string, middleware ...Middleware) Router {
	return Router{
		mux:        r.mux,
		basePath:   r.basePath + path,
		middleware: append(r.middleware, middleware...),
	}
}

// wrap handler with middleware handlers
func (r Router) wrap(handler http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	stack := append(r.middleware, middleware...)

	for i := range stack {
		if stack[len(stack)-1-i] == nil {
			continue
		}

		handler = stack[len(stack)-1-i](handler)
	}

	return handler
}

func (r Router) prefix(uri string) string {
	suffix := ""
	if len(uri) > 1 && uri[len(uri)-1] == '/' {
		suffix = "/"
	}

	if strings.Contains(uri, " ") {
		parts := strings.SplitN(uri, " ", 2)
		if parts[1] == "/" {
			suffix = ""
		}

		return parts[0] + " " + path.Join(r.basePath, strings.Trim(parts[1], " ")) + suffix
	}

	return path.Join(r.basePath, uri) + suffix
}

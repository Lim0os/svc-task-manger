package http_server

import (
	"context"
	"net/http"
	"strings"
)

type Router struct {
	routes map[string]map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]map[string]http.Handler)}
}

func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.Handle("GET", path, handler)
}

func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.Handle("PUT", path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.Handle("POST", path, handler)
}

func (r *Router) Handle(method, path string, handler http.Handler) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]http.Handler)
	}
	r.routes[path][method] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	if handlers, ok := r.routes[path]; ok {
		if handler, ok := handlers[method]; ok {
			handler.ServeHTTP(w, req)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	for route, handlers := range r.routes {
		if r.isWildcardMatch(route, path) || r.isParamMatch(route, path) {
			if handler, ok := handlers[method]; ok {
				if params := r.extractParams(route, path); params != nil {
					req = setPathParams(req, params)
				}
				handler.ServeHTTP(w, req)
				return
			}
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}

	http.NotFound(w, req)
}

func (r *Router) isWildcardMatch(pattern, path string) bool {
	if !strings.HasSuffix(pattern, "/*") {
		return false
	}

	base := strings.TrimSuffix(pattern, "/*")
	return strings.HasPrefix(path, base)
}

func (r *Router) isParamMatch(pattern, path string) bool {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, part := range patternParts {
		if !strings.HasPrefix(part, ":") && part != pathParts[i] {
			return false
		}
	}

	return true
}

func (r *Router) extractParams(pattern, path string) map[string]string {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return nil
	}

	params := make(map[string]string)

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			paramName := strings.TrimPrefix(part, ":")
			params[paramName] = pathParts[i]
		}
	}

	return params
}

func setPathParams(req *http.Request, params map[string]string) *http.Request {
	ctx := req.Context()
	for k, v := range params {
		ctx = context.WithValue(ctx, k, v)
	}
	return req.WithContext(ctx)
}

package router

import (
	"fmt"
	"net/http"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/internal/routerv2/internal/radix"
)

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

type Router interface {
	http.Handler
	Routes() []*Route
	Set(method string, route string, handler http.Handler) error
	Get(route string, handler http.Handler) error
	Post(route string, handler http.Handler) error
	Put(route string, handler http.Handler) error
	Patch(route string, handler http.Handler) error
	Delete(route string, handler http.Handler) error
}

func Provider(in di.Injector) {
	di.Provide[Router](in, provideRouter)
}

func provideRouter(in di.Injector) (Router, error) {
	return New(), nil
}

func New() Router {
	return &router{
		methods: map[string]*method{},
	}
}

type router struct {
	methods map[string]*method
}

type method struct {
	routes map[string]*Route
	tree   *radix.Tree
}

func (r *router) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		method, ok := r.methods[req.Method]
		if !ok {
			next.ServeHTTP(w, req)
			return
		}
		fmt.Println("got method", method)
	})
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Middleware(http.NotFoundHandler()).ServeHTTP(w, req)
}

func (r *router) Routes() (routes []*Route) {
	for _, m := range r.methods {
		for _, route := range m.routes {
			routes = append(routes, route)
		}
	}
	return routes
}

func (r *router) method(name string) *method {
	if m, ok := r.methods[name]; ok {
		return m
	}
	m := &method{
		routes: map[string]*Route{},
		tree:   radix.New(),
	}
	r.methods[name] = m
	return m
}

func (r *router) Set(method string, route string, handler http.Handler) error {
	m := r.method(method)
	if err := m.tree.Insert(route); err != nil {
		return err
	}
	m.routes[route] = &Route{
		Method:  method,
		Path:    route,
		Handler: handler,
	}
	return nil
}

func (r *router) Get(route string, handler http.Handler) error {
	return r.Set(http.MethodGet, route, handler)
}

func (r *router) Post(route string, handler http.Handler) error {
	return r.Set(http.MethodPost, route, handler)
}

func (r *router) Put(route string, handler http.Handler) error {
	return r.Set(http.MethodPut, route, handler)
}

func (r *router) Patch(route string, handler http.Handler) error {
	return r.Set(http.MethodPatch, route, handler)
}

func (r *router) Delete(route string, handler http.Handler) error {
	return r.Set(http.MethodDelete, route, handler)
}

func (r *router) route(method string, route string) *Route {
	m := r.method(method)
	if route, ok := m.routes[route]; ok {
		return route
	}
	m.routes[route] = &Route{
		Method: method,
		Path:   route,
	}
	return m.routes[route]
}

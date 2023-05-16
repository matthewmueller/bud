package router

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/matthewmueller/bud/internal/reflector"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/internal/socket"
	"github.com/matthewmueller/bud/web/router/radix"
)

type Route struct {
	Method string
	Path   string
	Name   string
}

// Interface interface
type Interface interface {
	Set(method, route string, fn func(w http.ResponseWriter, r *http.Request)) error
	Get(route string, fn func(w http.ResponseWriter, r *http.Request)) error
	Post(route string, fn func(w http.ResponseWriter, r *http.Request)) error
	Put(route string, fn func(w http.ResponseWriter, r *http.Request)) error
	Patch(route string, fn func(w http.ResponseWriter, r *http.Request)) error
	Delete(route string, fn func(w http.ResponseWriter, r *http.Request)) error
}

func Provider(in di.Injector) {
	di.Provide[*Router](in, provide)
}

func provide(in di.Injector) (*Router, error) {
	return New(), nil
}

// New router
func New() *Router {
	return &Router{
		methods: map[string]radix.Tree{},
	}
}

// Router struct
type Router struct {
	methods map[string]radix.Tree
	routes  []*Route
}

var _ http.Handler = (*Router)(nil)
var _ Interface = (*Router)(nil)

type Mounter interface {
	Mount(rt Interface)
}

func (rt *Router) Mount(m Mounter) {
	m.Mount(rt)
}

func (rt *Router) Listen(ctx context.Context, address string) error {
	ln, err := socket.Listen(address)
	if err != nil {
		return err
	}
	return rt.Serve(ctx, ln)
}

func (rt *Router) Serve(ctx context.Context, ln net.Listener) error {
	server := &http.Server{
		Addr:    ln.Addr().String(),
		Handler: rt,
	}
	// TODO: use the context to cancel the server
	return server.Serve(ln)
}

// Set a handler to a route
func (rt *Router) Set(method, route string, fn func(w http.ResponseWriter, r *http.Request)) error {
	if !isMethod(method) {
		return fmt.Errorf("router: %q is not a valid HTTP method", method)
	}
	return rt.set(method, route, fn)
}

func (rt *Router) set(method, route string, fn func(w http.ResponseWriter, r *http.Request)) error {
	if route == "/" {
		return rt.insert(method, route, fn)
	}
	// Trim any trailing slash and lowercase the route
	route = strings.TrimRight(strings.ToLower(route), "/")
	return rt.insert(method, route, fn)
}

// Insert the route into the method's radix tree
func (rt *Router) insert(method, route string, fn func(w http.ResponseWriter, r *http.Request)) error {
	if _, ok := rt.methods[method]; !ok {
		rt.methods[method] = radix.New()
	}
	fi, err := reflector.Func(fn)
	if err != nil {
		return err
	}
	rt.routes = append(rt.routes, &Route{
		Method: method,
		Path:   route,
		Name:   fi.Name(),
	})
	return rt.methods[method].Insert(route, fn)
}

// Get route
func (rt *Router) Get(route string, fn func(w http.ResponseWriter, r *http.Request)) error {
	return rt.set(http.MethodGet, route, fn)
}

// Post route
func (rt *Router) Post(route string, fn func(w http.ResponseWriter, r *http.Request)) error {
	return rt.set(http.MethodPost, route, fn)
}

// Put route
func (rt *Router) Put(route string, fn func(w http.ResponseWriter, r *http.Request)) error {
	return rt.set(http.MethodPut, route, fn)
}

// Patch route
func (rt *Router) Patch(route string, fn func(w http.ResponseWriter, r *http.Request)) error {
	return rt.set(http.MethodPatch, route, fn)
}

// Delete route
func (rt *Router) Delete(route string, fn func(w http.ResponseWriter, r *http.Request)) error {
	return rt.set(http.MethodDelete, route, fn)
}

// List all routes
func (rt *Router) List() []*Route {
	return rt.routes
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := rt.Middleware(http.NotFoundHandler())
	handler.ServeHTTP(w, r)
}

// Middleware implements the router middleware
func (rt *Router) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tree, ok := rt.methods[r.Method]
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		// Strip any trailing slash (e.g. /users/ => /users)
		urlPath := trimTrailingSlash(r.URL.Path)
		// Match the path
		match, ok := tree.Match(urlPath)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		// set the slots
		if len(match.Slots) > 0 {
			query := r.URL.Query()
			for _, slot := range match.Slots {
				query.Set(slot.Key, slot.Value)
			}
			r.URL.RawQuery = query.Encode()
		}
		// Call the handler
		match.Handler(w, r)
	})
}

func trimTrailingSlash(path string) string {
	if path == "/" {
		return path
	}
	return strings.TrimRight(path, "/")
}

// isMethod returns true if method is a valid HTTP method
func isMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPost,
		http.MethodPut, http.MethodPatch, http.MethodDelete,
		http.MethodConnect, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

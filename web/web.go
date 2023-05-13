package web

import "net/http"

type Route struct {
	Method string
	Path   string
	Name   string
}

// Router interface
type Router interface {
	http.Handler
	Set(method, route string, handler http.Handler) error
	Get(route string, handler http.Handler) error
	Post(route string, handler http.Handler) error
	Put(route string, handler http.Handler) error
	Patch(route string, handler http.Handler) error
	Delete(route string, handler http.Handler) error
	List() []*Route
}

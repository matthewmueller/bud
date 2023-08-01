package middleware

import (
	"net/http"

	"github.com/matthewmueller/bud/di"
)

func Provider(in di.Injector) {
	di.Provide[Stack](in, provideStack)
}

func provideStack(in di.Injector) (Stack, error) {
	return Stack{
		Function(di.Middleware(in)),
	}, nil
}

type Function func(http.Handler) http.Handler

func (f Function) Middleware(h http.Handler) http.Handler {
	return f(h)
}

type Middleware interface {
	Middleware(http.Handler) http.Handler
}

// compose a stack of middleware into a single middleware
func compose(middlewares ...Middleware) Middleware {
	return Function(func(h http.Handler) http.Handler {
		if len(middlewares) == 0 {
			return h
		}
		for i := len(middlewares) - 1; i >= 0; i-- {
			if middlewares[i] == nil {
				continue
			}
			h = middlewares[i].Middleware(h)
		}
		return h
	})
}

type Stack []Middleware

func (s Stack) Middleware(h http.Handler) http.Handler {
	return compose(s...).Middleware(h)
}

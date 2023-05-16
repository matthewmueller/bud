package middleware

import (
	"net/http"

	"github.com/matthewmueller/bud/di"
)

type Middleware func(http.Handler) http.Handler

func Provider(in di.Injector) {
	di.Provide[Middleware](in, provideMiddleware)
}

func provideMiddleware(in di.Injector) (Middleware, error) {
	return Compose(), nil
}

// Compose a stack of middleware into a single middleware
func Compose(middlewares ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		if len(middlewares) == 0 {
			return h
		}
		for i := len(middlewares) - 1; i >= 0; i-- {
			if middlewares[i] == nil {
				continue
			}
			h = middlewares[i](h)
		}
		return h
	}
}

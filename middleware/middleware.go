package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

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

package di

import (
	"context"
	"errors"
	"net/http"
)

type contextKey string

var key contextKey = "di"

func Middleware(parent Injector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			in := Clone(parent)
			Provide[context.Context](in, func(in Injector) (context.Context, error) {
				return r.Context(), nil
			})
			r = r.WithContext(context.WithValue(r.Context(), key, in))
			next.ServeHTTP(w, r)
		})
	}
}

var ErrNotInContext = errors.New("di: injector is not in the context")

// InjectorFrom returns the injector from a context if present
func InjectorFrom(ctx context.Context) (Injector, error) {
	in, ok := ctx.Value(key).(Injector)
	if !ok {
		return nil, ErrNotInContext
	}
	return in, nil
}

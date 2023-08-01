package handler

import (
	"net/http"

	"github.com/matthewmueller/bud/internal/httpwrap"
)

func Merge(handlers ...http.Handler) http.Handler {
	return httpwrap.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, h := range handlers {
			h.ServeHTTP(w, r)
		}
	}))
}

package welcome

import (
	"net/http"
)

func New() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<h1>Welcome to bud!</h1>"))
	})
}

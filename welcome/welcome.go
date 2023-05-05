package welcome

import (
	"net/http"
)

func New() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to bud!"))
	})
}

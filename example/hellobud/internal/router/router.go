package router

import (
	"net/http"

	"github.com/livebud/buddy/log"
	"github.com/livebud/buddy/router"
	"github.com/livebud/buddy/web"
	"hellobud.com/internal/env"
)

func New(env *env.Env, log log.Log) web.Router {
	r := router.New()
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	}))
	return r
}

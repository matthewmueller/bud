package server

import (
	"net/http"

	"app.com/internal/env"
	"github.com/livebud/buddy/log"
	"github.com/livebud/buddy/router"
	"github.com/livebud/buddy/web"
)

func New(env *env.Env, log log.Log) web.Server {
	r := router.New()
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world..."))
	}))
	return r
}

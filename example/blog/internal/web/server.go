package web

import (
	"database/sql"

	"github.com/livebud/buddy/example/blog/internal/controller/posts"
	"github.com/livebud/buddy/example/blog/internal/view"
	"github.com/livebud/buddy/log"
	"github.com/livebud/buddy/router"
	"github.com/livebud/buddy/transpiler"
	"github.com/livebud/buddy/web"
)

func New(log log.Log, db *sql.DB) web.Server {
	transpiler := transpiler.New()
	viewer := view.New(log, transpiler)
	router := router.New()
	router.Mount(posts.New(db, viewer))
	return router
}

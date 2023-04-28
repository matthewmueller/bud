package view

import (
	"embed"

	"github.com/livebud/buddy/log"
	"github.com/livebud/buddy/transpiler"
	"github.com/livebud/buddy/view"
	"github.com/livebud/buddy/view/gohtml"
)

//go:embed **/*.gohtml
var fsys embed.FS

func New(log log.Log, tr transpiler.Interface) view.Interface {
	return view.New(
		gohtml.New(fsys, log, tr),
	)
}

// type View struct {
// 	gohtml *gohtml.Viewer
// }

// var _ view.Interface = (*View)(nil)

// func (v *View) Render(w http.ResponseWriter, key string, props interface{}) error {
// 	return v.gohtml.Render(w, key, props)
// }

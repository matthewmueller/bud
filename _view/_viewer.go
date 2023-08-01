package view

import (
	"io/fs"
	"strings"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/gomod"
	"github.com/matthewmueller/bud/internal/reflector"
)

func provideViewer(in di.Injector) (*Viewer, error) {
	module, err := di.Load[*gomod.Module](in)
	if err != nil {
		return nil, err
	}
	return &Viewer{
		module:  module,
		fileMap: map[string]fs.FS{},
		pages:   map[string]*Page{},
	}, nil
}

type Viewer struct {
	module  *gomod.Module
	fileMap map[string]fs.FS
	pages   map[string]*Page
}

// type PageSet struct {
// 	fsys fs.
// }

func (v *Viewer) Mount(fsys fs.FS) error {
	modulePath, err := reflector.ModulePath(2)
	if err != nil {
		return err
	}
	relDir := strings.TrimPrefix(modulePath, v.module.Import("view")+"/")
	v.fileMap[relDir] = fsys
	return nil
}

// func (v *Viewer) Add(fsys fs.FS, pages ...*Page) {
// 	v.pages[p.Path] = p
// }

// Intended to be used by
// func (v *Viewer) Render(ctx context.Context, w io.Writer) {

// }

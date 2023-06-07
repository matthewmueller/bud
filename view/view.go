package view

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"testing/fstest"

	"github.com/matthewmueller/bud/view/internal/slot"

	"github.com/matthewmueller/bud/di"
)

type FS = fs.FS

type Input struct {
	Props any
	Slots *slot.Set
}

func (i *Input) Slot() (template.HTML, error) {
	slot, err := io.ReadAll(i.Slots)
	if err != nil {
		return "", err
	}
	return template.HTML(slot), nil
}

type Renderer interface {
	Render(ctx context.Context, w io.Writer, path string, in *Input) error
}

type Finder interface {
	Load(key string) (*Page, error)
}

type Page struct {
	*File
	Frames []*File
	Layout *File
	Error  *File
}

type File struct {
	Path string
	Key  string
	Ext  string
}

type Action interface {
	Layout(handler http.Handler) error
	Frame(handler http.Handler) error
	Error(handler http.Handler) error
}

func Provider(in di.Injector) {
	di.Provide[FS](in, provideFS)
	// di.Provide[Viewer](in, provideViewer)
	di.Provide[Finder](in, provideFinder)
	di.Provide[Renderer](in, provideRenderer)
	di.Provide[Renderers](in, provideRenderers)
	// di.Provide[]
	// di.Provide[Map](in, provideMap)
	// di.Provide[Pages](in, providePages)
	// di.Provide[Action](in, provideAction)
	// di.Provide[View](in, provideView)
}

func provideFS(in di.Injector) (FS, error) {
	return fstest.MapFS{}, nil
}

func provideFinder(in di.Injector) (Finder, error) {
	fsys, err := di.Load[FS](in)
	if err != nil {
		return nil, err
	}
	return &liveFinder{fsys}, nil
}

type liveFinder struct {
	fsys fs.FS
}

func (f *liveFinder) Load(key string) (*Page, error) {
	des, err := fs.ReadDir(f.fsys, path.Dir(key))
	if err != nil {
		return nil, err
	}
	for _, de := range des {
		if de.IsDir() {
			continue
		}
		name := de.Name()
		ext := filepath.Ext(de.Name())
		base := name[:len(name)-len(ext)]
		if !de.IsDir() && base == key {
			// TODO: handle frames and layouts
			return &Page{
				File: &File{
					Path: key,
					Key:  base,
					Ext:  ext,
				},
			}, nil
		}
	}
	return nil, fmt.Errorf("view: unable to load %q key: %w", key, fs.ErrNotExist)
}

func provideRenderers(in di.Injector) (Renderers, error) {
	return Renderers{}, nil
}

func provideRenderer(in di.Injector) (Renderer, error) {
	renderers, err := di.Load[Renderers](in)
	if err != nil {
		return nil, err
	}
	return renderer{renderers}, nil
}

// func provideViewer(in di.Injector) (Viewer, error) {
// 	Finder, err := di.Load[Finder](in)
// 	if err != nil {
// 		return nil, err
// 	}
// 	renderer, err := di.Load[Renderer](in)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &viewer{
// 		Finder:   Finder,
// 		renderer: renderer,
// 	}, nil
// }

type Renderers map[string]Renderer

// var _ = renderer{}

// func (r renderer) Register(extension string, renderer Renderer) {
// 	r[extension] = renderer
// }

type renderer struct {
	r Renderers
}

func (r renderer) Render(ctx context.Context, w io.Writer, path string, in *Input) error {
	ext := filepath.Ext(path)
	renderer, ok := r.r[ext]
	if !ok {
		return fmt.Errorf("view: no renderer for %q", ext)
	}
	return renderer.Render(ctx, w, path, in)
}

// type viewer struct {
// 	Finder   Finder
// 	renderer Renderer
// }

// func (v *viewer) Render(ctx context.Context, w io.Writer, key string, props any) error {
// 	page, err := v.Finder.Load(key)
// 	if err != nil {
// 		return err
// 	}
// 	// TODO: handle frames and layouts
// 	return v.renderer.Render(ctx, w, page.Path, Data{
// 		Props: props,
// 	})
// }

// func providePages(in di.Injector) (Pages, error) {
// 	fsys, err := di.Load[FS](in)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return view.Find(fsys)
// }

// func provideMap(in di.Injector) (Map, error) {
// 	pages, err := di.Load[Pages](in)
// 	if err != nil {
// 		return nil, err
// 	}
// 	m := Map{}
// 	for _, page := range pages {
// 		m[page.Key] = &pageView{page}
// 	}
// 	return m, nil
// }

// type Map map[string]View
// type Pages = map[string]*view.Page

// type pageView struct {
// 	page *view.Page
// }

// func (v *pageView) Render(ctx context.Context, w io.Writer, props any) error {
// 	fmt.Println("rendering", v.page.Key, v.page.Path, props)
// 	// return v.page.Render(ctx, w, props)
// 	return nil
// }

// func provideAction(in di.Injector) (Action, error) {
// 	return actionMap{}, nil
// }

// type actionMap map[string]http.Handler

// var _ Action = actionMap{}

// func (m actionMap) Layout(handler http.Handler) error {
// 	m["layout"] = handler
// 	return nil
// }

// func (m actionMap) Frame(handler http.Handler) error {
// 	info, err := reflector.Func(handler)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("frame", info.Name())
// 	return nil
// }

// func (m actionMap) Error(handler http.Handler) error {
// 	info, err := reflector.Func(handler)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("error", info.Name())
// 	return nil
// }

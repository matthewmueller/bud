package view

import (
	"context"
	"io"

	"github.com/matthewmueller/bud/di"
)

func Provider(in di.Injector) {
	di.Provide[Renderers](in, provideRenderers)
	// di.Provide[View](in, provideView)
	di.Provide[*Viewer](in, provideViewer)
}

func provideRenderers(in di.Injector) (Renderers, error) {
	return Renderers{}, nil
}

// func provideView(in di.Injector) (View, error) {
// 	pages, err := di.Load[Pages](in)
// 	if err != nil {
// 		return nil, err
// 	}
// 	renderers, err := di.Load[Renderers](in)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return New(pages, renderers), nil
// }

// type FS = fs.FS

type Page struct {
	Path   string
	Layout string
	Frames []string
	Error  string
}

type File struct {
	Path string
}

func (p *Page) Files() (files []*File) {
	if p.Layout != "" {
		files = append(files, &File{p.Layout})
	}
	for _, frame := range p.Frames {
		files = append(files, &File{frame})
	}
	files = append(files, &File{p.Path})
	return files
}

type Pages map[string]*Page

type Renderer interface {
	Render(ctx context.Context, w io.Writer, files []*File, props any) error
}

type Renderers map[string]Renderer

// type View interface {
// 	Render(ctx context.Context, w http.ResponseWriter, key string, props any) error
// }

// func New(pages Pages, renderers Renderers) *view {
// 	return &view{pages, renderers}
// }

// type view struct {
// 	pages     Pages
// 	renderers Renderers
// }

// var _ View = (*view)(nil)

// func (v *view) Render(ctx context.Context, w http.ResponseWriter, key string, props any) error {
// 	page, ok := v.pages[key]
// 	if !ok {
// 		return fmt.Errorf("unable to find page for key %q", key)
// 	}
// 	renderer, ok := v.renderers[path.Ext(page.Path)[1:]]
// 	if !ok {
// 		return fmt.Errorf("unable to find renderer for page %q", page.Path)
// 	}
// 	return renderer.Render(ctx, w, page.Files(), props)
// }

// func New(renderers ...Renderer) View {
// 	return &view{renderers}
// }

// type view struct {
// 	renderers []Renderer
// }

// // var _ Interface = (*View)(nil)

// func (v *view) Render(ctx context.Context, key string, props Props) ([]byte, error) {
// 	fmt.Println("rendering", key, props)
// 	return []byte(""), fmt.Errorf("not implemented")
// }

// func (v *view) RenderError(ctx context.Context, key string, props Props, err error) []byte {
// 	fmt.Println("rendering error", key, props, err)
// 	return []byte("not implemented")
// }

// type FS = fs.FS

// type view struct {
// 	viewers []Viewer
// }

// func (v *view) Render(w http.ResponseWriter, key string, props interface{}) error {
// 	for _, viewer := range v.viewers {
// 		html, err := viewer.Render(context.Background(), key, Props{key: props})
// 		if err != nil {
// 			html = viewer.RenderError(context.Background(), key, Props{key: props}, err)
// 		}
// 		_, err = w.Write(html)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

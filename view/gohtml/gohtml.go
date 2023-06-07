package gohtml

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"io/fs"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/logger"
	"github.com/matthewmueller/bud/transpiler"
	"github.com/matthewmueller/bud/view"
)

func Provider(in di.Injector) {
	di.Provide[*Renderer](in, provideRenderer)
	di.Register[view.Renderers](in, registerRenderer)
}

func provideRenderer(in di.Injector) (*Renderer, error) {
	fsys, err := di.Load[view.FS](in)
	if err != nil {
		return nil, err
	}
	log, err := di.Load[logger.Log](in)
	if err != nil {
		return nil, err
	}
	tr, err := di.Load[transpiler.Interface](in)
	if err != nil {
		return nil, err
	}
	return New(fsys, log, tr), nil
}

func registerRenderer(in di.Injector, renderers view.Renderers) error {
	renderer, err := di.Load[*Renderer](in)
	if err != nil {
		return err
	}
	renderers[".gohtml"] = renderer
	return nil
}

// MustParse panics if unable to parse
func MustParse(name, code string) *template.Template {
	template, err := Parse(name, code)
	if err != nil {
		panic(err)
	}
	return template
}

// Parse parses Go code
func Parse(name, code string) (*template.Template, error) {
	return template.New(name).Parse(code)
}

// func registerRenderer(in di.Injector, renderers view.Renderers) error {
// 	renderer, err := di.Load[*Renderer](in)
// 	if err != nil {
// 		return err
// 	}
// 	renderers["gohtml"] = renderer
// 	return nil
// }

func New(fsys fs.FS, log logger.Log, tr transpiler.Interface) *Renderer {
	return &Renderer{fsys, log, tr}
}

type Renderer struct {
	fsys fs.FS
	log  logger.Log
	tr   transpiler.Interface
}

var _ view.Renderer = (*Renderer)(nil)

func (r *Renderer) parseTemplate(ctx context.Context, templatePath string) (*template.Template, error) {
	code, err := fs.ReadFile(r.fsys, templatePath)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to parse template %q. %w", templatePath, err)
	}
	// TODO: don't transpile when embedded
	code, err = r.tr.Transpile(ctx, templatePath, ".gohtml", code)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to transpile %s: %w", templatePath, err)
	}
	tpl, err := template.New(templatePath).Parse(withProps(string(code)))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func withProps(template string) string {
	return `{{ with $.Props }}` + template + `{{ else }}` + template + `{{ end }}`
}

func (r *Renderer) Render(ctx context.Context, w io.Writer, path string, in *view.Input) error {
	tpl, err := r.parseTemplate(ctx, path)
	if err != nil {
		return err
	}
	if in.Props == nil {
		in.Props = struct{}{}
	}
	return tpl.Execute(w, in)
}

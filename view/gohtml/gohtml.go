package gohtml

import (
	"bytes"
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
	renderers["gohtml"] = renderer
	return nil
}

func New(fsys fs.FS, log logger.Log, tr transpiler.Interface) *Renderer {
	return &Renderer{fsys, log, tr}
}

type Renderer struct {
	fsys fs.FS
	log  logger.Log
	tr   transpiler.Interface
}

var _ view.Renderer = (*Renderer)(nil)

func (r *Renderer) parseTemplate(ctx context.Context, templatePath string, funcs template.FuncMap) (*template.Template, error) {
	code, err := fs.ReadFile(r.fsys, templatePath)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to parse template %q. %w", templatePath, err)
	}
	// TODO: don't transpile when embedded
	code, err = r.tr.Transpile(ctx, templatePath, ".gohtml", code)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to transpile %s: %w", templatePath, err)
	}
	tpl, err := template.New(templatePath).Funcs(funcs).Parse(string(code))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func (r *Renderer) Render(ctx context.Context, w io.Writer, files []*view.File, props any) error {
	funcs := template.FuncMap{
		"slot": func() template.HTML {
			return ""
		},
	}
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		r.log.Debugf("gohtml: parsing %s", file.Path)
		tpl, err := r.parseTemplate(ctx, file.Path, funcs)
		if err != nil {
			return err
		}
		// Write out the last template
		if i == 0 {
			r.log.Debugf("gohtml: rendering %s", file.Path)
			return tpl.Execute(w, props)
		}
		html := new(bytes.Buffer)
		r.log.Debugf("gohtml: rendering %s", file.Path)
		if err := tpl.Execute(html, props); err != nil {
			return err
		}
		funcs["slot"] = func() template.HTML {
			return template.HTML(html.String())
		}
	}
	return nil
}

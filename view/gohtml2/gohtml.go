package gohtml

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"io/fs"

	"github.com/livebud/buddy/log"
	"github.com/livebud/buddy/transpiler"
)

func New(fsys fs.FS, log log.Log, tr transpiler.Interface) *Viewer {
	return &Viewer{fsys, log, tr}
}

type Viewer struct {
	fsys fs.FS
	log  log.Log
	tr   transpiler.Interface
}

type Files []File

type File struct {
	Path  string
	Props any
}

func (v *Viewer) parseTemplate(ctx context.Context, templatePath string, funcs template.FuncMap) (*template.Template, error) {
	code, err := fs.ReadFile(v.fsys, templatePath)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to parse template %q. %w", templatePath, err)
	}
	// TODO: don't transpile when embedded
	code, err = v.tr.Transpile(ctx, templatePath, ".gohtml", code)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to transpile %s: %w", templatePath, err)
	}
	tpl, err := template.New(templatePath).Funcs(funcs).Parse(string(code))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func (v *Viewer) Render(ctx context.Context, w io.Writer, files Files) error {
	funcs := template.FuncMap{
		"slot": func() template.HTML {
			return ""
		},
	}
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		v.log.Debugf("gohtml: parsing %s", file.Path)
		tpl, err := v.parseTemplate(ctx, file.Path, funcs)
		if err != nil {
			return err
		}
		// Write out the last template
		if i == 0 {
			v.log.Debugf("gohtml: rendering %s", file.Path)
			return tpl.Execute(w, file.Props)
		}
		html := new(bytes.Buffer)
		v.log.Debugf("gohtml: rendering %s", file.Path)
		if err := tpl.Execute(html, file.Props); err != nil {
			return err
		}
		funcs["slot"] = func() template.HTML {
			return template.HTML(html.String())
		}
	}
	return nil
}

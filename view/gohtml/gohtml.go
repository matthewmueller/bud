package gohtml

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/livebud/buddy/log"
	"github.com/livebud/buddy/transpiler"
	"github.com/livebud/buddy/view"
)

func New(fsys fs.FS, log log.Log, tr transpiler.Interface) *Viewer {
	return &Viewer{fsys, log, tr}
}

type Viewer struct {
	fsys fs.FS
	log  log.Log
	tr   transpiler.Interface
}

var _ view.Viewer = (*Viewer)(nil)

func (v *Viewer) parseTemplate(ctx context.Context, templatePath string) (*template.Template, error) {
	// TODO: decide if we want to scope to the view path or module path
	code, err := fs.ReadFile(v.fsys, templatePath)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to parse template %q. %w", templatePath, err)
	}
	// TODO: don't transpile when embedded
	code, err = v.tr.Transpile(ctx, templatePath, ".gohtml", code)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to transpile %s: %w", templatePath, err)
	}
	tpl, err := template.New(templatePath).Parse(string(code))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func (v *Viewer) render(ctx context.Context, templatePath string, props interface{}) ([]byte, error) {
	tpl, err := v.parseTemplate(ctx, templatePath)
	if err != nil {
		return nil, err
	}
	return render(ctx, tpl, props)
}

func render(ctx context.Context, tpl *template.Template, props interface{}) ([]byte, error) {
	out := new(bytes.Buffer)
	// TODO: pass context through
	if err := tpl.Execute(out, props); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (v *Viewer) Render(ctx context.Context, key string, props view.Props) ([]byte, error) {
	// page, ok := v.pages[key]
	// if !ok {
	// 	return nil, fmt.Errorf("gohtml: unable to find page from key %q", key)
	// }
	// v.log.Info("gohtml: rendering", page.Path)
	// html, err := v.render(ctx, page.Path, propMap[page.Key])
	// if err != nil {
	// 	return nil, err
	// }
	// for _, frame := range page.Frames {
	// 	// TODO: support other props
	// 	html, err = v.render(ctx, frame.Path, template.HTML(html))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// if page.Layout != nil {
	// 	// TODO: support other props
	// 	html, err = v.render(ctx, page.Layout.Path, template.HTML(html))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// return html, nil
	return nil, fmt.Errorf("gohtml: render not implemented yet")
}

func (v *Viewer) RenderError(ctx context.Context, key string, props view.Props, originalError error) []byte {
	return []byte("gohtml: render error not implemented yet to render: " + originalError.Error())
	// page, ok := v.pages[key]
	// if !ok {
	// 	return []byte(fmt.Sprintf("gohtml: unable to find page from key %q to render error. %s", key, originalError))
	// }
	// if page.Error == nil {
	// 	return []byte(fmt.Sprintf("gohtml: page %q has no error page to render error. %s", key, originalError))
	// }
	// errorPage, ok := v.pages[page.Error.Key]
	// if !ok {
	// 	return []byte(fmt.Sprintf("gohtml: unable to find error page from key %q to render error. %s", page.Error.Key, originalError))
	// }
	// v.log.Info("gohtml: rendering error", errorPage.Path)
	// errorEntry, err := v.parseTemplate(errorPage.Path)
	// if err != nil {
	// 	return []byte(fmt.Sprintf("gohtml: unable to read error template %q to render error %s. %s", errorPage.Path, err, originalError))
	// }
	// frames := make([]*template.Template, len(errorPage.Frames))
	// for i, frame := range errorPage.Frames {
	// 	frameEntry, err := v.parseTemplate(frame.Path)
	// 	if err != nil {
	// 		return []byte(fmt.Sprintf("gohtml: unable to read frame template %q to render error %s. %s", frame.Path, err, originalError))
	// 	}
	// 	frames[i] = frameEntry
	// }
	// layout, err := v.parseTemplate(errorPage.Layout.Path)
	// if err != nil {
	// 	return []byte(fmt.Sprintf("gohtml: unable to parse layout template %q to render error %s. %s", errorPage.Path, err, originalError))
	// }
	// html, err := render(ctx, errorEntry, viewer.Error(originalError))
	// if err != nil {
	// 	return []byte(fmt.Sprintf("gohtml: unable to render error template %q to render error %s. %s", errorPage.Path, err, originalError))
	// }
	// for i, frame := range errorPage.Frames {
	// 	// TODO: support other props
	// 	html, err = render(ctx, frames[i], template.HTML(html))
	// 	if err != nil {
	// 		return []byte(fmt.Sprintf("gohtml: unable to render frame template %q to render error %s. %s", frame.Path, err, originalError))
	// 	}
	// }
	// html, err = render(ctx, layout, template.HTML(html))
	// if err != nil {
	// 	return []byte(fmt.Sprintf("gohtml: unable to render layout template %q to render error %s. %s", errorPage.Path, err, originalError))
	// }
	// return html
}

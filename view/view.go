package view

import (
	"context"
	"net/http"
)

type Props = map[string]interface{}

type Viewer interface {
	Render(ctx context.Context, key string, props Props) ([]byte, error)
	RenderError(ctx context.Context, key string, props Props, err error) []byte
}

type Interface interface {
	Render(w http.ResponseWriter, key string, props interface{}) error
}

func New(viewers ...Viewer) Interface {
	return &view{viewers}
}

type view struct {
	viewers []Viewer
}

func (v *view) Render(w http.ResponseWriter, key string, props interface{}) error {
	for _, viewer := range v.viewers {
		html, err := viewer.Render(context.Background(), key, Props{key: props})
		if err != nil {
			html = viewer.RenderError(context.Background(), key, Props{key: props}, err)
		}
		_, err = w.Write(html)
		if err != nil {
			return err
		}
	}
	return nil
}

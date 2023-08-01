package view

import (
	"context"
	"io"
	"io/fs"

	"github.com/matthewmueller/bud/di"
)

func Provider(in di.Injector) {

}

type FS = fs.FS

type View interface {
	Render(ctx context.Context, w io.Writer, key string, props any) error
}

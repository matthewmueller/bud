package view_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/matthewmueller/bud/injector"
	"github.com/matthewmueller/bud/view/internal/slot"
	"golang.org/x/sync/errgroup"

	"github.com/matryer/is"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/internal/virtual"
	"github.com/matthewmueller/bud/view"
)

func TestIndex(t *testing.T) {
	is := is.New(t)
	in := injector.New()
	di.Provide[view.FS](in, func(in di.Injector) (view.FS, error) {
		return virtual.Tree{
			"index.gohtml": &virtual.File{
				Data: []byte("Hello {{ .Planet }}!"),
			},
		}, nil
	})
	viewer, err := di.Load[view.Renderer](in)
	is.NoErr(err)
	html := new(bytes.Buffer)
	err = viewer.Render(context.Background(), html, "index.gohtml", &view.Input{
		Props: map[string]interface{}{
			"Planet": "Earth",
		},
	})
	is.NoErr(err)
	is.Equal(html.String(), "Hello Earth!")
}

func TestMultiple(t *testing.T) {
	is := is.New(t)
	in := injector.New()
	di.Provide[view.FS](in, func(in di.Injector) (view.FS, error) {
		return virtual.Tree{
			"index.gohtml": &virtual.File{
				Data: []byte("Hello {{ .Planet }}!"),
			},
			"frame.gohtml": &virtual.File{
				Data: []byte("<main>{{ $.Slot }}</main>"),
			},
			"layout.gohtml": &virtual.File{
				Data: []byte("<html><body>{{ $.Slot }}</body></html>"),
			},
		}, nil
	})
	viewer, err := di.Load[view.Renderer](in)
	is.NoErr(err)
	html := new(bytes.Buffer)
	eg := new(errgroup.Group)
	indexSet := slot.New()
	frameSet := indexSet.New()
	layoutSet := frameSet.New()
	eg.Go(func() error {
		defer indexSet.Close()
		return viewer.Render(context.Background(), indexSet, "index.gohtml", &view.Input{
			Props: map[string]interface{}{
				"Planet": "Earth",
			},
			Slots: indexSet,
		})
	})
	eg.Go(func() error {
		defer frameSet.Close()
		return viewer.Render(context.Background(), frameSet, "frame.gohtml", &view.Input{
			Props: map[string]interface{}{},
			Slots: frameSet,
		})
	})
	eg.Go(func() error {
		defer layoutSet.Close()
		return viewer.Render(context.Background(), html, "layout.gohtml", &view.Input{
			Props: map[string]interface{}{},
			Slots: layoutSet,
		})
	})
	is.NoErr(eg.Wait())
	is.NoErr(err)
	is.Equal(html.String(), "<html><body><main>Hello Earth!</main></body></html>")
}

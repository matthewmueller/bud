package web_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/hexops/valast"
	"github.com/matryer/is"
	"golang.org/x/sync/errgroup"
)

func IO[In, Out any](fn func(ctx context.Context, in In) (Out, error)) func(ctx context.Context, r *http.Request) (any, error) {
	return func(ctx context.Context, r *http.Request) (any, error) {
		var in In
		return fn(ctx, in)
	}
}

type LayoutIn struct {
}

type LayoutOut struct {
	Title string
}

func Layout(ctx context.Context, in *LayoutIn) (*LayoutOut, error) {
	return &LayoutOut{
		Title: "hello",
	}, nil
}

type FrameIn struct {
}

type FrameOut struct {
	Theme string
}

func Frame(ctx context.Context, in *FrameIn) (*FrameOut, error) {
	// time.Sleep(1 * time.Second)
	return &FrameOut{
		Theme: "light",
	}, nil
}

type ViewIn struct {
}

type ViewOut struct {
	Categories []string
}

func View(ctx context.Context, in *ViewIn) (*ViewOut, error) {
	// time.Sleep(1 * time.Second)
	return &ViewOut{
		Categories: []string{"soccer", "finance"},
	}, nil
}

func Load(r *http.Request, m map[string]func(ctx context.Context, r *http.Request) (any, error)) (result map[string]any, err error) {
	mu := sync.Mutex{}
	result = map[string]any{}
	eg, ctx := errgroup.WithContext(r.Context())
	for name, fn := range m {
		name := name
		fn := fn
		eg.Go(func() error {
			value, err := fn(ctx, r)
			if err != nil {
				return err
			}
			mu.Lock()
			result[name] = value
			mu.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}

type Map map[string]func(ctx context.Context, r *http.Request) (any, error)

func TestAction(t *testing.T) {
	is := is.New(t)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	is.NoErr(err)
	propMap, err := Load(req, Map{
		"layout.gohtml": IO(Layout),
		"frame.gohtml":  IO(Frame),
		"index.gohtml":  IO(View),
	})
	is.NoErr(err)
	fmt.Println(valast.String(propMap))
}

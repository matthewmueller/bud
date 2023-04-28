package controller

import (
	"net/http"

	"github.com/livebud/buddy/view"
)

// func New(view view.Viewer) *Controller {
// 	return &Controller{view}
// }

// type Controller struct {
// 	view view.Viewer
// }

// // func (c *Controller) Page(fn func(*Context[any, any]) error) http.Handler {
// // 	return &action[any, any]{}
// // }

// type Context[In, Out any] struct {
// }

// // func Action[In, Out any](fn func(*Context[In, Out]) error) http.Handler {
// // 	return &action[In, Out]{fn: fn, name: "action"}
// // }

// type action[In, Out any] struct {
// 	fn   func(*Context[In, Out]) error
// 	name string
// }

// var _ http.Handler = (*action[any, any])(nil)

// func (a *action[In, Out]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	a.fn(&Context[In, Out]{})
// }

// func (a *action[In, Out]) String() string {
// 	return a.name
// }

// // func (c *Controller)Action[In, Out any](fn func(*Context[In, Out]) error) *Action {
// // 	return &Action{}
// // }

// type Action struct {
// }

// func (a *Action) ServeHTTP(w http.ResponseWriter, r *http.Request) {

// }

func Action[In, Out any](view view.Interface, fn func(*Context[In, Out]) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context[In, Out]{
			Request: r,
			view:    view,
			w:       w,
		}
		_ = ctx
		w.Write([]byte("controller action!"))
	})
}

type Context[In, Out any] struct {
	*http.Request
	view view.Interface
	w    http.ResponseWriter
	Map  In
}

func (c *Context[In, Out]) Accepts(types ...string) bool {
	return true
}

func (c *Context[In, Out]) Status(status int) *Context[In, Out] {
	return c
}

func (c *Context[In, Out]) Render(path string, out Out) error {
	return nil
}

func (c *Context[In, Out]) Redirect(path string) error {
	return nil
}

func (c *Context[In, Out]) JSON(out Out) error {
	return nil
}

func (c *Context[In, Out]) HTML(out Out) error {
	return nil
}

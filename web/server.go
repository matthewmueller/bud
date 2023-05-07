package web

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/internal/signals"
	"github.com/matthewmueller/bud/internal/socket"
)

type Handler = http.Handler
type Server http.Server

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}

func (s *Server) Serve(ctx context.Context, ln net.Listener) error {
	return Serve(ctx, (*http.Server)(s), ln)
}

func (s *Server) Listen(ctx context.Context, address string) error {
	ln, err := socket.Listen(address)
	if err != nil {
		return err
	}
	return s.Serve(ctx, ln)
}

// Serve the handler at address
func Serve(ctx context.Context, server *http.Server, ln net.Listener) error {
	// Make the server shutdownable
	shutdown := shutdown(ctx, server)
	// Serve requests
	if err := server.Serve(ln); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	// Handle any errors that occurred while shutting down
	if err := <-shutdown; err != nil {
		if !errors.Is(err, context.Canceled) {
			return err
		}
	}
	return nil
}

// Shutdown the server when the context is canceled
func shutdown(ctx context.Context, server *http.Server) <-chan error {
	shutdown := make(chan error, 1)
	go func() {
		<-ctx.Done()
		// Wait for one more interrupt to force an immediate shutdown
		forceCtx := signals.Trap(context.Background(), os.Interrupt)
		if err := server.Shutdown(forceCtx); err != nil {
			shutdown <- err
		}
		close(shutdown)
	}()
	return shutdown
}

// Format a listener
func Format(l net.Listener) string {
	address := l.Addr().String()
	if l.Addr().Network() == "unix" {
		return address
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		// Give up trying to format.
		// TODO: figure out if this can occur.
		return address
	}
	// https://serverfault.com/a/444557
	if host == "::" {
		host = "0.0.0.0"
	}
	return fmt.Sprintf("http://%s:%s", host, port)
}

// type Server interface {
// 	http.Handler
// 	Serve(ctx context.Context, ln net.Listener) error
// 	Listen(ctx context.Context, address string) error
// }

type Request[In any] struct {
	r *http.Request
	di.Injector
	Params In
}

// type Router = router.Interface

// type Response[Out any] interface {
// 	// http.ResponseWriter
// 	Render(out Out) error
// }

// type Viewer[Out any] struct {
// 	http.ResponseWriter
// }

// var _ Response[any] = (*Viewer[any])(nil)

// func (v *Viewer[Out]) Render(out Out) error {
// 	v.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	v.ResponseWriter.Write([]byte("hello"))
// 	return nil
// }

// func NewRouter() *Router {
// 	return &Router{}
// }

// type Router = router.Interface

// type Router struct {
// }

// func (r *Router) Get(path string) *Route[any, any] {
// 	return &Route[any, any]{}
// }

// type Middleware = middleware.Middleware

// type Route[In, Out any] struct {
// }

// func (r *Route[In, Out]) Action(fn func(*Request[In], Response[Out]) error) *Route[In, Out] {
// 	return r
// }

// type Func[In, Out any] interface {
// 	func(w http.ResponseWriter, r *http.Request) |
// 		func(w http.ResponseWriter, r *http.Request) error |
// 		func(*Request[In], Response[Out]) error
// }

// func Handler[In, Out any, F Func[In, Out]](fn F) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("hello"))
// 	})
// }

// func Handler(fn func(w http.ResponseWriter, r *http.Request) error) http.Handler {

// 	return http.HandlerFunc(fn)
// }

// func Action2[In, Out any](view view.Interface, fn func(context.Context, In) (Out, error)) http.Handler {
// 	return &action2[In, Out]{fn}
// }

// type action2[In, Out any] struct {
// 	fn func(context.Context, In) (Out, error)
// }

// func (a *action2[In, Out]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	var in In
// 	if err := request.Unmarshal(r, &in); err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}
// 	out, err := a.fn(r.Context(), in)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}
// 	_ = out
// 	w.Write([]byte("hello"))
// }

// func Action[In, Out any](fn func(*Request[In], Response[Out]) error) http.Handler {
// 	return &action[In, Out]{fn: fn, name: "action"}
// }

// type action[In, Out any] struct {
// 	fn   func(*Request[In], Response[Out]) error
// 	name string
// }

// var _ http.Handler = (*action[any, any])(nil)

// func (a *action[In, Out]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	a.fn(&Request[In]{}, &Viewer[Out]{w})
// }

// func (a *action[In, Out]) String() string {
// 	return a.name
// }

// func funcKey(funcs ...interface{}) string {
// 	names := []string{}
// 	for _, f := range funcs {
// 		if n, ok := f.(RouteInfo); ok {
// 			names = append(names, n.HandlerName)
// 			continue
// 		}
// 		rv := reflect.ValueOf(f)
// 		ptr := rv.Pointer()
// 		keyMapMutex.Lock()
// 		if n, ok := keyMap[ptr]; ok {
// 			keyMapMutex.Unlock()
// 			names = append(names, n)
// 			continue
// 		}
// 		keyMapMutex.Unlock()
// 		n := ptrName(ptr)
// 		keyMapMutex.Lock()
// 		keyMap[ptr] = n
// 		keyMapMutex.Unlock()
// 		names = append(names, n)
// 	}
// 	return strings.Join(names, funcKeyDelimeter)
// }

// func ptrName(ptr uintptr) string {
// 	fnc := runtime.FuncForPC(ptr)
// 	n := fnc.Name()

// 	n = strings.Replace(n, "-fm", "", 1)
// 	n = strings.Replace(n, "(", "", 1)
// 	n = strings.Replace(n, ")", "", 1)
// 	return n
// }

// func Router() *router.Router {
// 	return router.New()
// }

// type Router struct {
// 	*router.Router
// }

// var _ Server = (*Router)(nil)

// func (r *Router) Listen(ctx context.Context, address string) error {
// 	ln, err := socket.Listen(address)
// 	if err != nil {
// 		return err
// 	}
// 	return r.Serve(ctx, ln)
// }

// func (r *Router) Serve(ctx context.Context, ln net.Listener) error {
// 	return r.Router.Serve(ctx, ln)
// }

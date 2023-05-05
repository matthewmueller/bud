package serve

import (
	"context"

	"github.com/livebud/buddy/middleware"
	"github.com/livebud/buddy/router"
	"github.com/livebud/buddy/welcome"

	"github.com/livebud/buddy/di"
	"github.com/livebud/buddy/internal/socket"
	"github.com/livebud/buddy/log"
	"github.com/livebud/buddy/web"

	"github.com/livebud/buddy/cli"
)

func New(in di.Injector) *Command {
	di.Provide(in, provideRouter)
	di.Provide(in, provideMiddleware)
	di.Provide(in, provideHandler)
	di.Provide(in, provideServer)
	return &Command{in: in}
}

type Command struct {
	in      di.Injector
	Address string
}

func (c *Command) Mount(cmd cli.Command) {
	cmd.Arg("address").String(&c.Address).Default(":3000")
	cmd.Run(c.Serve)
}

func (c *Command) Serve(ctx context.Context) error {
	log, err := di.Load[log.Log](c.in)
	if err != nil {
		return err
	}
	server, err := di.Load[*web.Server](c.in)
	if err != nil {
		return err
	}
	ln, err := socket.Listen(c.Address)
	if err != nil {
		return err
	}
	log.Infof("Listening on %s", web.Format(ln))
	return server.Serve(ctx, ln)
}

func provideRouter(in di.Injector) (web.Router, error) {
	r := router.New()
	r.Get("/", welcome.New())
	return r, nil
}

func provideMiddleware(in di.Injector) (web.Middleware, error) {
	return middleware.Compose(), nil
}

func provideHandler(in di.Injector) (web.Handler, error) {
	router, err := di.Load[web.Router](in)
	if err != nil {
		return nil, err
	}
	middleware, err := di.Load[web.Middleware](in)
	if err != nil {
		return nil, err
	}
	return middleware(router), nil
}

func provideServer(in di.Injector) (*web.Server, error) {
	handler, err := di.Load[web.Handler](in)
	if err != nil {
		return nil, err
	}
	server := &web.Server{
		Handler: handler,
	}
	return server, nil
}

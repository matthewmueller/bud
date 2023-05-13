package serve

import (
	"context"
	"net"

	"github.com/matthewmueller/bud/env"
	"github.com/matthewmueller/bud/middleware"
	"github.com/matthewmueller/bud/web/router"
	"github.com/matthewmueller/bud/welcome"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/logger"
	"github.com/matthewmueller/bud/web"

	"github.com/matthewmueller/bud/cli"
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
	cmd.Arg("address").String(&c.Address).Default(":" + env.Or("PORT", "3000"))
	cmd.Run(c.Serve)
}

func (c *Command) Serve(ctx context.Context) error {
	log, err := di.Load[logger.Log](c.in)
	if err != nil {
		return err
	}
	server, err := di.Load[*web.Server](c.in)
	if err != nil {
		return err
	}
	ln, err := net.Listen("tcp", c.Address)
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

func provideMiddleware(in di.Injector) (middleware.Middleware, error) {
	return middleware.Compose(
		di.Middleware(in),
	), nil
}

func provideHandler(in di.Injector) (web.Handler, error) {
	router, err := di.Load[web.Router](in)
	if err != nil {
		return nil, err
	}
	middleware, err := di.Load[middleware.Middleware](in)
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

package serve

import (
	"context"
	"net"

	"github.com/matthewmueller/bud/env"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/logger"
	"github.com/matthewmueller/bud/web"

	"github.com/matthewmueller/bud/internal/cli"
)

func Provide(in di.Injector) (*Command, error) {
	log, err := di.Load[logger.Log](in)
	if err != nil {
		return nil, err
	}
	server, err := di.Load[*web.Server](in)
	if err != nil {
		return nil, err
	}
	return New(log, server), nil
}

func Register(in di.Injector, cli *cli.CLI) error {
	cmd, err := di.Load[*Command](in)
	if err != nil {
		return err
	}
	sub := cli.Command("serve", "serve the app")
	sub.Arg("address").String(&cmd.Address).Default(":" + env.Or("PORT", "3000"))
	sub.Run(cmd.Serve)
	return nil
}

func New(log logger.Log, server *web.Server) *Command {
	return &Command{
		log:    log,
		server: server,
	}
}

type Command struct {
	log     logger.Log
	server  *web.Server
	Address string
}

func (c *Command) Serve(ctx context.Context) error {
	ln, err := net.Listen("tcp", c.Address)
	if err != nil {
		return err
	}
	c.log.Infof("Listening on %s", web.Format(ln))
	return c.server.Serve(ctx, ln)
}

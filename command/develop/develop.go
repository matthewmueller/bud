package develop

import (
	"context"

	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/command/serve"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/env"
)

func New(in di.Injector) *Command {
	serve := serve.New(in)
	return &Command{serve}
}

type Command struct {
	serve *serve.Command
}

func (c *Command) Mount(cmd cli.Command) {
	cmd.Arg("address").String(&c.serve.Address).Default(":" + env.Or("PORT", "3000"))
	cmd.Run(c.Develop)
}

func (c *Command) Develop(ctx context.Context) error {
	return c.serve.Serve(ctx)
}

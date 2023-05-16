package develop

import (
	"context"

	"github.com/matthewmueller/bud/command/serve"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/env"
	"github.com/matthewmueller/bud/internal/cli"
)

func Provide(in di.Injector) (*Command, error) {
	serve, err := di.Load[*serve.Command](in)
	if err != nil {
		return nil, err
	}
	return New(serve), nil
}

func Register(in di.Injector, cli *cli.CLI) error {
	cmd, err := di.Load[*Command](in)
	if err != nil {
		return err
	}
	cli.Arg("address").String(&cmd.serve.Address).Default(":" + env.Or("PORT", "3000"))
	cli.Run(cmd.Develop)
	return nil
}

func New(serve *serve.Command) *Command {
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

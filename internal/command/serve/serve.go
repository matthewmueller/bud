package serve

import (
	"context"

	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/internal/command"
)

func New(cmd *command.Command) *Command {
	return &Command{Command: cmd}
}

type Command struct {
	*command.Command
	Address string
}

func (c *Command) Mount(cmd cli.Command) {
	cmd.Arg("address").String(&c.Address).Default(":3000")
	cmd.Run(c.Serve)
}

func (c *Command) Serve(ctx context.Context) error {
	cmd := c.Shell()
	cmd.Run(ctx, "go", "run", "main.go")
	return nil
}

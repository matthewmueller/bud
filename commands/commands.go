package commands

import (
	"github.com/livebud/buddy/cli"
	"github.com/livebud/buddy/command/serve"
	"github.com/livebud/buddy/di"
	"github.com/livebud/buddy/log"
)

func New(in di.Injector) *Command {
	di.Provide(in, provideLog)
	return &Command{in}
}

type Command struct {
	in di.Injector
}

var _ cli.Mounter = (*Command)(nil)

func (c *Command) Mount(cmd cli.Command) {
	cmd.Command("serve", "serve app").Mount(serve.New(c.in))
}

func provideLog(in di.Injector) (log.Log, error) {
	return log.Default(), nil
}

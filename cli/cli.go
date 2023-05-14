package cli

import (
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/internal/cli"
)

type CLI = cli.CLI
type Command = cli.Command
type Mounter = cli.Mounter

func Provider(in di.Injector) {
	di.Provide[Command](in, provide)
	// TODO: figure out a way to alias Command to *CLI
	// then switch back to *CLI
}

func provide(in di.Injector) (Command, error) {
	return Default(), nil
}

func Default() *CLI {
	return New("app", "app cli")
}

func New(name, usage string) *CLI {
	cli := cli.New(name, usage)
	return cli
}

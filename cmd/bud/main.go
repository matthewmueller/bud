package main

import (
	"context"
	"os"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/internal/cli"
	"github.com/matthewmueller/bud/internal/command"
	"github.com/matthewmueller/bud/internal/command/new"
	"github.com/matthewmueller/bud/internal/command/serve"
	"github.com/matthewmueller/bud/logger"
	"github.com/matthewmueller/bud/program"
)

func main() {
	os.Exit(program.Run(run))
}

type Command struct {
	Dir string
}

func run(ctx context.Context, in di.Injector, args ...string) error {
	cmd := &command.Command{}
	cli := cli.New("bud", "bud web framework")
	cli.Flag("chdir", "change the working dir").Short('C').String(&cmd.Dir).Default(".")
	cli.Command("new", "create a new bud project").Mount(new.New(logger.Default()))
	cli.Command("serve", "serve your app").Mount(serve.New(cmd))
	return cli.Parse(ctx, args...)
}

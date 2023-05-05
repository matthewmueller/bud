package main

import (
	"context"
	"os"

	"github.com/livebud/buddy/cli"
	"github.com/livebud/buddy/internal/command"
	"github.com/livebud/buddy/internal/command/new"
	"github.com/livebud/buddy/internal/command/serve"
	"github.com/livebud/buddy/log"
	"github.com/livebud/buddy/program"
)

func main() {
	os.Exit(program.Run(run))
}

type Command struct {
	Dir string
}

func run(ctx context.Context, args ...string) error {
	cmd := &command.Command{}
	cli := cli.New("bud", "bud web framework")
	cli.Flag("chdir", "change the working dir").Short('C').String(&cmd.Dir).Default(".")
	cli.Command("new", "create a new bud project").Mount(new.New(log.Default()))
	cli.Command("serve", "serve your app").Mount(serve.New(cmd))
	return cli.Parse(ctx, args...)
}

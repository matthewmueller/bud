package main

import (
	"context"
	"os"

	"github.com/livebud/buddy/cli"
	"github.com/livebud/buddy/commands"
	"github.com/livebud/buddy/program"
	"hellobud.com/internal/di"
)

func main() {
	os.Exit(program.Run(run))
}

func run(ctx context.Context, args ...string) error {
	cli := cli.New(`helloapp`, `helloapp app`)
	cli.Mount(commands.New(di.New()))
	return cli.Parse(ctx, args...)
}

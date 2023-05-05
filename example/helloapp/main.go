package main

import (
	"context"
	"os"

	"app.com/internal/di"
	"github.com/livebud/buddy/cli"
	"github.com/livebud/buddy/command/serve"
	"github.com/livebud/buddy/program"
)

func main() {
	os.Exit(program.Run(run))
}

func run(ctx context.Context, args ...string) error {
	cli := cli.New(`helloapp`, `helloapp app`)
	injector := di.New()
	cli.Mount(serve.New(injector))
	return cli.Parse(ctx, args...)
}

package main

import (
	"context"

	"github.com/livebud/buddy/cli"
	"github.com/livebud/buddy/example/blog/internal/command/routes"
	"github.com/livebud/buddy/example/blog/internal/command/serve"
	"github.com/livebud/buddy/program"
)

func main() {
	program.Run(run)
}

func run(ctx context.Context, args ...string) error {
	cmd := cli.New("blog", "blog cli")
	cmd.Mount(serve.New())
	cmd.Mount(routes.New())
	return cmd.Parse(ctx, args...)
}

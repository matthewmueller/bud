package app

import (
	"context"
	"os"

	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/logger"
)

func Load(in di.Injector) (*cli.CLI, error) {
	return di.Load[*cli.CLI](in)
}

func Parse(ctx context.Context, in di.Injector, args ...string) error {
	cli, err := Load(in)
	if err != nil {
		return err
	}
	return cli.Parse(ctx, args...)
}

func Run(in di.Injector) int {
	ctx := context.Background()
	if err := Parse(ctx, in, os.Args[1:]...); err != nil {
		logger.Error(err.Error())
		return 1
	}
	return 0
}

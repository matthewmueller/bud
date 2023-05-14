package program

import (
	"context"
	"errors"
	"os"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/internal/cli"

	"github.com/matthewmueller/bud/internal/signals"
	"github.com/matthewmueller/bud/internal/stacktrace"
	"github.com/matthewmueller/bud/logger"
)

type Provider = func(in di.Injector)

func New(injector di.Injector) *Program {
	return &Program{injector}
}

type Program struct {
	injector di.Injector
}

func (p *Program) Run(ctx context.Context, args ...string) error {
	ctx = signals.Trap(ctx, os.Interrupt)
	cli, err := di.Load[cli.Parser](p.injector)
	if err != nil {
		return err
	}
	return cli.Parse(ctx, args...)
}

func Run(fn func(ctx context.Context, in di.Injector, args ...string) error) int {
	ctx := signals.Trap(context.Background(), os.Interrupt)
	log := logger.Default()
	in := di.New()
	if err := fn(ctx, in, os.Args[1:]...); err != nil && !errors.Is(err, context.Canceled) {
		log.Field("source", stacktrace.Source(2)).Error(err.Error())
		return 1
	}
	return 0
}

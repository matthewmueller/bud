package program

import (
	"context"
	"errors"
	"os"

	"github.com/livebud/buddy/cli"
	"github.com/livebud/buddy/di"

	"github.com/livebud/buddy/internal/signals"
	"github.com/livebud/buddy/internal/stacktrace"
	"github.com/livebud/buddy/log"
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
	log := log.Default()
	in := di.New()
	if err := fn(ctx, in, os.Args[1:]...); err != nil && !errors.Is(err, context.Canceled) {
		log.Field("source", stacktrace.Source(2)).Error(err.Error())
		return 1
	}
	return 0
}

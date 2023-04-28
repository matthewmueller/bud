package app

import (
	"context"
	"errors"
	"os"

	"github.com/livebud/buddy/internal/stacktrace"

	"github.com/livebud/buddy/internal/signals"
	"github.com/livebud/buddy/log"
)

func Run(fn func(ctx context.Context, args ...string) error) int {
	ctx := signals.Trap(context.Background(), os.Interrupt)
	log := log.Default()
	if err := fn(ctx, os.Args[1:]...); err != nil && !errors.Is(err, context.Canceled) {
		log.Field("source", stacktrace.Source(2)).Error(err.Error())
		return 1
	}
	return 0
}

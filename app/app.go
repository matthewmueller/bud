package app

import (
	"context"
	"fmt"
	"os"

	"github.com/matthewmueller/bud/db/sqlite"

	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/db"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/logger"
)

func Provider(in di.Injector) {
	cli.Provider(in)
	logger.Provider(in)
	di.Provide[db.DB](in, func(in di.Injector) (db.DB, error) {
		return sqlite.Open("bud.db")
	})
}

func Load(in di.Injector) (cli.Command, error) {
	return di.Load[cli.Command](in)
}

func Parse(ctx context.Context, in di.Injector, args ...string) error {
	cmd, err := Load(in)
	if err != nil {
		return err
	}
	cli, ok := cmd.(*cli.CLI)
	if !ok {
		return fmt.Errorf("TODO: not a cli")
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

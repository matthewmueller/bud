package serve

import (
	"context"
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
	"github.com/livebud/buddy/cli"
	"github.com/livebud/buddy/example/blog/internal/web"
	"github.com/livebud/buddy/log"
)

func New() *Command {
	return &Command{":8080", "info"}
}

type Command struct {
	Address string
	Log     string
}

func (c *Command) Mount(cmd cli.Command) {
	cmd = cmd.Command("serve", "serve web requests")
	cmd.Flag("log", "log level").String(&c.Log).Default(c.Log)
	cmd.Flag("listen", "listen on address").String(&c.Address).Default(c.Address)
	cmd.Run(c.Serve)
}

func (c *Command) Serve(ctx context.Context) error {
	// env, err := env.Load()
	// if err != nil {
	// 	return err
	// }
	log, err := log.Load(c.Log)
	if err != nil {
		return err
	}
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return err
	}
	defer db.Close()
	server := web.New(log, db)
	log.Infof("listening on %s", c.Address)
	return server.Listen(ctx, c.Address)
}

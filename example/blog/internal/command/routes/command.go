package routes

import (
	"context"
	"fmt"

	"github.com/livebud/buddy/cli"
)

func New() *Command {
	return &Command{}
}

type Command struct {
	Port int
}

func (c *Command) Mount(cmd cli.Command) {
	cmd = cmd.Command("routes", "list the web routes")
	cmd.Run(c.Routes)
}

func (c *Command) Routes(ctx context.Context) error {
	// log, err := log.Load(c.Log)
	// if err != nil {
	// 	return err
	// }
	// db, err := sql.Open("sqlite", ":memory:")
	// if err != nil {
	// 	return err
	// }
	// defer db.Close()
	// log.New()
	fmt.Println("routes!", c.Port)
	return nil
}

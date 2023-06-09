package command

import (
	"os"

	"github.com/matthewmueller/bud/internal/shell"
	"github.com/matthewmueller/bud/log"
)

type Command struct {
	Dir string
	Log string
}

func (c *Command) Logger() (log.Log, error) {
	return log.Default(), nil
}

func (c *Command) Shell() shell.Runner {
	return &shell.Command{
		Dir: c.Dir,
		Env: map[string]string{
			"HOME":   os.Getenv("HOME"),
			"PATH":   os.Getenv("PATH"),
			"GOPATH": os.Getenv("GOPATH"),
		},
		Stderr: os.Stderr,
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
	}
}

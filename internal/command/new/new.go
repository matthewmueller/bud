package new

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/internal/gotext"
	"github.com/matthewmueller/bud/internal/shell"
	"github.com/matthewmueller/bud/internal/task"
	"github.com/matthewmueller/bud/internal/task/filemaker"
	"github.com/matthewmueller/bud/log"
)

func New(log log.Log) *Command {
	return &Command{log: log}
}

type Command struct {
	log   log.Log
	Dir   string
	Dev   bool
	Force bool
}

func (c *Command) Mount(cmd cli.Command) {
	// new <dir>
	cmd.Arg("dir").String(&c.Dir)
	cmd.Flag("dev", "enable development mode").Bool(&c.Dev).Default(false)
	cmd.Flag("force", "force to overwrite existing files").Bool(&c.Force).Default(false)
	cmd.Run(c.New)

	{ // new:undo <dir>
		cmd = cmd.Command("undo", "undo the new command")
		cmd.Arg("dir").String(&c.Dir)
		cmd.Flag("dev", "enable development mode").Bool(&c.Dev).Default(false)
		cmd.Flag("force", "force to overwrite existing files").Bool(&c.Force).Default(false)
		cmd.Run(c.Undo)
	}
}

// New command
func (c *Command) New(ctx context.Context) error {
	cmd := shell.New(c.Dir)
	cmd.Env["HOME"] = os.Getenv("HOME")
	cmd.Env["PATH"] = os.Getenv("PATH")
	cmd.Env["GOPATH"] = os.Getenv("GOPATH")
	return task.Do(ctx,
		&maingoScaffold{c.log, c.Dir, c.Force, filepath.Base(c.Dir), fmt.Sprintf("%s app", filepath.Base(c.Dir))},
		&gomodScaffold{c.log, c.Force, c.Dir, "app.com", "1.20", "../.."},
		&gomodTidy{c.log, cmd},
	)
}

// Undo command
func (c *Command) Undo(ctx context.Context) error {
	return task.Undo(ctx)
}

type maingoScaffold struct {
	log   log.Log
	dir   string
	force bool
	Name  string
	Desc  string
}

var _ task.Interface = (*maingoScaffold)(nil)

//go:embed main.gotext
var maingoTemplate string

var maingoGenerator = gotext.MustParse("main.go", maingoTemplate)

func (s *maingoScaffold) Do(ctx context.Context) error {
	data, err := maingoGenerator.Generate(s)
	if err != nil {
		return err
	}
	path := filepath.Join(s.dir, "main.go")
	task := filemaker.New(s.log, path, data, s.force)
	return task.Do(ctx)
}

func (s *maingoScaffold) Undo(ctx context.Context) error {
	data, err := maingoGenerator.Generate(s)
	if err != nil {
		return err
	}
	path := filepath.Join(s.dir, "main.go")
	task := filemaker.New(s.log, path, data, s.force)
	return task.Undo(ctx)
}

//go:embed gomod.gotext
var gomodTemplate string

var gomodGenerator = gotext.MustParse("go.mod", gomodTemplate)

type gomodScaffold struct {
	log        log.Log
	force      bool
	Dir        string
	Module     string
	GoVersion  string
	ReplaceBud string
}

var _ task.Interface = (*gomodScaffold)(nil)

func (s *gomodScaffold) Do(ctx context.Context) error {
	data, err := gomodGenerator.Generate(s)
	if err != nil {
		return err
	}
	path := filepath.Join(s.Dir, "go.mod")
	task := filemaker.New(s.log, path, data, s.force)
	return task.Do(ctx)
}

func (s *gomodScaffold) Undo(ctx context.Context) error {
	data, err := gomodGenerator.Generate(s)
	if err != nil {
		return err
	}
	path := filepath.Join(s.Dir, "go.mod")
	task := filemaker.New(s.log, path, data, s.force)
	return task.Undo(ctx)
}

type gomodTidy struct {
	log log.Log
	cmd shell.Runner
}

var _ task.Interface = (*gomodTidy)(nil)

func (s *gomodTidy) Do(ctx context.Context) error {
	s.log.Info("running go mod tidy")
	return s.cmd.Run(ctx, "go", "mod", "tidy")
}

func (s *gomodTidy) Undo(ctx context.Context) error {
	return nil
}

package shell

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type Runner interface {
	Run(ctx context.Context, name string, args ...string) error
}

func New(dir string) *Command {
	return &Command{
		Dir:    dir,                     // Default to cwd
		Env:    make(map[string]string), // Default to empty
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
}

type Command struct {
	Dir    string
	Env    map[string]string
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

func (c *Command) clone() *Command {
	env := make(map[string]string, len(c.Env))
	for k, v := range c.Env {
		env[k] = v
	}
	return &Command{
		Dir:    c.Dir,
		Env:    env,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
		Stdin:  c.Stdin,
	}
}

func (c *Command) New(cmd *Command) *Command {
	clone := c.clone()
	clone.Dir = cmd.Dir
	for k, v := range cmd.Env {
		clone.Env[k] = v
	}
	clone.Stdout = cmd.Stdout
	clone.Stderr = cmd.Stderr
	clone.Stdin = cmd.Stdin
	return clone
}

func (c *Command) Run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = c.Dir
	cmd.Env = make([]string, 0, len(c.Env))
	for k, v := range c.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.Stdin = c.Stdin
	return cmd.Run()
}

// Log runner logs the command before running it
// func Log(log log.Log, r Runner) Runner {
// 	return &logRunner{log, r}
// }

// type logRunner struct {
// 	log log.Log
// 	Runner
// }

// func (r *logRunner) Run(ctx context.Context, name string, args ...string) error {
// 	for i, arg := range args {
// 		if strings.Contains(arg, " ") {
// 			args[i] = strconv.Quote(arg)
// 		}
// 	}
// 	r.log.Infof("running %s %s", name, strings.Join(args, " "))
// 	return r.Runner.Run(ctx, name, args...)
// }

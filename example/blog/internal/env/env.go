package env

import "github.com/livebud/buddy/env"

func Load() (*Env, error) {
	return env.Load[Env]()
}

type Env struct {
	Hi       string `env:"HI" default:"hello"`
	Database *Database
}

type Database struct {
	URL string `env:"DATABASE_URL" default:"postgres://localhost:5432/blog?sslmode=disable"`
}

func (e *Env) Development() error {
	return nil
}

func (e *Env) Test() error {
	return nil
}

func (e *Env) Preview() error {
	return nil
}

func (e *Env) Production() error {
	return nil
}

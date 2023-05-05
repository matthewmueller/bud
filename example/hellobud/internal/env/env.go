package env

import "github.com/livebud/buddy/env"

type Env struct {
	Log string `env:"LOG" default:"info"`
}

func Load() (*Env, error) {
	return env.Load[Env]()
}

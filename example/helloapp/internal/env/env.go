package env

import (
	"os"

	"github.com/livebud/buddy/env"
)

type Env struct {
	DatabaseURL string `env:"DATABASE_URL" default:"sqlite://development.db"`
	Log         string `env:"LOG" default:"info"`
}

var envs = env.Map[Env]{
	env.Development: development,
}

func Load() (*Env, error) {
	return envs.Load(os.Getenv("BUD_ENV"))
}

func Development() (*Env, error) {
	return envs.Load("development")
}

func development(e *Env) error {
	// e.DatabaseURL = "sqlite://development.db"
	return nil
}

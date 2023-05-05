package di

import (
	"github.com/livebud/buddy/di"
	"github.com/livebud/buddy/log"
	"github.com/livebud/buddy/web"
	"hellobud.com/internal/env"
	"hellobud.com/internal/router"
)

func New() di.Injector {
	in := di.New()
	di.Provide(in, provideEnv)
	di.Provide(in, provideLog)
	di.Provide(in, provideRouter)
	return in
}

func provideEnv(in di.Injector) (*env.Env, error) {
	return env.Load()
}

func provideLog(in di.Injector) (log.Log, error) {
	env, err := di.Load[*env.Env](in)
	if err != nil {
		return nil, err
	}
	return log.Parse(env.Log)
}

func provideRouter(in di.Injector) (web.Router, error) {
	env, err := di.Load[*env.Env](in)
	if err != nil {
		return nil, err
	}
	log, err := di.Load[log.Log](in)
	if err != nil {
		return nil, err
	}
	return router.New(env, log), nil
}

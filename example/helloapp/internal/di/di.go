package di

import (
	"database/sql"

	"app.com/internal/env"
	_ "github.com/glebarez/go-sqlite"
	"github.com/livebud/buddy/di"
	"github.com/livebud/buddy/log"
)

func New() di.Injector {
	in := di.New()
	di.Provide(in, provideEnv)
	di.Provide(in, provideLog)
	// di.Provide(in, provideServer)
	di.Provide(in, provideDatabase)
	return in
}

func provideEnv(i di.Injector) (*env.Env, error) {
	return env.Load()
}

func provideLog(i di.Injector) (log.Log, error) {
	env, err := di.Load[*env.Env](i)
	if err != nil {
		return nil, err
	}
	return log.Parse(env.Log)
}

func provideDatabase(i di.Injector) (*sql.DB, error) {
	env, err := di.Load[*env.Env](i)
	if err != nil {
		return nil, err
	}
	return sql.Open("sqlite", env.DatabaseURL)
}

// func provideServer(i di.Injector) (web.Server, error) {
// 	env, err := di.Load[*env.Env](i)
// 	if err != nil {
// 		return nil, err
// 	}
// 	log, err := di.Load[log.Log](i)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return server.New(env, log), nil
// }

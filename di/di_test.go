package di_test

import (
	"testing"

	"github.com/livebud/buddy/di"
	"github.com/matryer/is"
)

type Env struct {
}

func loadEnv(in di.Injector) (*Env, error) {
	return &Env{}, nil
}

type Log struct {
	env *Env
	V   string
}

func loadLog(in di.Injector) (*Log, error) {
	env, err := di.Load[*Env](in)
	if err != nil {
		return nil, err
	}
	return &Log{env: env, V: "hello"}, nil
}

func TestDI(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Provide(in, loadEnv)
	di.Provide(in, loadLog)
	log, err := di.Load[*Log](in)
	is.NoErr(err)
	is.Equal(log.V, "hello")
}

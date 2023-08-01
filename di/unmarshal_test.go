package di_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/bud/di"
)

type Context struct {
	Env *Env
	Log *Log
}

func TestUnmarshal(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Provide(in, loadEnv)
	di.Provide(in, loadLog)
	ctx, err := di.Unmarshal[Context](in)
	is.NoErr(err)
	is.Equal(ctx.Env.e, "hi")
	is.Equal(ctx.Log.V, "hello")
}

func TestUnmarshalPointer(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Provide(in, loadEnv)
	di.Provide(in, loadLog)
	ctx, err := di.Unmarshal[*Context](in)
	is.NoErr(err)
	is.Equal(ctx.Env.e, "hi")
	is.Equal(ctx.Log.V, "hello")
}

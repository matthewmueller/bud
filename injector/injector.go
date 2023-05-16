package injector

import (
	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/command/develop"
	"github.com/matthewmueller/bud/command/serve"
	"github.com/matthewmueller/bud/db"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/logger"
	"github.com/matthewmueller/bud/middleware"
	"github.com/matthewmueller/bud/web"
	"github.com/matthewmueller/bud/web/router"
)

func Provider(in di.Injector) {
	cli.Provider(in)
	logger.Provider(in)
	router.Provider(in)
	db.Provider(in)
	web.Provider(in)
	middleware.Provider(in)
	di.Provide[*serve.Command](in, serve.Provide)
	di.Register[*cli.CLI](in, serve.Register)
	di.Provide[*develop.Command](in, develop.Provide)
	di.Register[*cli.CLI](in, develop.Register)
}

func New() di.Injector {
	in := di.New()
	Provider(in)
	return in
}

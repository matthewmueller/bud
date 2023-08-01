package web

import "github.com/matthewmueller/bud/di"

func Provider(in di.Injector) {
	di.Provide[Handler](in, provideHandler)
	di.Provide[*Server](in, provideServer)
}

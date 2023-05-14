package logger

import "github.com/matthewmueller/bud/di"

func Provider(in di.Injector) {
	di.Provide[Level](in, provideLevel)
	di.Provide[Log](in, provideLog)
}

func provideLevel(in di.Injector) (Level, error) {
	return InfoLevel, nil
}

func provideLog(in di.Injector) (Log, error) {
	return Default(), nil
}

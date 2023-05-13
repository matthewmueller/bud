package logger_test

import (
	"errors"
	"os"

	"github.com/matthewmueller/bud/internal/color"
	"github.com/matthewmueller/bud/logger"
)

func ExampleConsole() {
	logger.Infof("hello %s", "mars")
	logger.Noticef("hello %s", "mars")
	logger.Warnf("hello %s", "mars")
	logger.Errorf("hello %s", "mars")
	logger.Error(errors.New("one"), "two", "three")
	logger := logger.New(logger.Console(color.Ignore(), os.Stdout))
	logger.Debug("hello", "world")
	logger.Error(errors.New("one"), 4, "three")
	// Output:
	// | hello world
	// | one 4 three
}

package log_test

import (
	"errors"
	"os"

	"github.com/matthewmueller/bud/internal/color"
	"github.com/matthewmueller/bud/log"
)

func ExampleConsole() {
	log.Infof("hello %s", "mars")
	log.Noticef("hello %s", "mars")
	log.Warnf("hello %s", "mars")
	log.Errorf("hello %s", "mars")
	log.Error(errors.New("one"), "two", "three")
	logger := log.New(log.Console(color.Ignore(), os.Stdout))
	logger.Debug("hello", "world")
	logger.Error(errors.New("one"), 4, "three")
	// Output:
	// | hello world
	// | one 4 three
}

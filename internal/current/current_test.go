package current_test

import (
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/bud/internal/current"
)

func TestDir(t *testing.T) {
	is := is.New(t)
	dirname, err := current.Directory()
	is.NoErr(err)
	is.Equal(filepath.Base(dirname), "current")
}

func TestFile(t *testing.T) {
	is := is.New(t)
	filename, err := current.Filename()
	is.NoErr(err)
	is.Equal(filepath.Base(filename), "current_test.go")
}

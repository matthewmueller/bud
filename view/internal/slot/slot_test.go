package slot_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/bud/view/internal/slot"
	"golang.org/x/sync/errgroup"
)

func TestList(t *testing.T) {
	is := is.New(t)
	a := slot.New()
	b := a.New()
	c := b.New()
	d := c.New()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		defer a.Close()
		buf := new(bytes.Buffer)
		buf.WriteByte('a')
		_, err := buf.WriteTo(a)
		return err
	})
	eg.Go(func() error {
		defer b.Close()
		buf := new(bytes.Buffer)
		data, err := io.ReadAll(b)
		if err != nil {
			return err
		}
		buf.Write(data)
		buf.WriteByte('b')
		_, err = buf.WriteTo(b)
		return err
	})
	eg.Go(func() error {
		defer c.Close()
		buf := new(bytes.Buffer)
		data, err := io.ReadAll(c)
		if err != nil {
			return err
		}
		buf.Write(data)
		buf.WriteByte('c')
		_, err = buf.WriteTo(c)
		return err
	})
	is.NoErr(eg.Wait())
	result, err := io.ReadAll(d)
	is.NoErr(err)
	is.Equal(string(result), "abc")
}

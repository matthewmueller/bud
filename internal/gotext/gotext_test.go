package gotext_test

import (
	"testing"

	"github.com/matthewmueller/bud/internal/gotext"
	"github.com/matryer/is"
)

func TestGenerateGoFile(t *testing.T) {
	is := is.New(t)
	template := `package main

	func main()  {
		  println("{{ .name }}")
}`
	expect := `package main

func main() {
	println("jason")
}
`
	generator := gotext.MustParse("test.gotext", template)
	b, err := generator.Generate(map[string]string{"name": "jason"})
	is.NoErr(err)
	is.Equal(string(b), expect)
}

func TestGenerateFreeText(t *testing.T) {
	is := is.New(t)
	template := `Hi {{ .name }}`
	expect := `Hi Kim`
	generator := gotext.MustParse("test.gotext", template)
	b, err := generator.Generate(map[string]string{"name": "Kim"})
	is.NoErr(err)
	is.Equal(string(b), expect)
}

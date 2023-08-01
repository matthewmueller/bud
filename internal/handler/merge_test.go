package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/bud/internal/handler"
)

func TestMerge(t *testing.T) {
	is := is.New(t)
	merged := handler.Merge(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("a", "aa")
			w.Write([]byte("a"))
		}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("b", "bb")
			w.Write([]byte("b"))
		}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("b", "cc")
			w.Write([]byte("c"))
		}),
	)
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	merged.ServeHTTP(w, r)
	res := w.Result()
	is.Equal(res.StatusCode, http.StatusOK)
	is.Equal(res.Header.Get("a"), "aa")
	is.Equal(res.Header.Get("b"), "cc")
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "abc")
}

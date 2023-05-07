package cookie

import (
	"net/http"
	"time"
)

// New default cookie
func New(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}
}

// Header is an interface for getting and setting cookies
type Header interface {
	Get(r *http.Request, name string) (*http.Cookie, error)
	Set(w http.ResponseWriter, cookie *http.Cookie) error
}

func Default() Header {
	return &defaultHeader{}
}

type defaultHeader struct{}

func (defaultHeader) Get(r *http.Request, name string) (*http.Cookie, error) {
	return r.Cookie(name)
}

func (defaultHeader) Set(w http.ResponseWriter, cookie *http.Cookie) error {
	http.SetCookie(w, cookie)
	return nil
}

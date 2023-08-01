package request

import (
	"net/http"

	"github.com/timewasted/go-accept-headers"
)

// func Wrap[In any](r *http.Request) (*Request[In], error) {
// 	var in In
// 	if err := Unmarshal(r, &in); err != nil {
// 		return nil, err
// 	}
// 	return &Request[In]{
// 		r:  r,
// 		In: in,
// 	}, nil
// }

// Request struct
// type Request[In any] struct {
// 	r  *http.Request
// 	In any
// }

// // Unmarshal the request body or parameters
// func (c *Request[In]) Unmarshal(r *http.Request, in interface{}) error {
// 	return Unmarshal(r, in)
// }

// Accepts a type
func Accepts(r *http.Request) Acceptable {
	return Acceptable(accept.Parse(r.Header.Get("Accept")))
}

// Acceptable types
type Acceptable accept.AcceptSlice

// Accepts checks if the content type is acceptable
func (as Acceptable) Accepts(ctype string) bool {
	return accept.AcceptSlice(as).Accepts(ctype)
}

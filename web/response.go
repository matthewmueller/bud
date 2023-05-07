package web

import "net/http"

type Response[Body any] interface {
	// http.ResponseWriter
	Status(status int) Response[Body]
	HTML(html string) error
	JSON(response Body) error
	Render(props Body) error
}

type response[Out any] struct {
	rw http.ResponseWriter
}

var _ Response[any] = &response[any]{}

func (r *response[Out]) Status(status int) Response[Out] {
	r.rw.WriteHeader(status)
	return r
}

func (r *response[Out]) HTML(html string) error {
	r.rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := r.rw.Write([]byte(html))
	return err
}

func (r *response[Out]) JSON(out Out) error {
	return nil
}

func (r *response[Out]) Render(out Out) error {
	return nil
}

// func (r *response) HTML(html string) error {
// 	_, err := r.rw.Write([]byte(html))
// 	return err
// }

// func (r *response) JSON(out Out) error {
// 	return nil
// }

// func (r *response) Render(out Out) error {
// 	return nil
// }

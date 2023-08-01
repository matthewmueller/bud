package web

import (
	"net/http"
	"reflect"
	"runtime"
)

func Action[In, Out any](fn func(req *Request[In], res Response[Out]) error) http.Handler {
	info := runtime.FuncForPC(reflect.ValueOf(fn).Pointer())
	return &action[In, Out]{fn, info}
}

type action[In, Out any] struct {
	fn   func(req *Request[In], res Response[Out]) error
	info *runtime.Func
}

func (a *action[In, Out]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// var in In
	request := &Request[In]{
		// r,
		// di.New(),
		// in,
	}
	response := &response[Out]{
		w,
	}
	if err := a.fn(request, response); err != nil {
		// TODO improve error handling
		http.Error(w, err.Error(), 500)
	}
}

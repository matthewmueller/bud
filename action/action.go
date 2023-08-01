package action

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/matthewmueller/bud/internal/reflector"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/view"
)

func New[Context, In, Out any](fn func(ctx Context, in In) (Out, error)) *Function[Context, In, Out] {
	return &Function[Context, In, Out]{fn, viewKey(fn)}
}

type Function[Context, In, Out any] struct {
	fn      func(ctx Context, in In) (Out, error)
	viewKey string
}

func Layout[Context, In, Out any](fn func(ctx Context, in struct{}) (struct{}, error)) *Function[Context, struct{}, struct{}] {
	return New(fn)
}

func (f *Function[Context, In, Out]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	injector, err := di.InjectorFrom(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var in In
	context, err := di.Unmarshal[Context](injector)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	out, err := f.fn(context, in)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	pages, err := di.Load[view.Pages](injector)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	actions, err := di.Load[view.Action](injector)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println("got actions", actions)
	page, ok := pages[f.viewKey]
	if !ok {
		body, err := json.Marshal(out)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
	pageData, err := json.MarshalIndent(page, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	fmt.Println("got page", string(pageData))
	// fmt.Println("got pages", pages[f.viewKey])
	// viewMap, err := di.Load[view.Map](injector)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// view, ok := viewMap[f.viewKey]
	// if !ok {
	// 	body, err := json.Marshal(out)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Write([]byte(err.Error()))
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(body)
	// 	return
	// }
	// if err := view.Render(ctx, w, out); err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
}

func viewKey[Context, In, Out any](fn func(ctx Context, in In) (Out, error)) string {
	info, err := reflector.Func(fn)
	if err != nil {
		return ""
	}
	name := info.Name()
	idx := strings.LastIndex(name, "controller")
	if idx == -1 {
		return ""
	}
	name = name[idx+10:]
	if name[0] != '/' && name[0] != '.' {
		return ""
	}
	name = name[1:]
	parts := strings.SplitN(name, ".", 3)
	switch len(parts) {
	case 2:
		return strings.ToLower(parts[1])
	case 3:
		return parts[0] + "/" + strings.ToLower(parts[2])
	default:
		return ""
	}
}

package di

import (
	"fmt"
	"sync"

	"github.com/matthewmueller/bud/internal/reflector"
)

func New() Injector {
	return &injector{
		fns:   make(map[string]any),
		cache: make(map[string]any),
	}
}

type Injector interface {
	setProvider(name string, provider any)
	getProvider(name string) (provider any, ok bool)
	setCache(name string, dep any)
	getCache(name string) (dep any, ok bool)
}

type injector struct {
	mu    sync.RWMutex
	fns   map[string]any
	cache map[string]any
}

var _ Injector = (*injector)(nil)

func (in *injector) setProvider(name string, provider any) {
	in.mu.Lock()
	defer in.mu.Unlock()
	in.fns[name] = provider
}

func (in *injector) getProvider(name string) (provider any, ok bool) {
	in.mu.RLock()
	defer in.mu.RUnlock()
	fn, ok := in.fns[name]
	return fn, ok
}

func (in *injector) setCache(name string, dep any) {
	in.mu.Lock()
	defer in.mu.Unlock()
	in.cache[name] = dep
}

func (in *injector) getCache(name string) (dep any, ok bool) {
	in.mu.RLock()
	defer in.mu.RUnlock()
	dep, ok = in.cache[name]
	return dep, ok
}

func Provide[Dep any](in Injector, fn func(in Injector) (d Dep, err error)) error {
	var dep Dep
	name, err := reflector.Name(dep)
	if err != nil {
		return err
	}
	in.setProvider(name, fn)
	return nil
}

func Load[Dep any](in Injector) (dep Dep, err error) {
	name, err := reflector.Name(dep)
	if err != nil {
		return dep, err
	}
	if dep, ok := in.getCache(name); ok {
		return dep.(Dep), nil
	}
	v, ok := in.getProvider(name)
	if !ok {
		return dep, fmt.Errorf("di: no provider for %s", name)
	}
	fn, ok := v.(func(in Injector) (Dep, error))
	if !ok {
		return dep, fmt.Errorf("invalid provider for %s", name)
	}
	d, err := fn(in)
	if err != nil {
		return dep, err
	}
	in.setCache(name, d)
	return d, nil
}

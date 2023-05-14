package di

import (
	"fmt"
	"sync"

	"github.com/matthewmueller/bud/internal/reflector"
)

func New() Injector {
	return &injector{
		fns:        make(map[string]any),
		cache:      make(map[string]any),
		registered: make(map[string][]any),
	}
}

type Injector interface {
	setProvider(name string, provider any) error
	addProvider(name string, provider any) error
	register(name string, provider any) error
	registrants(name string) (registrants []any, ok bool)
	getProvider(name string) (provider any, ok bool)
	listProviders() []string
	setCache(name string, dep any)
	getCache(name string) (dep any, ok bool)
}

type injector struct {
	mu         sync.RWMutex
	fns        map[string]any
	cache      map[string]any
	registered map[string][]any
}

var _ Injector = (*injector)(nil)

func (in *injector) setProvider(name string, provider any) error {
	in.mu.Lock()
	defer in.mu.Unlock()
	in.fns[name] = provider
	return nil
}

func (in *injector) addProvider(name string, provider any) error {
	in.mu.Lock()
	defer in.mu.Unlock()
	if _, ok := in.fns[name]; !ok {
		in.fns[name] = []any{}
	}
	fns, ok := in.fns[name].([]any)
	if !ok {
		return fmt.Errorf("unable to add provider to %s", name)
	}
	in.fns[name] = append(fns, provider)
	return nil
}

func (in *injector) register(name string, provider any) error {
	in.mu.Lock()
	defer in.mu.Unlock()
	if _, ok := in.registered[name]; !ok {
		in.registered[name] = []any{}
	}
	in.registered[name] = append(in.registered[name], provider)
	return nil
}

func (in *injector) registrants(name string) (registrants []any, ok bool) {
	in.mu.RLock()
	defer in.mu.RUnlock()
	registrants, ok = in.registered[name]
	return registrants, ok
}

func (in *injector) getProvider(name string) (provider any, ok bool) {
	in.mu.RLock()
	defer in.mu.RUnlock()
	if fn, ok := in.fns[name]; ok {
		return fn, ok
	}
	return nil, false
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

func (in *injector) listProviders() (providers []string) {
	in.mu.RLock()
	defer in.mu.RUnlock()
	for name := range in.fns {
		providers = append(providers, name)
	}
	return providers
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
	registrants, ok := in.registrants(name)
	if ok {
		for _, registrant := range registrants {
			fn, ok := registrant.(func(in Injector, dep Dep) error)
			if !ok {
				return dep, fmt.Errorf("invalid attachment for %s", name)
			}
			if err := fn(in, d); err != nil {
				return dep, err
			}
		}
	}
	in.setCache(name, d)
	return d, nil
}

func loadOne[Dep any](in Injector, name string, fn func(in Injector) (Dep, error)) (dep Dep, err error) {
	d, err := fn(in)
	if err != nil {
		return dep, err
	}
	in.setCache(name, d)
	return d, nil
}

// func loadAll[Dep any](in Injector, name string, fns []func(in Injector) (Dep, error)) (dep Dep, err error) {
// 	var deps []Dep
// 	for _, fn := range fns {
// 		d, err := fn(in)
// 		if err != nil {
// 			return dep, err
// 		}
// 		deps = append(deps, d)
// 	}
// 	in.setCache(name, deps)
// 	return deps, nil
// }

// func Add[Dep any](in Injector, fn func(in Injector) (d Dep, err error)) error {
// 	var dep []Dep
// 	name, err := reflector.Name(dep)
// 	if err != nil {
// 		return err
// 	}
// 	provider, ok := in.getProvider(name)
// 	if !ok {
// 		return in.setProvider(name, func(in Injector) ([]Dep, error) {
// 			dep, err := fn(in)
// 			return []Dep{dep}, err
// 		})
// 	}
// 	fns, ok := provider.(func(in Injector) ([]Dep, error))
// 	if !ok {
// 		return fmt.Errorf("invalid provider for %s", name)
// 	}
// 	fmt.Println("got it...", fns)
// 	// fns = append(fns, fn)
// 	return in.setProvider(name, fns)
// 	// return nil
// }

func Register[To any](in Injector, fn func(in Injector, to To) error) error {
	var to To
	name, err := reflector.Name(to)
	if err != nil {
		return err
	}
	return in.register(name, fn)
}

// func ProvideAs[Dep, As any](in Injector, fn func(in Injector) (d Dep, err error)) error {
// 	// dep := *new(Dep)
// 	// as := *new(As)
// 	// name, err := reflector.Name(dep)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// if !reflect.TypeOf(as).Implements(reflect.TypeOf(as)) {
// 	// 	return fmt.Errorf("di: %T does not implement %T", dep, as)
// 	// }
// 	// fmt.Println("got name!", name, as)
// 	// dep.(As)
// 	// var _ As = dep
// 	// if _, ok := dep.(As); !ok {
// 	// return fmt.Errorf("di: %s does not implement %s", reflector.Name(dep), reflector.Name(As))
// 	// }
// 	// name, err := reflector.Name(dep)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	return nil
// }

func Print(in Injector) string {
	providers := in.listProviders()
	return fmt.Sprintf("di: %d providers\n%s", len(providers), providers)
	// var s string
	// for name := range in.(*injector).fns {
	// 	s += name + "\n"
	// }
	// return s
}

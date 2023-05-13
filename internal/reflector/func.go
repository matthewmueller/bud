package reflector

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func Func[F any](fn F) (*FuncInfo, error) {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("reflector: expected a function, got %T", t)
	}
	ptr := reflect.ValueOf(fn).Pointer()
	info := runtime.FuncForPC(ptr)
	return &FuncInfo{info}, nil
}

type FuncInfo struct {
	info *runtime.Func
}

func (f *FuncInfo) Name() string {
	name := f.info.Name()
	name = strings.Replace(name, "-fm", "", 1)
	name = strings.Replace(name, "(", "", 1)
	name = strings.Replace(name, ")", "", 1)
	return name
}

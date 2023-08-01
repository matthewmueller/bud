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

func (f *FuncInfo) Path() string {
	name := f.info.Name()
	parts := strings.Split(name, ".")
	return strings.Join(parts[:len(parts)-1], ".")
}

// func ModulePath(skip int) (string, error) {
// 	pc, _, _, ok := runtime.Caller(skip)
// 	if !ok {
// 		return "", errors.New("unable to get the current filename")
// 	}
// 	info := runtime.FuncForPC(pc)
// 	name := info.Name()
// 	parts := strings.Split(name, ".")
// 	return strings.Join(parts[:len(parts)-1], "."), nil
// }

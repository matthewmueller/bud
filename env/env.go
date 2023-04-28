package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/matthewmueller/gotext"
)

var defaultLoader = &Loader{
	GetEnv: os.Getenv,
	Method: func() string {
		return gotext.Pascal(os.Getenv("BUD_ENV"))
	},
}

func Load[Env any]() (*Env, error) {
	var env Env
	if err := defaultLoader.Load(&env); err != nil {
		return nil, err
	}
	return &env, nil
}

type Loader struct {
	GetEnv func(string) string
	Method func() string
}

func (l *Loader) Load(env interface{}) error {
	t := reflect.TypeOf(env)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("env: expected a pointer to a struct, got %s", t)
	}
	v := reflect.ValueOf(env)
	return l.loadStruct(v.Elem())
}

func (l *Loader) loadStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		if fieldType.Anonymous {
			if err := l.loadStruct(fieldValue); err != nil {
				return err
			}
			continue
		}

		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			}
			fieldValue = fieldValue.Elem()
		}

		if fieldValue.Kind() == reflect.Struct {
			if err := l.loadStruct(fieldValue); err != nil {
				return err
			}
			continue
		}

		envKey := fieldType.Tag.Get("env")
		if envKey == "" {
			continue
		}

		envValue := os.Getenv(envKey)
		if envValue == "" {
			envValue = fieldType.Tag.Get("default")
		}
		if envValue == "" {
			return fmt.Errorf("env: missing required environment variable %q", envKey)
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(envValue)
		case reflect.Bool:
			if b, err := strconv.ParseBool(envValue); err == nil {
				fieldValue.SetBool(b)
			}
		case reflect.Int, reflect.Int64:
			if n, err := strconv.ParseInt(envValue, 10, 64); err == nil {
				fieldValue.SetInt(n)
			}
		default:
			return fmt.Errorf("env: unsupported field type: %s", fieldValue.Kind())
		}
	}

	// Call the environment function if it exists
	envFunc := v.Addr().MethodByName(l.Method())
	if envFunc.IsValid() {
		res := envFunc.Call(nil)
		if len(res) > 0 && !res[0].IsNil() {
			return res[0].Interface().(error)
		}
	}

	return nil
}

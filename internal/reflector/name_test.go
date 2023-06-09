package reflector_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/matthewmueller/bud/internal/reflector"
	"github.com/matryer/is"
)

// func TestName(t *testing.T) {
// 	reflector.Name()
// }

type customStruct struct {
	Field1 int
	Field2 string
}

type customInterface interface {
	Method1() int
}

type customFunc func() int

// // func Func() (int
// 	is.True(err != nil)
// 	is.True(errors.Is(err,reflector.ErrNoName))
// // is.Equal(name,er) {)
// // 	return 0
// is.True(err != nil)
// is.True(errors.Is(err,reflector.ErrNoName))
// is.Equal(name,ni)
// }

type customAlias = int

func TestName(t *testing.T) {
	is := is.New(t)
	name, err := reflector.Name(customStruct{})
	is.NoErr(err)
	is.Equal(name, "github.com/matthewmueller/bud/internal/reflector_test.customStruct")
	name, err = reflector.Name(&customStruct{})
	is.NoErr(err)
	is.Equal(name, "github.com/matthewmueller/bud/internal/reflector_test.*customStruct")
	name, err = reflector.Name(reflect.ValueOf((*customInterface)(nil)).Elem())
	is.NoErr(err)
	is.Equal(name, "reflect.Value")
	name, err = reflector.Name(customFunc(func() int { return 0 }))
	is.NoErr(err)
	is.Equal(name, "github.com/matthewmueller/bud/internal/reflector_test.customFunc")
}

func TestNoName(t *testing.T) {
	t.Skip("TODO: fix no builtins")
	is := is.New(t)
	name, err := reflector.Name(bool(false))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(int(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(int8(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(int16(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(int32(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(int64(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(uint(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(uint8(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(uint16(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(uint32(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(uint64(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(uintptr(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(float32(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(float64(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(complex64(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(complex128(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(string(""))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(byte(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(rune(0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(errors.New(""))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, ".*errorString")
	name, err = reflector.Name(new(int))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int")
	name, err = reflector.Name(new(int8))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int8")
	name, err = reflector.Name(new(int16))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int16")
	name, err = reflector.Name(new(int32))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int32")
	name, err = reflector.Name(new(int64))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int64")
	name, err = reflector.Name(new(uint))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint")
	name, err = reflector.Name(new(uint8))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint8")
	name, err = reflector.Name(new(uint16))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint16")
	name, err = reflector.Name(new(uint32))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint32")
	name, err = reflector.Name(new(uint64))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint64")
	name, err = reflector.Name(new(uintptr))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uintptr")
	name, err = reflector.Name(new(float32))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "float32")
	name, err = reflector.Name(new(float64))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "float64")
	name, err = reflector.Name(new(complex64))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "complex64")
	name, err = reflector.Name(new(complex128))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "complex128")
	name, err = reflector.Name(new(bool))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "bool")
	name, err = reflector.Name(new(string))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "string")
	name, err = reflector.Name(new(byte))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint8")
	name, err = reflector.Name(new(rune))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int32")
	name, err = reflector.Name(new(error))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "error")
	name, err = reflector.Name(new(interface{}))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "interface {}")
	name, err = reflector.Name([]int{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int")
	name, err = reflector.Name([]int8{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int8")
	name, err = reflector.Name([]int16{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int16")
	name, err = reflector.Name([]int32{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int32")
	name, err = reflector.Name([]int64{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int64")
	name, err = reflector.Name([]uint{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint")
	name, err = reflector.Name([]uint8{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint8")
	name, err = reflector.Name([]uint16{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint16")
	name, err = reflector.Name([]uint32{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint32")
	name, err = reflector.Name([]uint64{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint64")
	name, err = reflector.Name([]uintptr{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uintptr")
	name, err = reflector.Name([]float32{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "float32")
	name, err = reflector.Name([]float64{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "float64")
	name, err = reflector.Name([]complex64{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "complex64")
	name, err = reflector.Name([]bool{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "bool")
	name, err = reflector.Name([]string{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "string")
	name, err = reflector.Name([]byte{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "uint8")
	name, err = reflector.Name([]rune{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "int32")
	name, err = reflector.Name([]error{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "error")
	name, err = reflector.Name([]interface{}{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "interface {}")
	name, err = reflector.Name([1]int{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]int")
	name, err = reflector.Name([1]int8{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]int8")
	name, err = reflector.Name([1]int16{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]int16")
	name, err = reflector.Name([1]int32{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]int32")
	name, err = reflector.Name([1]int64{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]int64")
	name, err = reflector.Name([1]uint{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]uint")
	name, err = reflector.Name([1]uint8{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]uint8")
	name, err = reflector.Name([1]uint16{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]uint16")
	name, err = reflector.Name([1]uint32{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]uint32")
	name, err = reflector.Name([1]uint64{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]uint64")
	name, err = reflector.Name([1]uintptr{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]uintptr")
	name, err = reflector.Name([1]float32{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]float32")
	name, err = reflector.Name([1]float64{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]float64")
	name, err = reflector.Name([1]complex64{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]complex64")
	name, err = reflector.Name([1]complex128{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]complex128")
	name, err = reflector.Name([1]bool{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]bool")
	name, err = reflector.Name([1]string{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]string")
	name, err = reflector.Name([1]byte{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]uint8")
	name, err = reflector.Name([1]rune{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]int32")
	name, err = reflector.Name([1]error{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]error")
	name, err = reflector.Name([1]interface{}{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "1]interface {}")
	name, err = reflector.Name(struct{}{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, " {}")
	name, err = reflector.Name(struct {
		field1 int
		field2 string
	}{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, " { field1 int; field2 string }")
	name, err = reflector.Name(map[string]int{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "[string]int")
	name, err = reflector.Name(map[string]string{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "[string]string")
	name, err = reflector.Name(map[string]bool{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "[string]bool")
	name, err = reflector.Name(map[string]float32{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "[string]float32")
	name, err = reflector.Name(map[string]float64{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "[string]float64")
	name, err = reflector.Name(map[string]complex64{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "[string]complex64")
	name, err = reflector.Name(map[string]complex128{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "[string]complex128")
	name, err = reflector.Name(map[string]interface{}{})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "[string]interface {}")
	name, err = reflector.Name(int(42))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(float64(42.0))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(true)
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name("hello")
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name([]int{1, 2, 3})
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
	name, err = reflector.Name(customAlias(10))
	is.True(err != nil)
	is.True(errors.Is(err, reflector.ErrNoName))
	is.Equal(name, "")
}

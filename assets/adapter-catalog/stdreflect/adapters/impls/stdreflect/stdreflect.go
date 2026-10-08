package reflectdeps

import (
	"errors"
	"reflect"

	"{{.Module}}/sandbox/deps"
	reflectdeps "{{.Module}}/sandbox/deps/reflectdeps"
)

// Bind fills deps.Deps.Reflectdeps with the standard library's reflect.
func Bind(deps *deps.Deps) {
	deps.Reflectdeps = reflectdeps.Sandbox{
		NumIn:     numIn,
		NewIn:     newIn,
		Call:      call,
		NumField:  numField,
		FieldName: fieldName,
		FieldTag:  fieldTag,
		FieldKind: fieldKind,
		SetField:  setField,
	}
}

// funcType is the type of fn, nil when fn is not a function.
func funcType(fn any) reflect.Type {
	kind := reflect.TypeOf(fn)
	if kind == nil || kind.Kind() != reflect.Func {
		return nil
	}
	return kind
}

func numIn(fn any) int {
	kind := funcType(fn)
	if kind == nil {
		return -1
	}
	return kind.NumIn()
}

func newIn(fn any, index int) any {
	kind := funcType(fn)
	if kind == nil || index < 0 || index >= kind.NumIn() {
		return nil
	}
	param := kind.In(index)
	if param.Kind() == reflect.Pointer {
		return reflect.New(param.Elem()).Interface()
	}
	return reflect.Zero(param).Interface()
}

func call(fn any, args []any) []any {
	kind := funcType(fn)
	if kind == nil || kind.NumIn() != len(args) {
		return nil
	}

	in := make([]reflect.Value, len(args))
	for index, arg := range args {
		if arg == nil {
			in[index] = reflect.Zero(kind.In(index))
			continue
		}
		in[index] = reflect.ValueOf(arg)
	}

	out := []any{}
	for _, value := range reflect.ValueOf(fn).Call(in) {
		out = append(out, value.Interface())
	}
	return out
}

// structOf is the struct target points at, the zero Value when it is not a
// pointer to a struct.
func structOf(target any) reflect.Value {
	value := reflect.ValueOf(target)
	if value.Kind() != reflect.Pointer || value.IsNil() || value.Elem().Kind() != reflect.Struct {
		return reflect.Value{}
	}
	return value.Elem()
}

// fieldOf is the field at index of the struct target points at, false when
// there is none.
func fieldOf(target any, index int) (reflect.StructField, reflect.Value, bool) {
	value := structOf(target)
	if !value.IsValid() || index < 0 || index >= value.NumField() {
		return reflect.StructField{}, reflect.Value{}, false
	}
	return value.Type().Field(index), value.Field(index), true
}

func numField(target any) int {
	value := structOf(target)
	if !value.IsValid() {
		return -1
	}
	return value.NumField()
}

func fieldName(target any, index int) string {
	field, _, ok := fieldOf(target, index)
	if !ok {
		return ""
	}
	return field.Name
}

func fieldTag(target any, index int, key string) string {
	field, _, ok := fieldOf(target, index)
	if !ok {
		return ""
	}
	return field.Tag.Get(key)
}

func fieldKind(target any, index int) string {
	field, _, ok := fieldOf(target, index)
	if !ok {
		return ""
	}
	return field.Type.Kind().String()
}

func setField(target any, index int, value any) error {
	_, field, ok := fieldOf(target, index)
	if !ok {
		return errors.New("no such field")
	}
	if !field.CanSet() {
		return errors.New("field is not settable")
	}
	if value == nil {
		field.Set(reflect.Zero(field.Type()))
		return nil
	}
	given := reflect.ValueOf(value)
	if !given.Type().AssignableTo(field.Type()) {
		return errors.New("cannot assign " + given.Type().String() + " to " + field.Type().String())
	}
	field.Set(given)
	return nil
}

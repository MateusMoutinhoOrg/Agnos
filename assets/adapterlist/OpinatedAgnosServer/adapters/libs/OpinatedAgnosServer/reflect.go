package opinatedagnosserver

import (
	"errors"
	"reflect"
)

// The reflection the binder builds, fills and calls a route's own Entries
// with: every route declares a struct of its own, so the binder knows its
// type only at run time.

// funcType is the type of fn, nil when fn is not a function.
func funcType(fn any) reflect.Type {
	kind := reflect.TypeOf(fn)
	if kind == nil || kind.Kind() != reflect.Func {
		return nil
	}
	return kind
}

// newIn builds a fresh value for the parameter at index of fn: a pointer to a
// new zero value when that parameter is a pointer, the zero value otherwise,
// nil when fn is not a function or index is out of range.
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

// call calls fn with args, a nil arg passed as the zero value of its
// parameter, and returns what it returned; nil when fn is not a function or
// the argument count differs.
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

// structOf is the struct target points at, the zero Value when target is not
// a pointer to a struct.
func structOf(target any) reflect.Value {
	value := reflect.ValueOf(target)
	if value.Kind() != reflect.Pointer || value.IsNil() || value.Elem().Kind() != reflect.Struct {
		return reflect.Value{}
	}
	return value.Elem()
}

// fieldOf is the field at index of the struct target points at.
func fieldOf(target any, index int) (reflect.StructField, reflect.Value, bool) {
	value := structOf(target)
	if !value.IsValid() || index < 0 || index >= value.NumField() {
		return reflect.StructField{}, reflect.Value{}, false
	}
	return value.Type().Field(index), value.Field(index), true
}

// numField is how many fields the struct target points at has, -1 when target
// is not a pointer to a struct.
func numField(target any) int {
	value := structOf(target)
	if !value.IsValid() {
		return -1
	}
	return value.NumField()
}

// fieldName is the Go name of the field at index, "" when there is none.
func fieldName(target any, index int) string {
	field, _, ok := fieldOf(target, index)
	if !ok {
		return ""
	}
	return field.Name
}

// fieldTag is the value of key in the tag of the field at index, "" when the
// tag or the field is absent.
func fieldTag(target any, index int, key string) string {
	field, _, ok := fieldOf(target, index)
	if !ok {
		return ""
	}
	return field.Tag.Get(key)
}

// setField stores value in the field at index, refusing a field that is
// absent, not settable, or of a type value cannot be assigned to.
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

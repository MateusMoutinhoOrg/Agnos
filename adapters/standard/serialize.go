package standard

import (
	"encoding/json"
	"fmt"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/serializebles"
	"gopkg.in/yaml.v3"
)

// SerializeLibFactory returns the closure that fills deps.Deps.SerializeLib,
// providing the capability to create, parse, and serialize generic JSON/YAML structures.
func SerializeLibFactory(s *StandardAdapter) serializibles.Lib {
	return serializibles.Lib{
		CreateString: func(value string) serializibles.SerializibleObject {
			var v any = value
			return wrapValue(&v)
		},
		CreateInt: func(value int64) serializibles.SerializibleObject {
			var v any = value
			return wrapValue(&v)
		},
		CreateFloat: func(value float64) serializibles.SerializibleObject {
			var v any = value
			return wrapValue(&v)
		},
		CreateBool: func(value bool) serializibles.SerializibleObject {
			var v any = value
			return wrapValue(&v)
		},
		CreateNull: func() serializibles.SerializibleObject {
			var v any = nil
			return wrapValue(&v)
		},
		CreateObject: func() serializibles.SerializibleObject {
			var v any = make(map[string]any)
			return wrapValue(&v)
		},
		CreateArray: func() serializibles.SerializibleObject {
			var v any = make([]any, 0)
			return wrapValue(&v)
		},

		ParseJson: func(data string) (serializibles.SerializibleObject, error) {
			var v any
			if err := json.Unmarshal([]byte(data), &v); err != nil {
				return serializibles.SerializibleObject{}, err
			}
			return wrapValue(&v), nil
		},
		ParseYaml: func(data string) (serializibles.SerializibleObject, error) {
			var v any
			if err := yaml.Unmarshal([]byte(data), &v); err != nil {
				return serializibles.SerializibleObject{}, err
			}
			// YAML unmarshals into map[string]any but sometimes into map[any]any.
			// gopkg.in/yaml.v3 unmarshals into map[string]any.
			return wrapValue(&v), nil
		},

		SerializeToJson: func(data serializibles.SerializibleObject) (string, error) {
			raw := reconstruct(data)
			bytes, err := json.Marshal(raw)
			if err != nil {
				return "", err
			}
			return string(bytes), nil
		},
		SerializeToYaml: func(data serializibles.SerializibleObject) (string, error) {
			raw := reconstruct(data)
			bytes, err := yaml.Marshal(raw)
			if err != nil {
				return "", err
			}
			return string(bytes), nil
		},
	}
}

// wrapValue takes a pointer to an 'any' variable and returns a functional
// SerializibleObject contract around it.
func wrapValue(val *any) serializibles.SerializibleObject {
	return serializibles.SerializibleObject{
		IsInt: func() bool {
			if val == nil || *val == nil {
				return false
			}
			switch v := (*val).(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				return true
			case float64:
				return v == float64(int64(v))
			case float32:
				return v == float32(int32(v))
			default:
				return false
			}
		},
		IsString: func() bool {
			if val == nil || *val == nil {
				return false
			}
			_, ok := (*val).(string)
			return ok
		},
		IsFloat: func() bool {
			if val == nil || *val == nil {
				return false
			}
			switch (*val).(type) {
			case float32, float64:
				return true
			default:
				return false
			}
		},
		IsBool: func() bool {
			if val == nil || *val == nil {
				return false
			}
			_, ok := (*val).(bool)
			return ok
		},
		IsNull: func() bool {
			return val == nil || *val == nil
		},
		IsObject: func() bool {
			if val == nil || *val == nil {
				return false
			}
			_, ok := (*val).(map[string]any)
			return ok
		},
		IsArray: func() bool {
			if val == nil || *val == nil {
				return false
			}
			_, ok := (*val).([]any)
			return ok
		},

		GetInt: func() (int64, error) {
			if val == nil || *val == nil {
				return 0, fmt.Errorf("not an int")
			}
			switch v := (*val).(type) {
			case int:
				return int64(v), nil
			case int64:
				return v, nil
			case float64:
				return int64(v), nil
			case float32:
				return int64(v), nil
			default:
				return 0, fmt.Errorf("not an int")
			}
		},
		GetFloat: func() (float64, error) {
			if val == nil || *val == nil {
				return 0, fmt.Errorf("not a float")
			}
			switch v := (*val).(type) {
			case float32:
				return float64(v), nil
			case float64:
				return v, nil
			case int:
				return float64(v), nil
			case int64:
				return float64(v), nil
			default:
				return 0, fmt.Errorf("not a float")
			}
		},
		GetString: func() (string, error) {
			if val == nil || *val == nil {
				return "", fmt.Errorf("not a string")
			}
			if v, ok := (*val).(string); ok {
				return v, nil
			}
			return "", fmt.Errorf("not a string")
		},
		GetBool: func() (bool, error) {
			if val == nil || *val == nil {
				return false, fmt.Errorf("not a bool")
			}
			if v, ok := (*val).(bool); ok {
				return v, nil
			}
			return false, fmt.Errorf("not a bool")
		},

		GetObjectItem: func(key string) (serializibles.SerializibleObject, error) {
			if val == nil || *val == nil {
				return serializibles.SerializibleObject{}, fmt.Errorf("not an object")
			}
			m, ok := (*val).(map[string]any)
			if !ok {
				return serializibles.SerializibleObject{}, fmt.Errorf("not an object")
			}
			if v, exists := m[key]; exists {
				// We create a new pointer to the map's value. 
				// NOTE: modifying the returned child won't automatically update the map!
				// You must use ReplaceItemInObject if you want to mutate and persist the change.
				ptr := new(any)
				*ptr = v
				return wrapValue(ptr), nil
			}
			return serializibles.SerializibleObject{}, fmt.Errorf("key not found")
		},
		HasKey: func(key string) bool {
			if val == nil || *val == nil {
				return false
			}
			m, ok := (*val).(map[string]any)
			if !ok {
				return false
			}
			_, exists := m[key]
			return exists
		},
		GetKeys: func() ([]string, error) {
			if val == nil || *val == nil {
				return nil, fmt.Errorf("not an object")
			}
			m, ok := (*val).(map[string]any)
			if !ok {
				return nil, fmt.Errorf("not an object")
			}
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			return keys, nil
		},

		GetArrayItem: func(index int) (serializibles.SerializibleObject, error) {
			if val == nil || *val == nil {
				return serializibles.SerializibleObject{}, fmt.Errorf("not an array")
			}
			arr, ok := (*val).([]any)
			if !ok {
				return serializibles.SerializibleObject{}, fmt.Errorf("not an array")
			}
			if index < 0 || index >= len(arr) {
				return serializibles.SerializibleObject{}, fmt.Errorf("index out of bounds")
			}
			v := arr[index]
			ptr := new(any)
			*ptr = v
			return wrapValue(ptr), nil
		},
		GetArraySize: func() (int, error) {
			if val == nil || *val == nil {
				return 0, fmt.Errorf("not an array")
			}
			arr, ok := (*val).([]any)
			if !ok {
				return 0, fmt.Errorf("not an array")
			}
			return len(arr), nil
		},

		AddItemToObject: func(key string, item serializibles.SerializibleObject) error {
			if val == nil || *val == nil {
				return fmt.Errorf("not an object")
			}
			m, ok := (*val).(map[string]any)
			if !ok {
				return fmt.Errorf("not an object")
			}
			m[key] = reconstruct(item)
			return nil
		},
		ReplaceItemInObject: func(key string, item serializibles.SerializibleObject) error {
			if val == nil || *val == nil {
				return fmt.Errorf("not an object")
			}
			m, ok := (*val).(map[string]any)
			if !ok {
				return fmt.Errorf("not an object")
			}
			if _, exists := m[key]; !exists {
				return fmt.Errorf("key not found")
			}
			m[key] = reconstruct(item)
			return nil
		},
		DeleteItemFromObject: func(key string) error {
			if val == nil || *val == nil {
				return fmt.Errorf("not an object")
			}
			m, ok := (*val).(map[string]any)
			if !ok {
				return fmt.Errorf("not an object")
			}
			delete(m, key)
			return nil
		},

		AddItemToArray: func(item serializibles.SerializibleObject) error {
			if val == nil || *val == nil {
				return fmt.Errorf("not an array")
			}
			arr, ok := (*val).([]any)
			if !ok {
				return fmt.Errorf("not an array")
			}
			arr = append(arr, reconstruct(item))
			*val = arr
			return nil
		},
		DeleteItemFromArray: func(index int) error {
			if val == nil || *val == nil {
				return fmt.Errorf("not an array")
			}
			arr, ok := (*val).([]any)
			if !ok {
				return fmt.Errorf("not an array")
			}
			if index < 0 || index >= len(arr) {
				return fmt.Errorf("index out of bounds")
			}
			arr = append(arr[:index], arr[index+1:]...)
			*val = arr
			return nil
		},
	}
}

// reconstruct recursively turns a SerializibleObject back into a Go interface{} (any).
func reconstruct(item serializibles.SerializibleObject) any {
	if item.IsNull() {
		return nil
	}
	if item.IsInt() {
		v, _ := item.GetInt()
		return v
	}
	if item.IsFloat() {
		v, _ := item.GetFloat()
		return v
	}
	if item.IsString() {
		v, _ := item.GetString()
		return v
	}
	if item.IsBool() {
		v, _ := item.GetBool()
		return v
	}
	if item.IsArray() {
		size, _ := item.GetArraySize()
		arr := make([]any, size)
		for i := 0; i < size; i++ {
			child, _ := item.GetArrayItem(i)
			arr[i] = reconstruct(child)
		}
		return arr
	}
	if item.IsObject() {
		keys, _ := item.GetKeys()
		m := make(map[string]any)
		for _, k := range keys {
			child, _ := item.GetObjectItem(k)
			m[k] = reconstruct(child)
		}
		return m
	}
	return nil
}

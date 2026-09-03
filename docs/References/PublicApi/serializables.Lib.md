# `serializables.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	CreateString func(value string) *SerializibleObject
	CreateInt    func(value int64) *SerializibleObject
	CreateFloat  func(value float64) *SerializibleObject
	CreateBool   func(value bool) *SerializibleObject
	CreateNull   func() *SerializibleObject
	CreateObject func() *SerializibleObject
	CreateArray  func() *SerializibleObject

	ParseJson func(data string) (*SerializibleObject, error)
	ParseYaml func(data string) (*SerializibleObject, error)

	SerializeToJson func(data *SerializibleObject) string
	SerializeToYaml func(data *SerializibleObject) string
}

type SerializibleObject struct {
	IsInt, IsString, IsFloat, IsBool, IsNull, IsObject, IsArray func() bool

	GetInt    func() (int64, error)
	GetFloat  func() (float64, error)
	GetString func() (string, error)
	GetBool   func() (bool, error)

	GetObjectItem func(key string) (*SerializibleObject, error)
	HasKey        func(key string) bool
	GetKeys       func() ([]string, error)

	GetArrayItem func(index int) *SerializibleObject
	GetArraySize func() (int, error)

	AddItemToObject      func(key string, item any) error
	ReplaceItemInObject  func(key string, item any) error
	DeleteItemFromObject func(key string) error

	AddItemToArray      func(item any) error
	DeleteItemFromArray func(index int) error
}
```

## Description

A generic JSON/YAML value library injected whole as `Deps.Serializables`, so the sandbox imports no serialization module. A `SerializibleObject` wraps one dynamically typed value — scalar, object or array — with predicates, typed getters, and mutators; `Lib` creates them, parses text into them, and renders them back. Every parsable config under `sandbox/internal/parsables/` reads and writes through it (`ParseYaml` in `New`, `SerializeToYaml` in `Render`), which is why re-rendered YAML comes out with keys in the serializer's alphabetical order — see [HandleParsables.md](/docs/Tutorials/HandleParsables.md). The standard adapter backs it with `gopkg.in/yaml.v3` and `encoding/json`. Installed by the `serializables` dep.

## Fields

| Field | Description |
| :--- | :--- |
| `Create*` | Wrap a scalar, an empty object, or an empty array. |
| `ParseJson`, `ParseYaml` | Parse text; error on malformed input. |
| `SerializeToJson`, `SerializeToYaml` | Render a value tree to text. |
| `SerializibleObject.Is*` | Kind predicates. |
| `SerializibleObject.Get{Int,Float,String,Bool}` | Typed scalar access; error on kind mismatch. |
| `SerializibleObject.GetObjectItem`, `HasKey`, `GetKeys` | Object access. |
| `SerializibleObject.GetArrayItem`, `GetArraySize` | Array access. |
| `SerializibleObject.{Add,Replace,Delete}ItemFromObject`, `{Add,Delete}ItemFromArray` | Mutation; `item` is any Go value or another `*SerializibleObject`. |

## Examples

```go
root, err := deps.Serializables.ParseYaml(content)
if err != nil {
	return nil, err
}
if root.HasKey("version") {
	item, _ := root.GetObjectItem("version")
	version, _ := item.GetString()
	conf.Version = version
}

out := deps.Serializables.CreateObject()
out.AddItemToObject("name", conf.Name)
out.AddItemToObject("version", conf.Version)
text := deps.Serializables.SerializeToYaml(out)
```

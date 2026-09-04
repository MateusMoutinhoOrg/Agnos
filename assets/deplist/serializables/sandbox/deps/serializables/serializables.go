package serializibles

// SerializibleObject is one node of a parsed document — a scalar, an object or
// an array — and the whole tree is navigated and edited through its function
// fields. The same struct is what the Create* constructors of Lib return, so a
// document can be built in memory and serialized without ever being parsed.
type SerializibleObject struct {
	IsInt    func() bool
	IsString func() bool
	IsFloat  func() bool
	IsBool   func() bool
	IsNull   func() bool
	IsObject func() bool
	IsArray  func() bool

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

// Lib is the JSON/YAML codec injected whole as the Deps.Serializables field:
// constructors for every node kind, the two parsers and the two serializers.
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

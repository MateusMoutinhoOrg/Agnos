package serializabledeps

// SerializableObject is one node of a parsed document — a scalar, an object or
// an array — and the whole tree is navigated and edited through its function
// fields. The same struct is what the Create* constructors of Contract return, so a
// document can be built in memory and serialized without ever being parsed.
type SerializableObject struct {
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

	GetObjectItem func(key string) (*SerializableObject, error)
	HasKey        func(key string) bool
	GetKeys       func() ([]string, error)

	GetArrayItem func(index int) *SerializableObject
	GetArraySize func() (int, error)

	AddItemToObject      func(key string, item any) error
	ReplaceItemInObject  func(key string, item any) error
	DeleteItemFromObject func(key string) error

	AddItemToArray      func(item any) error
	DeleteItemFromArray func(index int) error
}

// Contract is the JSON/YAML codec injected whole as the Deps.SerializableDeps field:
// constructors for every node kind, the two parsers and the two serializers.
type Contract struct {
	CreateString func(value string) *SerializableObject
	CreateInt    func(value int64) *SerializableObject
	CreateFloat  func(value float64) *SerializableObject
	CreateBool   func(value bool) *SerializableObject
	CreateNull   func() *SerializableObject
	CreateObject func() *SerializableObject
	CreateArray  func() *SerializableObject

	ParseJson func(data string) (*SerializableObject, error)
	ParseYaml func(data string) (*SerializableObject, error)

	SerializeToJson func(data *SerializableObject) string
	SerializeToYaml func(data *SerializableObject) string
}

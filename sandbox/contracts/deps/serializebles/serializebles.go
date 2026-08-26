package serializibles

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

	GetObjectItem func(key string) (SerializibleObject, error)
	HasKey        func(key string) bool
	GetKeys       func() ([]string, error)

	GetArrayItem func(index int) (SerializibleObject, error)
	GetArraySize func() (int, error)

	AddItemToObject      func(key string, item SerializibleObject) error
	ReplaceItemInObject  func(key string, item SerializibleObject) error
	DeleteItemFromObject func(key string) error

	AddItemToArray      func(item SerializibleObject) error
	DeleteItemFromArray func(index int) error
}

type Lib struct {
	CreateString func(value string) SerializibleObject
	CreateInt    func(value int64) SerializibleObject
	CreateFloat  func(value float64) SerializibleObject
	CreateBool   func(value bool) SerializibleObject
	CreateNull   func() SerializibleObject
	CreateObject func() SerializibleObject
	CreateArray  func() SerializibleObject

	ParseJson func(data string) (SerializibleObject, error)
	ParseYaml func(data string) (SerializibleObject, error)

	SerializeToJson func(data SerializibleObject) (string, error)
	SerializeToYaml func(data SerializibleObject) (string, error)
}

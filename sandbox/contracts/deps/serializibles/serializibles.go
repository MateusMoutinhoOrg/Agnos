package serializibles

type SerializibleObject struct {
	IsInt    func() bool
	IsString func() bool
	IsFloat  func() bool
	IsBool   func() bool
	IsObject func() bool
	IsArray  func() bool
}

type Lib struct {
	ParseYaml func(data string) (SerializibleObject, error)
	ParseJson func(data string) (SerializibleObject, error)

	SerializeToYaml func(data SerializibleObject) (string, error)
	SerializeToJson func(data SerializibleObject) (string, error)
}

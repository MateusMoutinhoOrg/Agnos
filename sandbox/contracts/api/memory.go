package api

type Memory struct {
	GetProjectName func() string
}

type MemoryApi struct {
	NewMemory func(path string) Memory
}

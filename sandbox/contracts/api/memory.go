package api

type Memory struct {
}

type MemoryApi struct {
	NewMemory func(path string) Memory
}

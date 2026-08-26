package api

type CoreApi struct {
	Start func(path string) error
	Build func(path string) error
}

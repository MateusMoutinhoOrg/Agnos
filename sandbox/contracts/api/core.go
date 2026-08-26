package api

type CoreApi struct {
	Start func(path string) error
}

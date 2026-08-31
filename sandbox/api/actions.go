package api

type Actions struct {
	Build func(path string) error
}

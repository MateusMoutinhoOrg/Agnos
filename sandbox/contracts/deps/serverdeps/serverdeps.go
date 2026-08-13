package serverdeps

type Response struct {
	GetStatusCode func() int
	GetHeader     func(key string) string
	ReadBody      func(size int) ([]byte, error) // size == -1 lê tudo
	Close         func() error
}

type Request struct {
	AddHeader func(key string, value string)
	SetMethod func(method string)
	SetBody   func(body []byte)
	Fetch     func() (Response, error)
}

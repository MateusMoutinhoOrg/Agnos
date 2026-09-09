package api

// Server is the http surface of the sandbox. Serve is the generated
// route-and-dispatch entry point (see sandbox/internal/server/servermain.go),
// wired in by sandbox/binds/server.go. It blocks until the server stops.
type Server struct {
	Serve func(props ServeProps) error
}

// ServeProps describes one run of the http server: the address to listen on
// and the two timeouts, in milliseconds, a request and a response are held to.
type ServeProps struct {
	Addr           string
	ReadTimeoutMs  int
	WriteTimeoutMs int
}

const (
	// StatusOk reports that the route did what it was asked to do.
	StatusOk = 200
	// StatusCreated reports that the route created what it was asked for.
	StatusCreated = 201
	// StatusNoContent reports success with nothing to send back.
	StatusNoContent = 204
	// StatusBadRequest reports a request the dispatch could not bind: a
	// missing required field, an unparsable value, one out of range, or a
	// body the declared schema rejects.
	StatusBadRequest = 400
	// StatusNotFound reports that no declared route matches the path.
	StatusNotFound = 404
	// StatusMethodNotAllowed reports a path a route matches under another
	// method.
	StatusMethodNotAllowed = 405
	// StatusConflict reports a well-formed request the current state refuses.
	StatusConflict = 409
	// StatusPayloadTooLarge reports a body longer than the route's
	// max-bytes.
	StatusPayloadTooLarge = 413
	// StatusUnsupportedMedia reports a content type the route does not
	// declare.
	StatusUnsupportedMedia = 415
	// StatusFailure reports a well-formed request the route could not carry
	// out.
	StatusFailure = 500
)

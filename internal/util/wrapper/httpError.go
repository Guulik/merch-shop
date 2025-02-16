package wrapper

type HTTPError struct {
	InternalErr error
	Status      int
	Msg         string
}

func (e *HTTPError) Error() string {
	return e.Msg
}

func WrapHTTPError(err error, status int, msg string) *HTTPError {
	return &HTTPError{InternalErr: err, Status: status, Msg: msg}
}

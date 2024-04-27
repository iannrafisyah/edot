package utilities

type Error struct {
	StatusCode int
	Data       interface{}
	Message    error
}

// Error :
func (r *Error) Error() string {
	return r.Message.Error()
}

// ErrorRequest :
func ErrorRequest(err error, httpCode int, data ...interface{}) error {
	var errData interface{}

	if len(data) > 0 {
		errData = data[0]
	}

	return &Error{
		Message:    err,
		Data:       errData,
		StatusCode: httpCode,
	}
}

// ParseError :
func ParseError(r error) *Error {
	errInfo, ok := r.(*Error)
	if ok {
		return errInfo
	}
	return nil
}

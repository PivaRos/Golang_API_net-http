package err

import "net/http"

type Error struct {
	code int
	Err  error
}

func CreateInternalError(code int, Err error) Error {
	return Error{
		code: code,
		Err:  Err,
	}
}

func (e *Error) HandleResponse(w http.ResponseWriter) {
	http.Error(w, e.Err.Error(), e.code)
}

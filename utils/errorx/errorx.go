package errorx

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorWrapper struct {
	Code   int                `json:"-"` // HTTP Status code. We use `-` to skip json marshaling.
	Errors *map[string]string `json:"-"` // The original error. Same reason as above.
}

func NewError(statusCode int, err *map[string]string) error {
	return ErrorWrapper{
		Code:   statusCode,
		Errors: err,
	}
}

func NewErrorForStatusCode(statusCode int) error {
	return ErrorWrapper{
		Code:   statusCode,
		Errors: &map[string]string{},
	}
}

func NewErrorForNonFieldError(statusCode int, err error) error {
	return ErrorWrapper{
		Code:   statusCode,
		Errors: &map[string]string{"non_field_error": err.Error()},
	}
}

func NewErrorForKV(statusCode int, key string, value string) error {
	return ErrorWrapper{
		Code:   statusCode,
		Errors: &map[string]string{key: value},
	}
}

// Error function used to implement the `error` interface for returning errors.
func (err ErrorWrapper) Error() string {
	b, e := json.Marshal(err.Errors)
	if e != nil { // Defensive code
		return e.Error()
	}
	return string(b)
}

func NewBadRequestError(err *map[string]string) error {
	return ErrorWrapper{
		Code:   http.StatusBadRequest,
		Errors: err,
	}
}

func NewBadRequestErrorForNonField(err error) error {
	return ErrorWrapper{
		Code:   http.StatusBadRequest,
		Errors: &map[string]string{"non_field_error": err.Error()},
	}
}

func NewBadRequestErrorForKV(key string, value string) error {
	return ErrorWrapper{
		Code:   http.StatusBadRequest,
		Errors: &map[string]string{key: value},
	}
}

// ResponseError function returns the HTTP error response based on the httpcode used.
func ResponseError(rw http.ResponseWriter, err error) {
	// Copied from:
	// https://dev.to/tigorlazuardi/go-creating-custom-error-wrapper-and-do-proper-error-equality-check-11k7

	rw.Header().Set("Content-Type", "Application/json")

	//
	// CASE 1 OF 2: Handle API Errors.
	//

	var ew ErrorWrapper
	if errors.As(err, &ew) {
		rw.WriteHeader(ew.Code)
		_ = json.NewEncoder(rw).Encode(ew.Errors)
		return
	}

	//
	// CASE 2 OF 2: Handle non ErrorWrapper types.
	//

	rw.WriteHeader(http.StatusInternalServerError)

	_ = json.NewEncoder(rw).Encode(err.Error())
}

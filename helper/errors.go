package helper

import "errors"

var ErrNotFound = errors.New("not found der")
var ErrUnauthorize = errors.New("unauthorize")

type AppError struct {
	// Code is the error code.
	Code int `json:"code"`

	// Message is the error message.
	Message string `json:"message"`

	// Description is the error description.
	Description string `json:"description,omitempty"`

	// Detail is the error detail of validation
	Detail *ErrDetail `json:"detail,omitempty"`
}

// ErrDetail is used to represent an error detail.
type ErrDetail struct {
	// Kind is the error kind.
	Kind string `json:"kind,omitempty"`

	// Field is the error field.
	Field string `json:"field,omitempty"`

	// Description is the error description.
	Description string `json:"description,omitempty"`
}

func (e AppError) Error() string {
	return ""
}

func GetHttpStatus(e error) int {
	v, ok := e.(AppError)
	if ok {
		return v.Code
	} else {
		switch {
		case e == nil:
			return 200
		case errors.Is(e, ErrNotFound):
			return 401
		case errors.Is(e, ErrUnauthorize):
			return 403
		default:
			return 500
		}
	}
}

func GetErrMessage(e error) error {
	v, ok := e.(AppError)
	if ok {
		if v.Message == "" {
			v.Message = GetMessage(v.Code)
			return v
		}
		return v
	} else {
		return AppError{
			Code:    GetHttpStatus(e),
			Message: e.Error(),
		}
	}
}

func GetMessage(n int) string {
	switch n {
	case 401:
		return "not found"
	case 403:
		return "unauthorize"
	case 400:
		return "bad request"
	default:
		return "internal server error"
	}
}

package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

var (
	errFailedToEncodeResponse = errors.New("failed to encode response")
)

// NewResponse creates a new response
func NewResponse[T any](data T) *Response[T] {
	switch t := any(data).(type) {
	case tools.ErrorsMap:
		return &Response[T]{
			Success: false,
			Data:    data,
		}
	case error:
		return &Response[T]{
			Success: false,
			Message: t.Error(),
		}
	default:
		return &Response[T]{
			Success: true,
			Data:    data,
		}
	}
}

// Response is a generic response structure
type (
	Response[T any] struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
		Data    T      `json:"data,omitempty"`
	}

	isError interface {
		IsError() error
	}
)

// Write writes the response to the writer
func (r *Response[T]) Write(code int, wr http.ResponseWriter) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
	if err != nil {
		if r.Message != errFailedToEncodeResponse.Error() {
			NewResponse(errFailedToEncodeResponse).
				Write(http.StatusInternalServerError, wr)
		}

		return
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(code)
	buf.WriteTo(wr)
}

func (r *Response[T]) isError() error {
	isErrorData, ok := any(r.Data).(isError)
	if ok {
		return isErrorData.IsError()
	}

	if len(r.Message) > 0 {
		return errors.New("response errors: " + r.Message)
	}

	return errors.New("unknown response error")
}

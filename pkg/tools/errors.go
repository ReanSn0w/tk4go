package tools

import (
	"bytes"
	"fmt"
)

// ErrorsMap is a map of errors
func NewErrorsMap() ErrorsMap {
	return make(ErrorsMap)
}

type ErrorsMap map[string]error

// IsError returns nil if there are no errors
func (e ErrorsMap) IsError() error {
	if len(e) == 0 {
		return nil
	}

	return e
}

// Error returns a string representation of the errors map
func (e ErrorsMap) Error() string {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf("errors map has %v errors\n", len(e)))

	for k, v := range e {
		buf.WriteString(k)
		buf.WriteString(": ")
		buf.WriteString(v.Error())
		buf.WriteString("\n")
	}

	return buf.String()
}

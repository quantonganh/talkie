package talkie

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	ErrInvalid      = "invalid"
	ErrUnauthorized = "unauthorized"
	ErrNotFound     = "not_found"
	ErrConflict     = "conflict"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrCommentNotFound = errors.New("comment not found")
)

// Error represents an app error
// This can be used to distinguish between client and server error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Err     error  `json:"err,omitempty"`
}

func (e *Error) Error() string {
	var buf bytes.Buffer

	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			fmt.Fprintf(&buf, "<%s>", e.Code)
		}
		buf.WriteString(e.Message)
	}
	return buf.String()
}

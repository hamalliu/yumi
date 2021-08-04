package status

import (
	"fmt"

	"yumi/pkg/codes"
)

// Status ...
type Status struct {
	code    int32
	message I18nMessageID
	err     error
}

// New returns a Status representing c and msg.
func new(c codes.Code) *Status {
	return &Status{code: int32(c)}
}

// Code ...
func (s *Status) Code() codes.Code {
	if s == nil {
		return codes.OK
	}
	return codes.Code(s.code)
}

// Message ...
func (s *Status) Message(language string) string {
	if s == nil {
		return ""
	}

	return s.message.T(language)
}

// Err ...
func (s *Status) Err() error {
	return s.err
}

// Error ...
func (s *Status) Error() string {
	if s.err == nil {
		return ""
	}

	return s.err.Error()
}

// WithError ...
func (s *Status) WithError(err error) *Status {
	s.err = fmt.Errorf("status: code = %d; message = %s; err = %w", s.code, s.message, err)

	return s
}

// WithMessage ...
func (s *Status) WithMessage(msg I18nMessageID) *Status {
	s.message = msg
	return s
}

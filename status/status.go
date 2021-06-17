package status

import (
	"fmt"

	"yumi/pkg/codes"
)

// Status ...
type Status struct {
	code    int32
	message string
	details []string
}

// New returns a Status representing c and msg.
func new(c codes.Code, msg string) *Status {
	return &Status{code: int32(c), message: msg}
}

// Code ...
func (s *Status) Code() codes.Code {
	if s == nil {
		return codes.OK
	}
	return codes.Code(s.code)
}

// Message ...
func (s *Status) Message() string {
	if s == nil {
		return ""
	}

	return s.message
}

// Details ...
func (s *Status) Details() []string {
	if s == nil {
		return []string{}
	}
	return s.details
}

// Err ...
func (s *Status) Err() error {
	if s.Code() == codes.OK {
		return nil
	}

	return s
}

// Error ...
func (s *Status) Error() string {
	if s.Code() == codes.OK {
		return ""
	}

	return fmt.Sprintf("status: code =%d desc = %s", s.code, s.message)
}

// WithDetails ...
func (s *Status) WithDetails(details ...error) *Status {
	for _, detail := range details {
		s.details = append(s.details, fmt.Sprintf("%+v", detail))
	}

	return s
}

// WithMessage ...
func (s *Status) WithMessage(msg messageEntry) *Status {
	// TODO 实现i18n
	s.message = ""
	return s
}

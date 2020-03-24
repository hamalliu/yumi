package internal_error

import (
	goerrors "errors"
	"fmt"

	"github.com/pkg/errors"

	"yumi/utils/log"
)

//内部错误，日志级别：error
type Error struct {
	s string
}

func (intn *Error) Error() string {
	return intn.s
}

func New(message string) error {
	err := fmt.Errorf(message)
	err = errors.WithStack(err)
	log.Error2(fmt.Sprintf("%s\n", err.Error()))
	return err
}

func Newf(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	err = errors.WithStack(err)
	log.Error2(fmt.Sprintf("%s\n", err.Error()))
	return err
}

func With(err error) error {
	var targetErr *Error
	if !goerrors.As(err, &targetErr) {
		err = &Error{s: fmt.Sprintf("%+v", errors.WithStack(err))}
	}
	log.Error2(fmt.Sprintf("%s\n", err.Error()))
	return err
}

func Critical(err error) error {
	var targetErr *Error
	if !goerrors.As(err, &targetErr) {
		err = &Error{s: fmt.Sprintf("%+v", errors.WithStack(err))}
	}
	log.Critical2(fmt.Sprintf("%s\n", err.Error()))
	return err
}

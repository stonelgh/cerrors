package cerrors

import "fmt"

// Structs and functions defined here are mainly inspired by
// https://github.com/pkg/errors/blob/master/errors.go
// Some of them are just a simple copy-paste.

type errWrapper struct {
	cause error
	msg   string
}

func (e *errWrapper) Error() string {
	return e.msg + ": " + e.cause.Error()
}

func (e *errWrapper) Cause() error {
	return e.cause
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}
	for {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

// Wrap returns an error annotating err with the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &errWrapper{err, msg}
}

// Wrapf returns an error annotating err with the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	return Wrap(err, fmt.Sprintf(format, args...))
}

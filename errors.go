package cerrors

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// Error represents an error condition.
// It offers more facilities to accurately identify and debug the condition.
type Error interface {
	error

	// Msg simply describes the error in a string.
	Msg() string

	// Code returns the identity of the error.
	Code() interface{}
	// CodeIs returns if the code of the error is same as the one specified.
	CodeIs(code interface{}) bool
	// CodeIn returns if the code of the error is one of the specified codes.
	CodeIn(codes ...interface{}) bool

	// StackV returns the stack of the error in a string vector, with each
	// frame in one string.
	StackV() []string
	// Stack returns the stack of the error in a string, with frames separated
	// by '\n'.
	Stack() string
}

// cerror implements the Error interface.
type cerror struct {
	code  interface{}
	msg   string
	stack []uintptr
}

// New returns an Error with the supplied code and message,
// and the stack trace at the point it was called.
func New(code interface{}, msg string) Error {
	stack := [10]uintptr{}
	n := runtime.Callers(2, stack[:])
	return &cerror{
		code:  code,
		msg:   msg,
		stack: stack[:n],
	}
}

// Newf formats according to a format specifier and calls New() to
// return an Error.
func Newf(code interface{}, format string, args ...interface{}) Error {
	return New(code, fmt.Sprintf(format, args...))
}

func (e *cerror) Error() string {
	if e.code == nil {
		return e.msg
	}
	if len(e.msg) == 0 {
		return fmt.Sprintf("Error %v", e.code)
	}
	return fmt.Sprintf("Error %v: %v", e.code, e.msg)
}

func (e *cerror) Code() interface{} {
	return e.code
}

func (e *cerror) Msg() string {
	return e.msg
}

func (e *cerror) CodeIs(code interface{}) bool {
	return e.code == code
}

func (e *cerror) CodeIn(codes ...interface{}) bool {
	for _, code := range codes {
		if e.code == code {
			return true
		}
	}
	return false
}

func (e *cerror) StackV() []string {
	if len(e.stack) == 0 {
		return []string{}
	}

	a := make([]string, len(e.stack))
	frames := runtime.CallersFrames(e.stack)
	for i := 0; i < len(a); i++ {
		frame, more := frames.Next()
		name := frame.Function
		if n := strings.LastIndex(name, "/"); n != -1 {
			name = name[n+1:]
		}
		if n := strings.Index(name, "."); n != -1 {
			name = name[n+1:]
		}
		a[i] = frame.File + ":" + strconv.Itoa(frame.Line) + "\t" + name
		if !more {
			break
		}
	}
	return a
}

func (e *cerror) Stack() string {
	return strings.Join(e.StackV(), "\n")
}

// Code first calls Cause() to find the cause of the error.
// Then it returns the code of the cause.
// If the cause implements the following interface:
//
//     type coder interface {
//            Code() interface{}
//     }
//
// its code is the return value of Code().
// Otherwise, nil will be returned.
func Code(err error) interface{} {
	err = Cause(err)

	type coder interface {
		Code() interface{}
	}
	if err, ok := err.(coder); ok {
		return err.Code()
	}
	return nil
}

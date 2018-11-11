package cerrors

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func fooErr(code interface{}, msg string) Error {
	return New(code, msg)
}

type foo int

func (f foo) err(code interface{}, msg string) Error {
	return New(code, msg)
}

func TestNew(t *testing.T) {
	tests := []struct {
		maker      func(interface{}, string) Error
		code       interface{}
		msg        string
		want       error
		wantFrame  string
		unwantCode interface{}
	}{
		{
			fooErr,
			nil,
			"nil code",
			errors.New("nil code"),
			"/errors_test.go:11\tfooErr",
			"non-nil",
		}, {
			fooErr,
			123,
			"int code",
			errors.New("Error 123: int code"),
			"/errors_test.go:11\tfooErr",
			321,
		}, {
			fooErr,
			"NotFound",
			"",
			errors.New("Error NotFound"),
			"/errors_test.go:11\tfooErr",
			321,
		}, {
			foo(0).err,
			"BAD-INPUT",
			"string code",
			errors.New("Error BAD-INPUT: string code"),
			"/errors_test.go:17\tfoo.err", "GOOD-INPUT",
		},
	}

	for _, tt := range tests {
		got := tt.maker(tt.code, tt.msg)
		if got.Code() != tt.code || got.Msg() != tt.msg ||
			got.Error() != tt.want.Error() {
			t.Errorf("New(): got: %v, want %v", got, tt.want)
		}
		if got.Code() != tt.code || got.Code() == tt.unwantCode ||
			!got.CodeIs(tt.code) || got.CodeIs(tt.unwantCode) ||
			!got.CodeIn(tt.code) || got.CodeIn(tt.unwantCode) ||
			!got.CodeIn(tt.unwantCode, tt.code) || got.CodeIn(tt.unwantCode, "NEVER-BE") {
			t.Errorf("New.Code*() check failed: got: %v, want %v, unwant %v", got.Code(), tt.code, tt.unwantCode)
		}
		if !strings.HasSuffix(got.StackV()[0], tt.wantFrame) ||
			!strings.Contains(got.Stack(), tt.wantFrame+"\n") {
			t.Errorf("New.Stack*(): got stack:\n%v\nlack of frame: ...%v", got.Stack(), tt.wantFrame)
		}
	}
}

func TestNewf(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{
			Newf(123, "read error without format specifiers"),
			"Error 123: read error without format specifiers",
		}, {
			Newf("EBADF", "read error with %d format specifier", 1),
			"Error EBADF: read error with 1 format specifier",
		},
	}

	for _, tt := range tests {
		got := tt.err.Error()
		if got != tt.want {
			t.Errorf("Newf(%v): got: %q, want %q", tt.err, got, tt.want)
		}
	}
}

func TestCode(t *testing.T) {
	tests := []struct {
		err  error
		want interface{}
	}{
		{
			nil,
			nil,
		}, {
			errors.New("simple error"),
			nil,
		}, {
			Wrap(Wrap(io.EOF, "read error"), "client error"),
			nil,
		}, {
			New(123, "error with an int code"),
			123,
		}, {
			Wrap(New("EBADF", "error with a string code"), "read error"),
			"EBADF",
		},
	}

	for _, tt := range tests {
		got := Code(tt.err)
		if got != tt.want {
			t.Errorf("Code(%v): got: %q, want %q", tt.err, got, tt.want)
		}
	}
}

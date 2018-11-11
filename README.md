# cerrors
[![Build Status](https://travis-ci.org/stonelgh/cerrors.svg?branch=master)](https://travis-ci.org/stonelgh/cerrors)
[![GoDoc](https://godoc.org/github.com/stonelgh/cerrors?status.svg)](http://godoc.org/github.com/stonelgh/cerrors)
[![Report card](https://goreportcard.com/badge/github.com/stonelgh/cerrors)](https://goreportcard.com/report/github.com/stonelgh/cerrors)

Package cerrors provides simple yet flexible error handling primitives.

It allows you to categorize errors using error codes. It helps you to debug
by recording the call stack when the error occurs and allowing you to add
context to the failure path.

To create and test an error with an error code:
```go
const ETimedout = 110
err1 := cerrors.New(ETimedout, "opertion timed out")
if err1.CodeIs(ETimedout) {
    // handle the case operation timed out
}

const EBusy = "BUSY"
err2 := cerrors.New(EBusy, "device is busy")
if cerrors.Code(err2) == EBusy {
    // retry after a while
}
```

Print the call stack when the error happens:
```go
fmt.Println(err1.Stack())
```

Add context to an error:
```go
err3 := cerrors.Wrap(err2, "read from disk a")
fmt.Println(err3)
```

## Note on error codes
Returning errors with a code can be handy sometimes. However, don't abuse them.
They'll create dependencies between packages and probably cause more troubles
than benefits in large projects.

## Credits
Some parts are inspired greatly by [pkg/errors](https://github.com/pkg/errors).

## License
[MIT License](LICENSE)

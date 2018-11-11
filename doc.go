/*
Package cerrors provides simple yet flexible error handling primitives.

It allows you to categorize errors using error codes. It helps you to debug
by recording the call stack when the error occurs and allowing you to add
context to the failure path.

To create and test an error with an error code:

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

Print the call stack when the error happens:

    fmt.Println(err1.Stack())

Add context to an error:

    err3 := cerrors.Wrap(err2, "read from disk a")
    fmt.Println(err3)

*/
package cerrors

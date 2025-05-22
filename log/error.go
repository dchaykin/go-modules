package log

import (
	"fmt"
	"runtime"
)

type errorWithPos struct {
	File string
	Line int
	Err  error
}

func (e *errorWithPos) Error() string {
	return fmt.Sprintf("%s:%d: %v", e.File, e.Line, e.Err)
}

func WrapError(err error) error {
	if err == nil {
		return nil
	}
	_, file, line, _ := runtime.Caller(1)
	return &errorWithPos{File: file, Line: line, Err: err}
}

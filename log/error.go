package log

import (
	"fmt"
	"runtime"
)

type stackTrace struct {
	file string
	line int
}

type errorWithPos struct {
	stackTrace []stackTrace
	err        error
}

func (e *errorWithPos) Error() string {
	result := fmt.Sprintf("%v", e.err)
	for _, s := range e.stackTrace {
		result += fmt.Sprintf("\n%s:%d", s.file, s.line)
	}
	return result
}

func (e *errorWithPos) appendTrace(file string, line int) {
	e.stackTrace = append(e.stackTrace, stackTrace{file: file, line: line})
}

func WrapError(err error) error {
	_, file, line, _ := runtime.Caller(1)
	switch e := err.(type) {
	case *errorWithPos:
		e.appendTrace(file, line)
		return e
	}
	result := &errorWithPos{err: err}
	result.appendTrace(file, line)
	return result
}

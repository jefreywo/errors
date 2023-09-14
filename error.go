package util

import (
	"fmt"
	"io"
)

var depthNum int = 3
var errorCode int = 14
var warningCode int = 333

func InitErrors(depth, errCode, warnCode int) {
	depthNum = depth
	errorCode = errCode
	warningCode = warnCode
}

func ErrorNew(format string, args ...interface{}) error {
	return &fundamentalError{
		code:  errorCode,
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

type fundamentalError struct {
	code int
	msg  string
	*stack
}

func (f *fundamentalError) Error() string {
	return fmt.Sprintf("msg:%v", f.msg)
}

func (fund *fundamentalError) Format(fs fmt.State, verb rune) {
	switch verb {
	case 'v':
		if fs.Flag('+') {
			io.WriteString(fs, fund.msg)
			fund.stack.Format(fs, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(fs, fund.msg)
	case 'q':
		fmt.Fprintf(fs, "%q", fund.msg)
	}
}

// WarnNew()返回一个warn级别的错误
func WarnNew(format string, args ...interface{}) error {
	return &warnError{
		code:  warningCode,
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// warn级别的错误
type warnError struct {
	code int
	msg  string
	*stack
}

func (w *warnError) Error() string {
	return fmt.Sprintf("msg:%v", w.msg)
}

func (fund *warnError) Format(fs fmt.State, verb rune) {
	switch verb {
	case 'v':
		if fs.Flag('+') {
			io.WriteString(fs, fund.msg)
			fund.stack.Format(fs, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(fs, fund.msg)
	case 'q':
		fmt.Fprintf(fs, "%q", fund.msg)
	}
}

// 判断是否为warn级别的错误，true-是，false-不是
func IsWarningError(err error) bool {
	if err == nil {
		return false
	}
	er, ok := err.(*warnError)
	if !ok {
		return false
	}

	if er.code == warningCode {
		return true
	}

	return false
}

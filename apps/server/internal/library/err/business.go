// Package err 业务异常：携带错误码的领域错误，由统一响应中间件映射为三段式契约。
package err

import "fmt"

// Business 业务异常（实现 error；Code() 供中间件提取错误码）。
type Business struct {
	code int
	msg  string
	err  error
}

func New(code int, msg string) *Business {
	return &Business{code: code, msg: msg}
}

func Wrap(code int, msg string, err error) *Business {
	return &Business{code: code, msg: msg, err: err}
}

func (e *Business) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

func (e *Business) Code() int      { return e.code }
func (e *Business) Message() string { return e.msg }
func (e *Business) Unwrap() error  { return e.err }

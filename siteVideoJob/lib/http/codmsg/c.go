package codmsg

import (
	"fmt"
	"io"
)

var (
	DefaultReturn = &CodMsg{Code: 6020, Msg: "请求失败，请稍后重试!"}
	InvalidParam  = &CodMsg{Code: 6010, Msg: "参数错误"}
	RequestFail   = &CodMsg{Code: 6030, Msg: "请求失败"}
	SendFail      = &CodMsg{Code: 6040, Msg: "发送失败"}
	VerifyFail    = &CodMsg{Code: 6050, Msg: "校验失败"}
	SendFailError = Wrap(SendFail, nil)
)

type CodMsg struct {
	Code int
	Msg  string
}

type cmError struct {
	*CodMsg
	cause error
}

func (e cmError) Error() string {
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Msg)
}

func WithMsg(msg string) error {
	return &cmError{
		CodMsg: &CodMsg{
			Code: DefaultReturn.Code,
			Msg:  msg,
		},
		cause: nil,
	}
}

func WithMsgErr(msg string, err error) error {
	return &cmError{
		CodMsg: &CodMsg{
			Code: DefaultReturn.Code,
			Msg:  msg,
		},
		cause: err,
	}
}

func WithMsgCode(msg string, code int) error {
	return &cmError{
		CodMsg: &CodMsg{
			Code: code,
			Msg:  msg,
		},
		cause: nil,
	}
}

func Wrap(cm *CodMsg, err error) error {
	return &cmError{
		CodMsg: cm,
		cause:  err,
	}
}

func IsCmError(err error) (*cmError, bool) {
	if b, ok := err.(*cmError); ok {
		return b, true
	}
	return nil, false
}

func (ae *cmError) Cause() error { return ae.cause }

func (ae *cmError) Unwrap() error { return ae.cause }

func (ae *cmError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", ae.Cause())
			//ae.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, ae.Error())
	case 'q':
		fmt.Fprintf(s, "%q", ae.Error())
	}
}

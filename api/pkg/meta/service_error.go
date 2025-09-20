package meta

import (
	"errors"
	"microservice/pkg/meta/status"
)

type (
	IServiceErr interface {
		Data(items map[string]any) *Error
		SetErr(err string) *Error
		Error() string
		Is(target error) bool
	}

	Error struct {
		Msg    status.HttpMappedStatus
		Err    error
		Detail map[string]any
	}
)

func ServiceErr(msg status.HttpMappedStatus, err ...error) *Error {
	se := &Error{Msg: msg, Err: errors.New("")}
	if len(err) > 0 {
		se.Err = err[0]
	}

	return se
}

func (svc *Error) Data(items map[string]any) *Error {
	svc.Detail = items
	return svc
}

func (svc *Error) SetErr(err string) *Error {
	svc.Err = errors.New(err)
	return svc
}

// "default error" methods

func (svc *Error) Error() string {
	return svc.Err.Error()
}

func (svc *Error) Is(target error) bool {
	targetErr, ok := target.(*Error)
	return ok && svc.Msg == targetErr.Msg
}

func (svc *Error) Unwrap() error {
	return svc.Err
}

//

func EvalTxErr(txErr error) (err *Error) {
	customErr, ok := txErr.(*Error)
	if ok == true {
		err = customErr
		return
	}

	err = Failed
	return
}

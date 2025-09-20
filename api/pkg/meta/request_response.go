package meta

import (
	"github.com/labstack/echo/v4"
	"microservice/internal/adapter/locale"
	st "microservice/pkg/meta/status"
	"net/http"
)

type (
	Result struct {
		Status  int         `json:"status" example:"200"`
		Message string      `json:"message" example:"response message"`
		Data    interface{} `json:"data,omitempty"`
		Error   string      `json:"error,omitempty"`
	}

	Response struct {
		Result
		l   locale.ILocale
		ctx echo.Context
	}
)

func Resp(c echo.Context, l locale.ILocale) *Response {
	resp := &Response{}
	resp.l = l
	resp.ctx = c
	return resp
}

func (r *Response) Status(status st.HttpMappedStatus) *Response {
	r.Result.Status = st.MappedStatuses[status]
	r.Result.Message = r.l.Get(string(status))
	return r
}

func (r *Response) Msg(message string) *Response {
	if len(message) > 0 {
		r.Result.Message = message
	}

	return r
}

func (r *Response) Data(data interface{}) *Response {
	r.Result.Data = data
	return r
}

func (r *Response) Err(err error) *Response {
	r.Result.Error = err.Error()
	return r
}

func (r *Response) ServiceErr(err error) *Response {

	if se, ok := err.(*Error); ok == true {
		r.Result.Status = st.MappedStatuses[se.Msg]
		r.Result.Message = r.l.Get(string(se.Msg))

		if se.Err != nil {
			r.Result.Error = se.Err.Error()
		}

		if len(se.Detail) > 0 {
			r.Result.Data = se.Detail
		}
	} else {
		r.Result.Status = http.StatusBadRequest
		r.Result.Message = r.l.Get("resp_fail")
	}

	return r
}

func (r *Response) Json() error {
	if r.Result.Status == 0 {
		r.Status(st.Success)
	}
	return r.ctx.JSON(r.Result.Status, r.Result)
}

package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/safe-homie/backend/internal/constant"
)

type Response struct {
	Data  any
	Code  int
	Error any
}

type SuccessResponseWrapper struct {
	Success bool `json:"success" default:"true"`
	Code    int  `json:"code"`
	Data    any  `json:"data"`
}

type ErrorResponseWrapper struct {
	Success bool `json:"success" default:"false"`
	Code    int  `json:"code"`
	Error   any  `json:"error"`
}

func NewErrorResponse(code int, err any) *Response {
	return &Response{
		Code:  code,
		Error: err,
	}
}

func NewDataResponse(code int, data any) *Response {
	return &Response{
		Code: code,
		Data: data,
	}
}

func Wrap(fn func(c echo.Context) *Response) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := fn(c)
		payload := map[string]interface{}{
			constant.SuccessField: true,
			constant.CodeField:    http.StatusOK,
		}
		if res.Code > 0 {
			payload[constant.CodeField] = res.Code
		}
		if res.Error != nil {
			payload[constant.SuccessField] = false
			payload[constant.ErrorField] = res.Error
		} else {
			payload[constant.DataField] = res.Data
		}
		return c.JSON(res.Code, payload)
	}
}

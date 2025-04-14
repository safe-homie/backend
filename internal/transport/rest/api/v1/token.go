package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/safe-homie/backend/internal/domain"
)

func (a *APIV1) SaveToken(ctx echo.Context) *Response {
	var req domain.SaveTokenRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorInvalidRequestBody
	}
	_, err := a.notifyService.SaveToken(&req)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to save push token"+err.Error())
	}
	return NewDataResponse(http.StatusOK, req)
}

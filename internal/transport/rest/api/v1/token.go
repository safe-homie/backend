package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/safe-homie/backend/internal/domain"
)

// @Summary		Save Expo push token
// @Description	Update Expo push token
// @Tags			token
// @Accept			json
// @Produce		json
// @Param			request	body		domain.SaveTokenRequest	true	"save token request body"
// @Success		200		{object}	SuccessResponseWrapper{data=domain.SaveTokenRequest}
// @Failure		400		{object}	ErrorResponseWrapper{error=string}
// @Failure		500		{object}	ErrorResponseWrapper{error=string}
// @Router			/notify/token [post]
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

package auth

import (
	"net/http"

	"github.com/anousoneFS/clean-architecture/helper"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	authUsecase AuthUsecase
}

func NewHandler(app *echo.Echo, uc AuthUsecase) {
	h := Auth{
		authUsecase: uc,
	}
	api := app.Group("/v1")
	api.POST("/login", h.login)
}

func (h Auth) login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(401, helper.GetErrMessage(err))
	}
	ctx := c.Request().Context()
	res, err := h.authUsecase.Login(ctx, req)
	if err != nil {
		return c.JSON(helper.GetHttpStatus(err), helper.GetErrMessage(err))
	}
	return c.JSON(http.StatusOK, res)
}

package user

import (
	"net/http"

	"github.com/anousoneFS/clean-architecture/helper"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserUC UserUsecase
}

func NewHandler(e *echo.Echo, userUC UserUsecase) {
	h := &UserHandler{
		UserUC: userUC,
	}
	api := e.Group("/v1/users")
	api.GET("", h.listUsers)
	api.POST("", h.createUser)
	api.GET("/:id", h.get)
}

func (h *UserHandler) listUsers(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.UserUC.List(ctx)
	if err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, res)
}

func (h *UserHandler) createUser(c echo.Context) error {
	var req User
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "bad request"})
	}
	ctx := c.Request().Context()
	if err := h.UserUC.Create(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "error"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "created"})
}

func (h *UserHandler) get(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	res, err := h.UserUC.Get(ctx, id)
	if err != nil {
		return c.JSON(helper.GetHttpStatus(err), helper.GetErrMessage(err))
	}
	return c.JSON(http.StatusOK, res)
}

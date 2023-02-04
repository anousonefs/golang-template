package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userUC userUC
}

func NewHandler(e *echo.Echo, userUC userUC) {
	h := &userHandler{
		userUC: userUC,
	}
	api := e.Group("/v1/users")
	api.GET("", h.listUsers)
	api.POST("", h.createUser)
}

func (h *userHandler) listUsers(c echo.Context) error {
	ctx := c.Request().Context()
	resp, err := h.userUC.List(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "error"})
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *userHandler) createUser(c echo.Context) error {
	var req User
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "bad request"})
	}
	ctx := c.Request().Context()
	if err := h.userUC.Create(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "error"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "created"})
}

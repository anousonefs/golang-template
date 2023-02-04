package user

import (
	"net/http"

	helper "github.com/anousoneFS/clean-architecture/helper"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/encoding/protojson"
)

type UserHandler struct {
	UserUC *UserUC
}

func NewHandler(e *echo.Echo, userUC *UserUC) {
	h := &UserHandler{
		UserUC: userUC,
	}
	api := e.Group("/v1/users")
	api.GET("", h.ListUsers)
	api.POST("", h.CreateUser)
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	ctx := c.Request().Context()
	resp, err := h.UserUC.List(ctx)
	if err != nil {
		hs := helper.HttpStatusPbFromRPC(helper.GRPCStatusFromErr(err))
		b, _ := protojson.Marshal(hs)
		return c.JSONBlob(int(hs.Error.Code), b)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
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

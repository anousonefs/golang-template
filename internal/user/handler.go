package user

import (
	helper "github.com/anousoneFS/clean-architecture/helper"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
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
		hs := helper.HttpStatusPbFromRPC(helper.GRPCStatusFromErr(err))
		b, _ := protojson.Marshal(hs)
		return c.JSONBlob(int(hs.Error.Code), b)
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

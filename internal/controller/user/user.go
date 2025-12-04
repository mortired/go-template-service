package user

import (
	"net/http"
	"users/internal/model/user"
	"users/internal/service"

	response "github.com/mortired/appsap-response"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	userService service.User
}

func New(
	service service.User,
) *Controller {

	return &Controller{
		userService: service,
	}
}

func (c *Controller) ListUsers(ctx echo.Context) error {
	var filter user.Filter
	if err := ctx.Bind(&filter); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.InvalidRequest("Invalid request parameters: "+err.Error(), ctx.Request().URL.Path))
	}

	if err := filter.Validate(); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.ValidationError(err, ctx.Request().URL.Path))
	}

	users, err := c.userService.ListUsers(filter)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.InternalError(err.Error(), ctx.Request().URL.Path))
	}

	return ctx.JSON(http.StatusOK, users)
}

func (c *Controller) CreateUser(ctx echo.Context) error {
	var req user.CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.InvalidRequest("Invalid data format: "+err.Error(), ctx.Request().URL.Path))
	}

	if err := req.Validate(); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.ValidationError(err, ctx.Request().URL.Path))
	}

	createdUser, err := c.userService.CreateUser(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.InternalError(err.Error(), ctx.Request().URL.Path))
	}

	return ctx.JSON(http.StatusCreated, createdUser)
}

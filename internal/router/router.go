package router

import (
	"users/internal/controller/user"

	"github.com/labstack/echo/v4"
)

const apiV1 = "/api/v1"

type Route struct {
	Method      string
	Path        string
	Handler     echo.HandlerFunc
	Description string
}

func SetupRoutes(e *echo.Echo, c *user.Controller) {
	apiGroup := e.Group(apiV1)

	UserRoutes(apiGroup, c)
}

func UserRoutes(e *echo.Group, c *user.Controller) {
	routes := []Route{
		{
			Method:      "GET",
			Path:        "/users",
			Handler:     c.ListUsers,
			Description: "List users (HMAC protected)",
		},
		{
			Method:      "POST",
			Path:        "/users",
			Handler:     c.CreateUser,
			Description: "Create user (HMAC protected)",
		},
	}

	for _, route := range routes {
		e.Add(route.Method, route.Path, route.Handler)
	}
}

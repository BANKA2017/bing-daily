package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ApiTemplate[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Version string `json:"version"`
}

func apiTemplate[T any](code int, message string, data T, version string) ApiTemplate[T] {
	var template ApiTemplate[T]
	template.Code = code
	template.Message = message
	template.Data = data
	template.Version = version
	return template
}

func echoReject(c echo.Context) error {
	var response = apiTemplate(403, "Invalid Request", EmptyObject, "bd@api")
	return c.JSON(http.StatusForbidden, response)
}

func echoFavicon(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func echoRobots(c echo.Context) error {
	return c.String(http.StatusOK, "User-agent: *\nDisallow: /*")
}

var EmptyObject = make(map[string]interface{})
var EmptyArray = make([]struct{}, 0)

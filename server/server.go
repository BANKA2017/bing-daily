package server

import (
	"github.com/labstack/echo/v4"
)

func Server(host string) {
	e := echo.New()
	//e.Use(middleware.Logger())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("X-Powered-By", "bd@api")
			c.Response().Header().Add("Access-Control-Allow-Origin", "*")
			return next(c)
		}
	})
	e.Any("/*", echoReject)
	e.Any("/favicon.ico", echoFavicon)
	e.Any("/robots.txt", echoRobots)

	e.GET("/v1/data/list/", ImageList)
	e.GET("/docs", Docs)

	e.Logger.Fatal(e.Start(host))
}

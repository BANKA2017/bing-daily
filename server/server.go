package main

import (
	"fmt"
	"net/http"
	"strconv"

	"bing.image.nest.moe/m/v2/types"

	"bing.image.nest.moe/m/v2/dbio"
	"github.com/labstack/echo/v4"
)

func apiTemplate[T any](code int, message string, data T, version string) types.ApiTemplate {
	var template types.ApiTemplate
	template.Code = code
	template.Message = message
	template.Data = data
	template.Version = version
	return template
}

func echoReject(c echo.Context) error {
	var response = apiTemplate(403, "Invalid Request", make(map[string]interface{}, 0), "global_api")
	return c.JSON(http.StatusForbidden, response)
}

func echoFavicon(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func echoRobots(c echo.Context) error {
	return c.String(http.StatusOK, "User-agent: *\nDisallow: /*")
}

func findImageData(c echo.Context) error {
	var countN int64
	var dateN int64
	count := c.QueryParams().Get("count")
	countN, err := strconv.ParseInt(count, 10, 0)
	if err != nil {
		countN = 16
	} else if countN < 1 {
		countN = 1
	} else if countN > 100 {
		countN = 100
	}

	date := c.QueryParams().Get("date")
	dateN, err = strconv.ParseInt(date, 10, 0)
	if err != nil {
		dateN = 90000101
	}

	var imgData []types.SavedData
	err = dbio.DbioRead(&imgData)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, apiTemplate(500, "Unable to read db", make(map[string]interface{}, 0), "global_api"))
	}

	// find day
	// fmt.Println(dateN, countN)
	var tmpImgData []types.SavedData
	for k, v := range imgData {
		if v.Startdate == int(dateN) {
			tmpImgData = imgData[k : k+int(countN)]
		}
	}

	if len(tmpImgData) == 0 {
		tmpImgData = make([]types.SavedData, 0)
	}

	// TODO more...

	return c.JSON(http.StatusOK, apiTemplate(200, "OK", tmpImgData, "global_api"))
}

func main() {
	e := echo.New()
	//e.Use(middleware.Logger())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("X-Powered-By", "Bing daily")
			c.Response().Header().Add("Access-Control-Allow-Methods", "*")
			c.Response().Header().Add("Access-Control-Allow-Credentials", "true")
			c.Response().Header().Add("Access-Control-Allow-Origin", "*")
			return next(c)
		}
	})
	e.GET("/v1/data/list/", findImageData)
	e.Any("/*", echoReject)
	e.Any("/favicon.ico", echoFavicon)
	e.Any("/robots.txt", echoRobots)

	e.Logger.Fatal(e.Start(":1323"))
}

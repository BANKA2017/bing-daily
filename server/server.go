package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BANKA2017/bing-daily/types"
	"golang.org/x/exp/slices"

	"github.com/BANKA2017/bing-daily/dbio"
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
	var response = apiTemplate(403, "Invalid Request", make(map[string]interface{}, 0), "bing")
	return c.JSON(http.StatusForbidden, response)
}

func echoFavicon(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func echoRobots(c echo.Context) error {
	return c.String(http.StatusOK, "User-agent: *\nDisallow: /*")
}

func isValidDate(_date int64) bool {
	year := _date / 10000
	month := (_date % 10000) / 100
	date := _date % 100

	isLeapYear := (year%4 == 0 && year%100 != 0) || (year%400 != 0)

	if month < 0 || month > 12 || date < 0 || date > 31 || (slices.Contains([]int64{4, 6, 9, 11}, month) && date > 30) || month == 2 && ((isLeapYear && date > 29) || (!isLeapYear && date > 28)) {
		return false
	} else {
		return true
	}
}

func findImgIndex(imgData []types.SavedData, date int64, lIndex, rIndex int) int {
	var cIndex = (rIndex-lIndex)/2 + lIndex

	if !isValidDate(date) || int(date) < imgData[0].Startdate || int(date) > imgData[len(imgData)-1].Startdate {
		return -1
	}

	if imgData[cIndex].Startdate > int(date) {
		return findImgIndex(imgData, date, lIndex, cIndex)
	} else if imgData[cIndex].Startdate == int(date) {
		return cIndex
	} else if cIndex == lIndex && imgData[cIndex].Startdate != int(date) {
		return -1
	} else {
		return findImgIndex(imgData, date, cIndex, rIndex)
	}
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
		return c.JSON(http.StatusInternalServerError, apiTemplate(500, "Unable to read db", make(map[string]interface{}, 0), "bing"))
	}

	// find day
	// fmt.Println(dateN, countN)
	tmpIndex := findImgIndex(imgData, dateN, 0, len(imgData)-1)
	var tmpImgData []types.SavedData
	if tmpIndex > -1 {
		rIndex := tmpIndex + int(countN+1)
		if rIndex > len(imgData) {
			rIndex = len(imgData)
		}
		tmpImgData = imgData[tmpIndex:rIndex]
	}

	if len(tmpImgData) == 0 {
		tmpImgData = make([]types.SavedData, 0)
	}

	var imgList types.ApiImgList
	imgList.More = len(tmpImgData) == int(countN)+1
	imgList.Image = tmpImgData
	if imgList.More {
		imgList.Image = imgList.Image[:len(tmpImgData)-1]
	}

	return c.JSON(http.StatusOK, apiTemplate(200, "OK", imgList, "bing"))
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

package server

import (
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/BANKA2017/bing-daily/bing"
	"github.com/BANKA2017/bing-daily/dbio"
	"github.com/BANKA2017/bing-daily/dbio/model"
	"github.com/labstack/echo/v4"
)

type ApiImgList struct {
	Image []*bing.SavedData2 `json:"image"`
	More  bool               `json:"more"`
}

type ImageDTO struct {
	Count int64  `query:"count"`
	Date  int64  `query:"date"`
	Color string `query:"color"`
	Mkt   string `query:"mkt"`
}

func ImageList(c echo.Context) error {
	q := new(ImageDTO)
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusInternalServerError, apiTemplate(500, "Invalid Query", EmptyArray, "bd@api"))
	}

	var countN = q.Count
	var dateN = q.Date
	if countN < 1 {
		countN = 1
	} else if countN > 100 {
		countN = 100
	}

	mkt := strings.ToUpper(q.Mkt)
	if !slices.Contains(bing.ValidMkt, mkt) && mkt != "ROW" {
		mkt = "ZH-CN"
	} else if mkt == "EN-NZ" || mkt == "EN-AU" {
		mkt = "ROW"
	}

	mktDate := bing.LatestDate[mkt]

	// find day
	if dateN > mktDate {
		return c.JSON(http.StatusOK, apiTemplate(404, "Out of range", EmptyArray, "bd@api"))
	}

	if dateN == 0 {
		dateN = mktDate
	}

	var DBImg []*model.Img2
	// var DBColor []*model.Color

	if err := dbio.GormMemCacheDB.R.Model(&model.Img2{}).Where("date >= ? AND market = ?", dateN, mkt).Limit(int(countN) + 1).Order("date").Find(&DBImg).Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, apiTemplate(500, "Failed", EmptyArray, "bd@api"))
	}

	slices.Reverse(DBImg)

	more := false
	if len(DBImg) >= int(countN+1) {
		more = true
		DBImg = DBImg[1:]
	}

	var SavedData = make([]*bing.SavedData2, len(DBImg))

	for i := 0; i < len(DBImg); i++ {
		d := DBImg[i]
		SavedData[i] = &bing.SavedData2{
			Blurhash:    d.Blurhash,
			Color:       strings.Split(d.Color, ","),
			Height:      int(d.Height),
			Width:       int(d.Width),
			Title:       d.Title,
			Headline:    d.Headline,
			Description: d.Description,
			// QuickFact:    d.QuickFact,
			Copyright:    d.Copyright,
			TriviaUrl:    d.TriviaURL,
			BackstageUrl: d.BackstageURL,
			Name:         d.Name,
			Market:       d.Market,
			Hash:         d.Hash,
			Url:          d.URL,
			Date:         int(d.Date),
		}
	}

	var imgList ApiImgList
	imgList.More = more
	imgList.Image = SavedData

	return c.JSON(http.StatusOK, apiTemplate(200, "OK", imgList, "bd@api"))
}

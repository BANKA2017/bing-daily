package bing

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/BANKA2017/bing-daily/dbio"
)

var LatestDate = InitLatestDate()

func InitLatestDate() map[string]int64 {
	var ld = make(map[string]int64, len(ValidMkt))

	for _, mkt := range ValidMkt {
		ld[mkt] = 0
	}

	return ld
}

type SavedData struct {
	Blurhash      string   `json:"blurhash"`
	Bot           int      `json:"bot"`
	Color         []string `json:"color"`
	Copyright     string   `json:"copyright"`
	Copyrightlink string   `json:"copyrightlink"`
	Drk           int      `json:"drk"`
	Height        int      `json:"height"`
	Hs            string   `json:"hs"`
	Hsh           string   `json:"hsh"`
	Quiz          string   `json:"quiz"`
	Startdate     int      `json:"startdate"`
	Title         string   `json:"title"`
	Top           int      `json:"top"`
	URL           string   `json:"url"`
	Urlbase       string   `json:"urlbase"`
	Width         int      `json:"width"`
	Wp            int      `json:"wp"`
}

type SavedData2 struct {
	// img meta data
	Blurhash string   `json:"blurhash"`
	Color    []string `json:"color"`
	Height   int      `json:"height"`
	Width    int      `json:"width"`

	// story data
	Title       string `json:"title"`
	Headline    string `json:"headline"`
	Description string `json:"description"`
	// QuickFact   string `json:"quick_fact"`
	Copyright string `json:"copyright"`

	// bing links
	TriviaUrl    string `json:"trivia_url"` // quiz
	BackstageUrl string `json:"backstage_url"`

	// data
	// OHR.{Name}_{Market}{Hash}
	Name   string `json:"name"`
	Market string `json:"market"`
	Hash   string `json:"hash"`
	Url    string `json:"url"`
	Date   int    `json:"date"`
}

type BingImageInfoImage struct {
	Startdate     string `json:"startdate"`
	Fullstartdate string `json:"fullstartdate"`
	Enddate       string `json:"enddate"`
	URL           string `json:"url"`
	Urlbase       string `json:"urlbase"`
	Copyright     string `json:"copyright"`
	Copyrightlink string `json:"copyrightlink"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Desc2         string `json:"desc2,omitempty"`
	Desc3         string `json:"desc3,omitempty"`
	Desc4         string `json:"desc4,omitempty"`
	Desc5         string `json:"desc5,omitempty"`
	Desc6         string `json:"desc6,omitempty"`
	Desc7         string `json:"desc7,omitempty"`
	Desc8         string `json:"desc8,omitempty"`
	Desc9         string `json:"desc9,omitempty"`
	Desc10        string `json:"desc10,omitempty"`
	Quiz          string `json:"quiz"`
	Wp            bool   `json:"wp"`
	Hsh           string `json:"hsh"`
	Drk           int    `json:"drk"`
	Top           int    `json:"top"`
	Bot           int    `json:"bot"`
	Hs            []any  `json:"hs"`
}

type BingImageInfo struct {
	Market struct {
		Mkt string `json:"mkt"`
	} `json:"market"`
	Images []BingImageInfoImage `json:"images"`
}

type BingImageInfoImage2 struct {
	ImageContent struct {
		Description string `json:"Description,omitempty"`
		Image       struct {
			URL          string `json:"Url,omitempty"`
			Wallpaper    string `json:"Wallpaper,omitempty"`
			Downloadable bool   `json:"Downloadable,omitempty"`
		} `json:"Image,omitempty"`
		Headline   string `json:"Headline,omitempty"`
		Title      string `json:"Title,omitempty"`
		Copyright  string `json:"Copyright,omitempty"`
		SocialGood any    `json:"SocialGood,omitempty"`
		MapLink    struct {
			URL  string `json:"Url,omitempty"`
			Link string `json:"Link,omitempty"`
		} `json:"MapLink,omitempty"`
		QuickFact    json.RawMessage `json:"QuickFact,omitempty"`
		TriviaURL    string          `json:"TriviaUrl,omitempty"`
		BackstageURL string          `json:"BackstageUrl,omitempty"`
		TriviaID     string          `json:"TriviaId,omitempty"`
	} `json:"ImageContent,omitempty"`
	Ssd            string `json:"Ssd,omitempty"`
	Name           string `json:"Name,omitempty"`
	Market         string `json:"Market,omitempty"`
	Hash           string `json:"Hash,omitempty"`
	FullDateString string `json:"FullDateString,omitempty"`
}
type BingImageInfo2 struct {
	MediaContents []BingImageInfoImage2 `json:"MediaContents,omitempty"`
}

// https://www.bing.com/HPImageArchive.aspx?n=7&idx=0&format=js&pid=opal&mbl=1&desc=1&mkt=zh-cn
func GetImgInfo(locale string) (*BingImageInfo, error) {
	/// var bingImageInfo BingImageInfo
	/// readJson(ROOTPATH+"/bing.json", &bingImageInfo)
	/// return &bingImageInfo, nil
	var bingImageInfo BingImageInfo
	return dbio.FetchJson("https://www.bing.com/HPImageArchive.aspx?idx=0&n=10&format=js&pid=opal&mbl=1&desc=1&video=1&mkt="+locale, "GET", nil, nil, bingImageInfo)
}

func GetImgInfo2(locale string) (*BingImageInfo2, error) {
	/// var bingImageInfo BingImageInfo
	/// readJson(ROOTPATH+"/bing.json", &bingImageInfo)
	/// return &bingImageInfo, nil
	var bingImageInfo2 BingImageInfo2
	return dbio.FetchJson("https://www.bing.com/hp/api/model?mkt="+locale, "GET", nil, nil, bingImageInfo2)
}

// https://wallpaper.sapphire.microsoftapp.net/api/v1/app/photo/list/bing/1/10000

// var MemBingImgMetaCache = make(map[int]*SavedData)
// var MemBingImgMetaCache []*SavedData2
//
// func InitCache() error {
// 	return dbio.DbioRead(&MemBingImgMetaCache)
// }

var ValidMkt = strings.Split(strings.ToUpper("en-us,zh-cn,ja-jp,es-es,en-ca,en-au,de-de,fr-fr,it-it,en-gb,en-in,pt-br"), ",")

func PDate(_date int64) int {
	var year = int(_date / 10000)
	var month = int((_date % 10000) / 100)
	var day = int(_date % 100)

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	t = t.AddDate(0, 0, -1)
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}

func NDate(_date int64) int {
	var year = int(_date / 10000)
	var month = int((_date % 10000) / 100)
	var day = int(_date % 100)

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	t = t.AddDate(0, 0, 1)
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}

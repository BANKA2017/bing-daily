package workers

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BANKA2017/bing-daily/b2"
	"github.com/BANKA2017/bing-daily/bing"
	"github.com/BANKA2017/bing-daily/dbio"
	"github.com/BANKA2017/bing-daily/dbio/model"
	"github.com/BANKA2017/bing-daily/image2"
	"gorm.io/gorm"
)

func GetImgsWorker(B2ApplicationKeyId, B2ApplicationKey, WorkersLocale string) error {
	noB2 := B2ApplicationKey == "" || B2ApplicationKeyId == ""

	WorkersLocale = strings.ToUpper(WorkersLocale)

	if WorkersLocale == "" || !slices.Contains(bing.ValidMkt, WorkersLocale) {
		log.Println("bing-daily: Some environment variables are empty")
		noB2 = true
	}

	var b2UploadUrl *b2.B2UploadUrl
	var uploadResponse *b2.B2UploadResponse

	if !noB2 {
		b2AuthorizeAccount, err := b2.GetB2AuthorizeAccount(B2ApplicationKeyId, B2ApplicationKey)
		if err != nil {
			return err
		} else if b2AuthorizeAccount.Status != 0 {
			return errors.New(b2AuthorizeAccount.Message)
		}
		b2UploadUrl, err = b2.GetB2UploadUrl(b2AuthorizeAccount)
		if err != nil {
			return err
		} else if b2UploadUrl.Status != 0 {
			return errors.New(b2UploadUrl.Message)
		}
	}

	// if err != nil {
	// 	fmt.Println(len(savedData))
	// 	return errors.New(err)
	// }

	// if len(bing.MemBingImgMetaCache) > 0 {
	// 	startDate = bing.MemBingImgMetaCache[len(bing.MemBingImgMetaCache)-1].Date
	// }

	for _, mkt := range bing.ValidMkt {
		mktLatestDate := bing.LatestDate[mkt]

		var tmpDataList []*bing.SavedData2

		bingData, err := bing.GetImgInfo2(mkt)
		if err != nil {
			return err
		}

		for _, v := range bingData.MediaContents {
			tmpStartDate, _ := strconv.ParseInt(v.Ssd, 10, 64)
			tmpStartDate = int64(bing.PDate(tmpStartDate))
			if tmpStartDate > mktLatestDate {
				u := v.ImageContent.Image.Wallpaper
				if u == "" {
					u = v.ImageContent.Image.URL
				}

				qf, _ := v.ImageContent.QuickFact.MarshalJSON()

				tmpDataList = append(tmpDataList, &bing.SavedData2{
					Title:        v.ImageContent.Title,
					Headline:     v.ImageContent.Headline,
					Description:  v.ImageContent.Description,
					QuickFact:    string(qf),
					Copyright:    v.ImageContent.Copyright,
					TriviaUrl:    v.ImageContent.TriviaURL,
					BackstageUrl: v.ImageContent.BackstageURL,
					Name:         v.Name,
					Market:       v.Market,
					Hash:         v.Hash,
					Url:          u,
					Date:         int(tmpStartDate),
				})
			}
		}

		if len(tmpDataList) == 0 {
			fmt.Println("bing-daily: No updated", mkt)
			continue
		}

		// sort
		sort.Slice(tmpDataList[:], func(i, j int) bool {
			return tmpDataList[i].Date < tmpDataList[j].Date
		})

		bing.LatestDate[mkt] = int64(tmpDataList[0].Date)

		var DBImg []*model.Img2
		var DBColor []*model.Color

		for _, v := range tmpDataList {
			//https://www.bing.com${img.urlbase}_UHD.jpg
			//https://www.bing.com${img.urlbase}_UHD.jpg&rf=LaDigue_UHD.jpg&pid=hp&w=256&h=128&rs=1&c=4

			// /th?id=OHR.SnowdoniaDolwyddelan_ZH-CN0238391772

			var meta image2.ImageMeta

			if m := image2.ImageMetaCache.Get(v.Name); m != nil && !(!noB2 && mkt == WorkersLocale) {
				meta = m.Value()
			} else {
				urlbase := "/th?id=OHR." + v.Name + "_" + v.Market + v.Hash
				bingDailyImgBuffer, err := dbio.FetchFile("https://www.bing.com" + urlbase + "_UHD.jpg")
				if err != nil {
					return err
				}

				if !noB2 && mkt == WorkersLocale {
					uploadResponse, err = b2.UploadToB2(b2UploadUrl, bingDailyImgBuffer, "bing/"+strconv.Itoa(v.Date)+".jpg", "image2/jpeg")
					if err != nil {
						return err
					}
					fmt.Println(uploadResponse)
				}

				img, err := image2.GetImg(bingDailyImgBuffer)
				if err != nil {
					return err
				}

				meta, err = image2.GetImgMeta(img, v.Name)
				if err != nil {
					return err
				}
			}

			v.Blurhash = meta.Blurhash
			v.Color = meta.Color
			v.Width = meta.Width
			v.Height = meta.Height

			// bing.MemBingImgMetaCache = append(bing.MemBingImgMetaCache, v)

			DBImg = append(DBImg, &model.Img2{
				Blurhash:     v.Blurhash,
				Color:        strings.Join(v.Color, ","),
				Height:       int64(v.Height),
				Width:        int64(v.Width),
				Headline:     v.Headline,
				Description:  v.Description,
				QuickFact:    v.QuickFact,
				TriviaURL:    v.TriviaUrl,
				URL:          v.Url,
				Date:         int32(v.Date),
				Name:         v.Name,
				Market:       strings.ToUpper(v.Market),
				Hash:         v.Hash,
				BackstageURL: v.BackstageUrl,
				Copyright:    v.Copyright,
				Title:        v.Title,
			})

			for i := 0; i < len(v.Color); i++ {
				DBColor = append(DBColor, &model.Color{
					Date: int32(v.Date),
					Mkt:  strings.ToUpper(v.Market),
					Hex:  v.Color[i],
					RgbR: int32(meta.Rgb[i][0]),
					RgbG: int32(meta.Rgb[i][1]),
					RgbB: int32(meta.Rgb[i][2]),
					LabL: meta.Lab[i][0],
					LabA: meta.Lab[i][1],
					LabB: meta.Lab[i][2],
				})
			}
		}

		// dataBuffer, err := json.Marshal(bing.MemBingImgMetaCache)
		// if err != nil {
		// 	return err
		// }
		// dbio.DbioWrite(dataBuffer)

		// upload to b2
		// uploadResponse, err = b2.UploadToB2(b2UploadUrl, dataBuffer, "bing/bing.json", "text/json")
		// if err != nil {
		// 	return err
		// }
		// fmt.Println(uploadResponse)

		// save to db
		err = dbio.GormDB.W.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&model.Img2{}).Create(&DBImg).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.Color{}).Create(&DBColor).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}

		err = dbio.GormMemCacheDB.W.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&model.Img2{}).Create(&DBImg).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.Color{}).Create(&DBColor).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}

		time.Sleep(time.Second)
	}

	return nil
}

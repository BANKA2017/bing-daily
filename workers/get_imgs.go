package workers

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
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

func ParseCopyright(s string) (title, copyright string) {
	s = strings.TrimSpace(s)
	if i := strings.LastIndex(s, "("); i != -1 && strings.HasSuffix(s, ")") {
		title = strings.TrimSpace(s[:i])
		copyright = strings.TrimSuffix(strings.TrimSpace(s[i+1:]), ")")
	} else {
		title = s
	}
	return
}

// "/th?id=OHR.LagoonNebula_ZH-CN3890147543"
func ParseURLBase(urlbase string) (name, market, hash string, ok bool) {
	re := regexp.MustCompile(`OHR\.([^_]+)_(\D+)(\d+)`)
	m := re.FindStringSubmatch(urlbase)
	if len(m) == 4 {
		name, market, hash, ok = m[1], m[2], m[3], true
	}
	return
}

func JoinDescs(m *bing.BingImageInfoImage) string {
	var parts []string
	descs := []string{
		m.Desc, m.Desc2, m.Desc3, m.Desc4, m.Desc5,
		m.Desc6, m.Desc7, m.Desc8, m.Desc9, m.Desc10,
	}
	for _, d := range descs {
		if d != "" {
			parts = append(parts, strings.TrimSpace(d))
		}
	}
	return strings.Join(parts, "\n")
}

func BuildSearchLink(copyrightLink, date string) (string, error) {
	u, err := url.Parse(copyrightLink)
	if err != nil {
		return "", err
	}

	q := u.Query().Get("q")
	form := u.Query().Get("form")
	if q == "" {
		return "", fmt.Errorf("no q param in url")
	}
	if form == "" {
		form = "hpcapt"
	}

	result := fmt.Sprintf(`/search?q=%s&form=%s&filters=HpDate:"%s_1600"+mgzv3configlist:"BingQA_Encyclopedia_Layout"`,
		url.QueryEscape(q), form, date)
	return result, nil
}

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
		if mkt == "EN-AU" {
			mkt = "ROW"
		}
		mktLatestDate := bing.LatestDate[mkt]

		var tmpDataList []*bing.SavedData2

		bingData, err := bing.GetImgInfo(mkt)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, v := range bingData.Images {
			tmpStartDate, _ := strconv.ParseInt(v.Startdate, 10, 64)
			if tmpStartDate > mktLatestDate {
				// qf, _ := v.ImageContent.QuickFact.MarshalJSON()

				title, cpy := ParseCopyright(v.Copyright)
				name, market, hash, _ := ParseURLBase(v.Urlbase)

				desc := JoinDescs(&v)

				backstageUrl, _ := BuildSearchLink(v.Copyrightlink, strconv.Itoa(int(tmpStartDate)))
				quizLink := strings.ReplaceAll(v.Quiz, v.Startdate, strconv.Itoa(int(bing.NDate(tmpStartDate))))

				tmpDataList = append(tmpDataList, &bing.SavedData2{
					Title:       title,
					Headline:    v.Title,
					Description: desc,
					// QuickFact:    string(qf),
					Copyright:    cpy,
					TriviaUrl:    quizLink,
					BackstageUrl: backstageUrl,
					Name:         name,
					Market:       market,
					Hash:         hash,
					Url:          v.URL,
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
					log.Println(err)
					continue
				}

				if !noB2 && mkt == WorkersLocale {
					uploadResponse, err = b2.UploadToB2(b2UploadUrl, bingDailyImgBuffer, "bing/"+strconv.Itoa(v.Date)+".jpg", "image2/jpeg")
					if err != nil {
						log.Println(err)
						continue
					}
					fmt.Println(uploadResponse)
				}

				img, err := image2.GetImg(bingDailyImgBuffer)
				if err != nil {
					log.Println(err)
					continue
				}

				meta, err = image2.GetImgMeta(img, v.Name)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			v.Blurhash = meta.Blurhash
			v.Color = meta.Color
			v.Width = meta.Width
			v.Height = meta.Height

			// bing.MemBingImgMetaCache = append(bing.MemBingImgMetaCache, v)

			DBImg = append(DBImg, &model.Img2{
				Blurhash:    v.Blurhash,
				Color:       strings.Join(v.Color, ","),
				Height:      int64(v.Height),
				Width:       int64(v.Width),
				Headline:    v.Headline,
				Description: v.Description,
				// QuickFact:    v.QuickFact,
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
			log.Println(err)
			continue
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
			log.Println(err)
			continue
		}

		fmt.Println("bing-daily: Done", mkt)

		time.Sleep(time.Second)
	}

	return nil
}

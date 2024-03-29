package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/BANKA2017/bing-daily/types"

	"github.com/BANKA2017/bing-daily/dbio"
	Blurhash "github.com/bbrks/go-blurhash"
	m "github.com/ericpauley/go-quantize/quantize"
	"github.com/spf13/viper"
)

func padStart(str string, length int, pad string) string {
	if len(str) >= length {
		return str
	}
	return strings.Repeat(pad, length-len(str)) + str
}

func fetchJson[T any](_url string, _method string, _body []byte, _headers map[string]string, responseTemplate T) (*T, error) {
	var req *http.Request
	var err error
	if _body != nil {
		body := &bytes.Buffer{}
		writer := io.Writer(body)
		_, err = writer.Write(_body)

		if err != nil {
			panic(err)
		}
		req, err = http.NewRequest(_method, _url, body)
	} else {
		req, err = http.NewRequest(_method, _url, nil)
	}

	//if _body != nil &&  {
	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for k, v := range _headers {
		req.Header.Set(k, v)
	}
	/// fmt.Println(req.Header.Values("Content-Type"))
	client := &http.Client{}
	/// fmt.Println(req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	/// fmt.Println(string(response))

	if err = json.Unmarshal(response, &responseTemplate); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &responseTemplate, err
}

func fetchFile(_url string) ([]byte, error) {
	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func getImgInfo(locale string) (*types.BingImageInfo, error) {
	/// var bingImageInfo BingImageInfo
	/// readJson(ROOTPATH+"/bing.json", &bingImageInfo)
	/// return &bingImageInfo, nil
	var bingImageInfo types.BingImageInfo
	return fetchJson("https://www.bing.com/HPImageArchive.aspx?idx=0&n=10&format=js&mkt="+locale, "GET", nil, nil, bingImageInfo)
}

func getB2AuthorizeAccount(applicationKeyId string, applicationKey string) (*types.B2AuthorizeAccount, error) {
	/// var b2AuthorizeAccount B2AuthorizeAccount
	/// readJson(ROOTPATH+"/b2.json", &b2AuthorizeAccount)
	/// return &b2AuthorizeAccount, nil
	var b2AuthorizeAccount types.B2AuthorizeAccount
	return fetchJson("https://api.backblazeb2.com/b2api/v2/b2_authorize_account", "GET", nil, map[string]string{
		"Authorization": "Basic" + base64.StdEncoding.EncodeToString([]byte(applicationKeyId+":"+applicationKey)),
	}, b2AuthorizeAccount)
}

func getB2UploadUrl(b2AuthorizeAccount *types.B2AuthorizeAccount) (*types.B2UploadUrl, error) {
	/// var b2UploadUrl B2UploadUrl
	/// readJson(ROOTPATH+"/b2upload.json", &b2UploadUrl)
	/// return &b2UploadUrl, nil
	var b2UploadUrl types.B2UploadUrl
	return fetchJson(b2AuthorizeAccount.APIURL+"/b2api/v2/b2_get_upload_url?bucketId="+b2AuthorizeAccount.Allowed.BucketID, "GET", nil, map[string]string{
		"Authorization": b2AuthorizeAccount.AuthorizationToken,
	}, b2UploadUrl)
}

func uploadToB2(b2AuthorizeAccount *types.B2AuthorizeAccount, b2UploadUrl *types.B2UploadUrl, fileBuffer []byte, fileName string, contentType string) (*types.B2UploadResponse, error) {
	var uploadResponse types.B2UploadResponse
	/// readJson(ROOTPATH+"/b2response.json", &uploadResponse)
	/// return &uploadResponse, nil
	_sha1 := sha1.Sum(fileBuffer)
	return fetchJson(b2UploadUrl.UploadURL, "POST", fileBuffer, map[string]string{
		"Authorization":     b2UploadUrl.AuthorizationToken,
		"Content-Type":      contentType,
		"X-Bz-File-Name":    fileName,
		"Content-Length":    strconv.Itoa(len(fileBuffer)),
		"X-Bz-Content-Sha1": hex.EncodeToString(_sha1[:]),
	}, uploadResponse)
}

type ImageMeta struct {
	Width    int      `json:"width"`
	Height   int      `json:"height"`
	Color    []string `json:"color"`
	Blurhash string   `json:"blurhash"`
}

func getImageMeta(imgBuffer, smallImgBuffer []byte) (ImageMeta, error) {
	var ImageMeta ImageMeta
	imageMeta, _, err := image.Decode(bytes.NewReader(imgBuffer))
	if err != nil {
		return ImageMeta, err
	}

	imgSize := imageMeta.Bounds().Size()

	ImageMeta.Width = imgSize.X
	ImageMeta.Height = imgSize.Y

	q := m.MedianCutQuantizer{}
	p := q.Quantize(make([]color.Color, 0, 5), imageMeta)
	for _, color := range p {
		tmpColorR, tmpColorG, tmpColorB, _ := color.RGBA()
		ImageMeta.Color = append(ImageMeta.Color, padStart(strconv.FormatInt(int64(tmpColorR>>8), 16), 2, "0")+padStart(strconv.FormatInt(int64(tmpColorG>>8), 16), 2, "0")+padStart(strconv.FormatInt(int64(tmpColorB>>8), 16), 2, "0"))
	}

	smallImageMeta, _, err := image.Decode(bytes.NewReader(smallImgBuffer))
	if err != nil {
		return ImageMeta, err
	}
	hash, err := Blurhash.Encode(4, 4, smallImageMeta)

	if err != nil {
		return ImageMeta, err
	}

	ImageMeta.Blurhash = hash

	return ImageMeta, nil
}

func main() {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
	if err := viper.BindEnv("B2_APPLICATION_KEY_ID"); err != nil {
		panic(err)
	}
	B2_APPLICATION_KEY_ID := viper.GetString("B2_APPLICATION_KEY_ID")
	if err := viper.BindEnv("B2_APPLICATION_KEY"); err != nil {
		panic(err)
	}
	B2_APPLICATION_KEY := viper.GetString("B2_APPLICATION_KEY")
	if err := viper.BindEnv("WORKERS_LOCALE"); err != nil {
		panic(err)
	}
	WORKERS_LOCALE := viper.GetString("WORKERS_LOCALE")

	if B2_APPLICATION_KEY == "" || B2_APPLICATION_KEY_ID == "" || WORKERS_LOCALE == "" {
		panic("bing-daily: Some environment variables are empty")
	}

	var savedData []types.SavedData
	_ = dbio.DbioRead(&savedData)

	// if err != nil {
	// 	fmt.Println(len(savedData))
	// 	panic(err)
	// }
	var startDate = 0
	if len(savedData) > 0 {
		startDate = savedData[len(savedData)-1].Startdate
	}

	var tmpDataList []types.SavedData

	bingData, err := getImgInfo(WORKERS_LOCALE)
	if err != nil {
		panic(err)
	}

	for _, v := range bingData.Images {
		tmpStartDate, _ := strconv.ParseInt(v.Startdate, 10, 0)
		if tmpStartDate > int64(startDate) {
			var Wp int
			if v.Wp {
				Wp = 1
			} else {
				Wp = 0
			}

			if err != nil {
				panic(err)
			}
			tmpHs, err := json.Marshal(v.Hs)
			if err != nil {
				panic(err)
			}
			tmpDataList = append(tmpDataList, types.SavedData{
				Startdate:     int(tmpStartDate),
				URL:           v.URL,
				Urlbase:       v.Urlbase,
				Copyright:     v.Copyright,
				Copyrightlink: v.Copyrightlink,
				Title:         v.Title,
				Quiz:          v.Quiz,
				Wp:            Wp,
				Hsh:           v.Hsh,
				Drk:           v.Drk,
				Top:           v.Top,
				Bot:           v.Bot,
				Hs:            string(tmpHs), //JSON.stringify(v.hs)
			})
		}
	}

	// sort
	sort.Slice(tmpDataList[:], func(i, j int) bool {
		return tmpDataList[i].Startdate < tmpDataList[j].Startdate
	})

	if len(tmpDataList) == 0 {
		fmt.Println("bing-daily: No updated")
		os.Exit(0)
	}

	b2AuthorizeAccount, err := getB2AuthorizeAccount(B2_APPLICATION_KEY_ID, B2_APPLICATION_KEY)
	if err != nil {
		panic(err)
	} else if b2AuthorizeAccount.Status != 0 {
		panic(b2AuthorizeAccount.Message)
	}
	b2UploadUrl, err := getB2UploadUrl(b2AuthorizeAccount)
	if err != nil {
		panic(err)
	} else if b2UploadUrl.Status != 0 {
		panic(b2UploadUrl.Message)
	}

	var uploadResponse *types.B2UploadResponse

	for _, v := range tmpDataList {
		//https://www.bing.com${img.urlbase}_UHD.jpg
		//https://www.bing.com${img.urlbase}_UHD.jpg&rf=LaDigue_UHD.jpg&pid=hp&w=256&h=128&rs=1&c=4
		bingDailyImgBuffer, err := fetchFile("https://www.bing.com" + v.Urlbase + "_UHD.jpg")
		if err != nil {
			panic(err)
		}
		bingDailyImgSmallBuffer, err := fetchFile("https://www.bing.com" + v.Urlbase + "_UHD.jpg&rf=LaDigue_UHD.jpg&pid=hp&w=256&h=128&rs=1&c=4")
		if err != nil {
			panic(err)
		}

		uploadResponse, err = uploadToB2(b2AuthorizeAccount, b2UploadUrl, bingDailyImgBuffer, "bing/"+strconv.Itoa(v.Startdate)+".jpg", "image/jpeg")
		if err != nil {
			panic(err)
		}
		fmt.Println(uploadResponse)
		meta, err := getImageMeta(bingDailyImgBuffer, bingDailyImgSmallBuffer)
		if err != nil {
			panic(err)
		}

		v.Blurhash = meta.Blurhash
		v.Color = meta.Color
		v.Width = meta.Width
		v.Height = meta.Height

		savedData = append(savedData, v)
	}

	dataBuffer, err := json.Marshal(savedData)
	if err != nil {
		panic(err)
	}
	dbio.DbioWrite(dataBuffer)

	// upload to b2
	uploadResponse, err = uploadToB2(b2AuthorizeAccount, b2UploadUrl, dataBuffer, "bing/bing.json", "text/json")
	if err != nil {
		panic(err)
	}
	fmt.Println(uploadResponse)
}

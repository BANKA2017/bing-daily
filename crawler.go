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
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	Blurhash "github.com/bbrks/go-blurhash"
	"github.com/spf13/viper"
)

var ROOTPATH = ""

type BingImageInfoImage struct {
	Startdate     string `json:"startdate"`
	Fullstartdate string `json:"fullstartdate"`
	Enddate       string `json:"enddate"`
	URL           string `json:"url"`
	Urlbase       string `json:"urlbase"`
	Copyright     string `json:"copyright"`
	Copyrightlink string `json:"copyrightlink"`
	Title         string `json:"title"`
	Quiz          string `json:"quiz"`
	Wp            bool   `json:"wp"`
	Hsh           string `json:"hsh"`
	Drk           int    `json:"drk"`
	Top           int    `json:"top"`
	Bot           int    `json:"bot"`
	Hs            []any  `json:"hs"`
}

type BingImageInfo struct {
	Images   []BingImageInfoImage `json:"images"`
	Tooltips struct {
		Loading  string `json:"loading"`
		Previous string `json:"previous"`
		Next     string `json:"next"`
		Walle    string `json:"walle"`
		Walls    string `json:"walls"`
	} `json:"tooltips"`
}

type B2AuthorizeAccount struct {
	Code                    string `json:"code,omitempty"`
	Message                 string `json:"message,omitempty"`
	Status                  int    `json:"status,omitempty"`
	AbsoluteMinimumPartSize int    `json:"absoluteMinimumPartSize,omitempty"`
	AccountID               string `json:"accountId,omitempty"`
	Allowed                 struct {
		BucketID     string   `json:"bucketId"`
		BucketName   string   `json:"bucketName"`
		Capabilities []string `json:"capabilities"`
		NamePrefix   any      `json:"namePrefix"`
	} `json:"allowed,omitempty"`
	APIURL              string `json:"apiUrl,omitempty"`
	AuthorizationToken  string `json:"authorizationToken,omitempty"`
	DownloadURL         string `json:"downloadUrl,omitempty"`
	RecommendedPartSize int    `json:"recommendedPartSize,omitempty"`
	S3APIURL            string `json:"s3ApiUrl,omitempty"`
}

type B2UploadUrl struct {
	AuthorizationToken string `json:"authorizationToken,omitempty"`
	BucketID           string `json:"bucketId,omitempty"`
	UploadURL          string `json:"uploadUrl,omitempty"`
	Code               string `json:"code,omitempty"`
	Message            string `json:"message,omitempty"`
	Status             int    `json:"status,omitempty"`
}

type B2UploadResponse struct {
	Code          string `json:"code,omitempty"`
	Message       string `json:"message,omitempty"`
	Status        int    `json:"status,omitempty"`
	AccountID     string `json:"accountId,omitempty"`
	Action        string `json:"action,omitempty"`
	BucketID      string `json:"bucketId,omitempty"`
	ContentLength int    `json:"contentLength,omitempty"`
	ContentMd5    string `json:"contentMd5,omitempty"`
	ContentSha1   string `json:"contentSha1,omitempty"`
	ContentType   string `json:"contentType,omitempty"`
	FileID        string `json:"fileId,omitempty"`
	FileInfo      struct {
	} `json:"fileInfo,omitempty"`
	FileName      string `json:"fileName,omitempty"`
	FileRetention struct {
		IsClientAuthorizedToRead bool `json:"isClientAuthorizedToRead,omitempty"`
		Value                    any  `json:"value,omitempty"`
	} `json:"fileRetention,omitempty"`
	LegalHold struct {
		IsClientAuthorizedToRead bool `json:"isClientAuthorizedToRead,omitempty"`
		Value                    any  `json:"value,omitempty"`
	} `json:"legalHold,omitempty"`
	ServerSideEncryption struct {
		Algorithm any `json:"algorithm,omitempty"`
		Mode      any `json:"mode,omitempty"`
	} `json:"serverSideEncryption,omitempty"`
	UploadTimestamp int64 `json:"uploadTimestamp,omitempty"`
}

type SavedData struct {
	Blurhash      string `json:"blurhash"`
	Bot           int    `json:"bot"`
	Color         string `json:"color"`
	Copyright     string `json:"copyright"`
	Copyrightlink string `json:"copyrightlink"`
	Drk           int    `json:"drk"`
	Height        int    `json:"height"`
	Hs            string `json:"hs"`
	Hsh           string `json:"hsh"`
	Quiz          string `json:"quiz"`
	Startdate     int    `json:"startdate"`
	Title         string `json:"title"`
	Top           int    `json:"top"`
	URL           string `json:"url"`
	Urlbase       string `json:"urlbase"`
	Width         int    `json:"width"`
	Wp            int    `json:"wp"`
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

func readJson[T any](path string, template *T) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, template); err != nil {
		return err
	}
	return err
}

func saveTo(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	n, _ := file.Seek(0, io.SeekEnd)
	_, err = file.WriteAt(content, n)
	defer file.Close()
	return err
}

func getImgInfo(locale string) (*BingImageInfo, error) {
	/// var bingImageInfo BingImageInfo
	/// readJson(ROOTPATH+"/bing.json", &bingImageInfo)
	/// return &bingImageInfo, nil
	var bingImageInfo BingImageInfo
	return fetchJson("https://www.bing.com/HPImageArchive.aspx?idx=0&n=10&format=js&mkt="+locale, "GET", nil, nil, bingImageInfo)
}

func getB2AuthorizeAccount(applicationKeyId string, applicationKey string) (*B2AuthorizeAccount, error) {
	/// var b2AuthorizeAccount B2AuthorizeAccount
	/// readJson(ROOTPATH+"/b2.json", &b2AuthorizeAccount)
	/// return &b2AuthorizeAccount, nil
	var b2AuthorizeAccount B2AuthorizeAccount
	return fetchJson("https://api.backblazeb2.com/b2api/v2/b2_authorize_account", "GET", nil, map[string]string{
		"Authorization": "Basic" + base64.StdEncoding.EncodeToString([]byte(applicationKeyId+":"+applicationKey)),
	}, b2AuthorizeAccount)
}

func getB2UploadUrl(b2AuthorizeAccount *B2AuthorizeAccount) (*B2UploadUrl, error) {
	/// var b2UploadUrl B2UploadUrl
	/// readJson(ROOTPATH+"/b2upload.json", &b2UploadUrl)
	/// return &b2UploadUrl, nil
	var b2UploadUrl B2UploadUrl
	return fetchJson(b2AuthorizeAccount.APIURL+"/b2api/v2/b2_get_upload_url?bucketId="+b2AuthorizeAccount.Allowed.BucketID, "GET", nil, map[string]string{
		"Authorization": b2AuthorizeAccount.AuthorizationToken,
	}, b2UploadUrl)
}

func uploadToB2(b2AuthorizeAccount *B2AuthorizeAccount, b2UploadUrl *B2UploadUrl, imageBuffer []byte, startDate int) (*B2UploadResponse, error) {
	var uploadResponse B2UploadResponse
	/// readJson(ROOTPATH+"/b2response.json", &uploadResponse)
	/// return &uploadResponse, nil
	_sha1 := sha1.Sum(imageBuffer)
	/// fmt.Println(hex.EncodeToString(_sha1[:]))
	/// fmt.Println(map[string]string{
	/// 	"Authorization":     b2UploadUrl.AuthorizationToken,
	/// 	"Content-Type":      "image/jpeg",
	/// 	"X-Bz-File-Name":    "bing/" + "20230819" + ".jpg",
	/// 	"Content-Length":    strconv.Itoa(len(imageBuffer)),
	/// 	"X-Bz-Content-Sha1": hex.EncodeToString(_sha1[:]),
	/// })
	return fetchJson(b2UploadUrl.UploadURL, "POST", imageBuffer, map[string]string{
		"Authorization":     b2UploadUrl.AuthorizationToken,
		"Content-Type":      "image/jpeg",
		"X-Bz-File-Name":    "bing/" + strconv.Itoa(startDate) + ".jpg",
		"Content-Length":    strconv.Itoa(len(imageBuffer)),
		"X-Bz-Content-Sha1": hex.EncodeToString(_sha1[:]),
	}, uploadResponse)
}

type ImageMeta struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Color    string `json:"color"`
	Blurhash string `json:"blurhash"`
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

	var redSum float64
	var greenSum float64
	var blueSum float64

	for x := 0; x < imgSize.X; x++ {
		for y := 0; y < imgSize.Y; y++ {
			pixel := imageMeta.At(x, y)
			col := color.RGBAModel.Convert(pixel).(color.RGBA)

			redSum += float64(col.R)
			greenSum += float64(col.G)
			blueSum += float64(col.B)
		}
	}

	imgArea := float64(imgSize.X * imgSize.Y)

	redAverage := math.Round(redSum / imgArea)
	greenAverage := math.Round(greenSum / imgArea)
	blueAverage := math.Round(blueSum / imgArea)

	ImageMeta.Color = strconv.FormatInt(int64(redAverage), 16) + strconv.FormatInt(int64(greenAverage), 16) + strconv.FormatInt(int64(blueAverage), 16)

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
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	ROOTPATH = filepath.Dir(ex)

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
	if err := viper.BindEnv("ROOTPATH"); err != nil {
		fmt.Println(err)
	}
	if viper.GetString("ROOTPATH") != "" {
		ROOTPATH = viper.GetString("ROOTPATH")
	}

	var savedData []SavedData
	err = readJson(ROOTPATH+"/bing.json", &savedData)

	if err != nil {
		panic(err)
	}
	var startDate = 0
	if len(savedData) > 0 {
		startDate = savedData[len(savedData)-1].Startdate
	}

	var tmpDataList []SavedData

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
			tmpDataList = append(tmpDataList, SavedData{
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

		var uploadResponse *B2UploadResponse

		uploadResponse, err = uploadToB2(b2AuthorizeAccount, b2UploadUrl, bingDailyImgBuffer, v.Startdate)
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
	saveTo(ROOTPATH+"/bing.json", dataBuffer)
}

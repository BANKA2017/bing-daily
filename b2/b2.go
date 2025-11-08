package b2

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"strconv"

	"github.com/BANKA2017/bing-daily/dbio"
)

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

func GetB2AuthorizeAccount(applicationKeyId string, applicationKey string) (*B2AuthorizeAccount, error) {
	/// var b2AuthorizeAccount B2AuthorizeAccount
	/// readJson(ROOTPATH+"/b2.json", &b2AuthorizeAccount)
	/// return &b2AuthorizeAccount, nil
	var b2AuthorizeAccount B2AuthorizeAccount
	return dbio.FetchJson("https://api.backblazeb2.com/b2api/v2/b2_authorize_account", "GET", nil, map[string]string{
		"Authorization": "Basic" + base64.StdEncoding.EncodeToString([]byte(applicationKeyId+":"+applicationKey)),
	}, b2AuthorizeAccount)
}

func GetB2UploadUrl(b2AuthorizeAccount *B2AuthorizeAccount) (*B2UploadUrl, error) {
	/// var b2UploadUrl B2UploadUrl
	/// readJson(ROOTPATH+"/b2upload.json", &b2UploadUrl)
	/// return &b2UploadUrl, nil
	var b2UploadUrl B2UploadUrl
	return dbio.FetchJson(b2AuthorizeAccount.APIURL+"/b2api/v2/b2_get_upload_url?bucketId="+b2AuthorizeAccount.Allowed.BucketID, "GET", nil, map[string]string{
		"Authorization": b2AuthorizeAccount.AuthorizationToken,
	}, b2UploadUrl)
}

func UploadToB2(b2UploadUrl *B2UploadUrl, fileBuffer []byte, fileName string, contentType string) (*B2UploadResponse, error) {
	var uploadResponse B2UploadResponse
	/// readJson(ROOTPATH+"/b2response.json", &uploadResponse)
	/// return &uploadResponse, nil
	_sha1 := sha1.Sum(fileBuffer)
	return dbio.FetchJson(b2UploadUrl.UploadURL, "POST", fileBuffer, map[string]string{
		"Authorization":     b2UploadUrl.AuthorizationToken,
		"Content-Type":      contentType,
		"X-Bz-File-Name":    fileName,
		"Content-Length":    strconv.Itoa(len(fileBuffer)),
		"X-Bz-Content-Sha1": hex.EncodeToString(_sha1[:]),
	}, uploadResponse)
}

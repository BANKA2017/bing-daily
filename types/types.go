package types

type ApiTemplate struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Version string `json:"version"`
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

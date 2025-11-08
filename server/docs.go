package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var docsContent = `# Bing daily docs

---

## Image List

- /v1/data/list/[?count=<int64>[&date=<int64>[&mkt=<string>]]]
  - count: [1 -> 100]
  - date: yyyymmdd (20090102)
  - mkt: zh-cn/en-us/ja-jp/es-es/en-ca/en-au/de-de/fr-fr/it-it/en-gb/en-in/pt-br
`

func Docs(c echo.Context) error {
	c.Response().Header().Set("content-type", "text/txt")
	return c.String(http.StatusOK, docsContent)
}

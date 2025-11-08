package dbio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var client = http.Client{}

func FetchJson[T any](_url string, _method string, _body []byte, _headers map[string]string, responseTemplate T) (*T, error) {
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

func FetchFile(_url string) ([]byte, error) {
	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

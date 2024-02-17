package util

import (
	"bytes"
	"framework/errors"
	json "github.com/json-iterator/go"
	"net/http"
	"net/url"
	"strings"
)

type ResponseDecoder interface {
	Decode(res *http.Response, err error) error
}

type JsonResponseHandler func(res *http.Response) error

func (handler JsonResponseHandler) Decode(res *http.Response, err error) error {
	if err != nil {
		return err
	}
	return handler(res)
}

func JsonResponseDecoder(data interface{}) ResponseDecoder {
	return JsonResponseHandler(func(res *http.Response) error {
		if res.StatusCode != http.StatusOK {
			return errors.Status(res.StatusCode, res.Status)
		}

		var resData = data
		if resData == nil {
			return nil
		} else {
			return json.NewDecoder(res.Body).Decode(resData)
		}
	})
}

func HttpJsonGet(url string, resData interface{}, header ...http.Header) error {
	var req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")

	for _, h := range header {
		for k, vs := range h {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}

	return JsonResponseDecoder(resData).Decode(http.DefaultClient.Do(req))
}

func HttpJsonPost(url string, reqData, resData interface{}, header ...http.Header) error {
	var err error
	var rb []byte
	switch reqData.(type) {
	case []byte:
		rb = reqData.([]byte)
	case string:
		rb = []byte(reqData.(string))
	default:
		rb, err = json.Marshal(reqData)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(rb))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	for _, h := range header {
		for k, vs := range h {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}

	return JsonResponseDecoder(resData).Decode(http.DefaultClient.Do(req))
}

// HttpPostWhenTheContentTypeOfRequestIsUrlencodedButTheContentTypeOfResponseIsJson ??? (¯︵¯)
func HttpPostWhenTheContentTypeOfRequestIsUrlencodedButTheContentTypeOfResponseIsJson(url string, values url.Values, resData interface{}, header ...http.Header) error {
	var req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	for _, h := range header {
		for k, vs := range h {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}

	return JsonResponseDecoder(resData).Decode(http.DefaultClient.Do(req))
}

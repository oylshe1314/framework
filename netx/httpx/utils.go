package httpx

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type responseDecoder interface {
	decode(res *http.Response, err error) error
}

type jsonResponseHandler func(res *http.Response) error

func (handler jsonResponseHandler) decode(res *http.Response, err error) error {
	if err != nil {
		return err
	}
	return handler(res)
}

func jsonResponseDecoder(data interface{}) responseDecoder {
	return jsonResponseHandler(func(res *http.Response) error {
		if res.StatusCode != http.StatusOK {
			return StatusError(res.StatusCode, res.Status)
		}

		var resData = data
		if resData == nil {
			return nil
		} else {
			return json.NewDecoder(res.Body).Decode(resData)
		}
	})
}

func JsonGet(url string, resData interface{}, headers ...http.Header) error {
	var req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")

	for _, header := range headers {
		for k, vs := range header {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}

	return jsonResponseDecoder(resData).decode(http.DefaultClient.Do(req))
}

func JsonPost(url string, reqData, resData interface{}, headers ...http.Header) error {
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

	for _, header := range headers {
		for k, vs := range header {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}

	return jsonResponseDecoder(resData).decode(http.DefaultClient.Do(req))
}

// PostWhenTheContentTypeOfRequestIsUrlencodedButTheContentTypeOfResponseIsJson ??? (¯︵¯)
func PostWhenTheContentTypeOfRequestIsUrlencodedButTheContentTypeOfResponseIsJson(url string, values url.Values, resData interface{}, headers ...http.Header) error {
	var req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	for _, header := range headers {
		for k, vs := range header {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}

	return jsonResponseDecoder(resData).decode(http.DefaultClient.Do(req))
}

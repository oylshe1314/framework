package client

import (
	"bytes"
	json "github.com/json-iterator/go"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/message"
	"github.com/oylshe1314/framework/util"
	"net/http"
	"net/url"
	"time"
)

type HttpHeader http.Header

func (hh HttpHeader) Add(key, value string) {
	http.Header(hh).Add(key, value)
}

func (hh HttpHeader) Set(key, value string) {
	http.Header(hh).Set(key, value)
}

func (hh HttpHeader) Get(key string) string {
	return http.Header(hh).Get(key)
}

func (hh HttpHeader) Values(key string) []string {
	return http.Header(hh).Values(key)
}

type HttpClient struct {
	ssl     bool
	network string
	address string

	logger log.Logger

	httpUrl    *url.URL
	httpClient *http.Client
}

func (this *HttpClient) WithNetwork(network string) {
	this.network = network
}

func (this *HttpClient) WithAddress(address string) {
	this.address = address
}

func (this *HttpClient) Network() string {
	return this.address
}

func (this *HttpClient) Address() string {
	return this.address
}

func (this *HttpClient) SetLogger(logger log.Logger) {
	this.logger = logger
}

func (this *HttpClient) Init() (err error) {
	if this.logger == nil {
		this.logger = log.DefaultLogger
	}

	if len(this.network) == 0 {
		return errors.Error("'network' cannot be empty")
	}

	if len(this.address) == 0 {
		return errors.Error("'address' cannot be empty")
	}

	if this.address[len(this.address)-1] == '/' {
		this.address = this.address[:len(this.address)-1]
	}

	this.httpUrl, err = url.Parse(this.address)
	if err != nil {
		return err
	}

	this.httpClient = &http.Client{Timeout: time.Second * 30}
	return nil
}

func (this *HttpClient) Close() error {
	this.httpClient.CloseIdleConnections()
	return nil
}

func (this *HttpClient) Get(path string, query url.Values, data interface{}, headers ...HttpHeader) (*message.Reply, error) {
	reqUrl := this.httpUrl.JoinPath(path)
	reqUrl.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	for _, header := range headers {
		for key, values := range header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	var reply = &message.Reply{Data: data}
	err = util.JsonResponseDecoder(reply).Decode(this.httpClient.Do(req))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (this *HttpClient) Post(path string, query url.Values, msg, data interface{}, headers ...HttpHeader) (*message.Reply, error) {
	reqUrl := this.httpUrl.JoinPath(path)
	reqUrl.RawQuery = query.Encode()

	var buf bytes.Buffer
	var err = json.NewEncoder(&buf).Encode(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, reqUrl.String(), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	for _, header := range headers {
		for key, values := range header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	var reply = &message.Reply{Data: data}
	err = util.JsonResponseDecoder(reply).Decode(this.httpClient.Do(req))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

package http

import (
	json "github.com/json-iterator/go"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/message"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type Message struct {
	R *http.Request
	W http.ResponseWriter
}

func (this *Message) Read(v interface{}) error {
	if v == nil {
		return nil
	}

	switch this.R.Method {
	case http.MethodGet:
		return UrlValues(this.R.URL.Query()).read(v)
	case http.MethodPost:
		var contentType = this.R.Header.Get("Content-Type")
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			return err
		}

		switch mediaType {
		case "application/json":
			return json.NewDecoder(this.R.Body).Decode(v)
		case "application/x-www-form-urlencoded":
			return this.readUrlencoded(v, params)
		case "multipart/form-data":
			return this.readFormData(v, params)
		default:
			return errors.Error(http.StatusText(http.StatusUnsupportedMediaType))
		}
	default:
		return errors.Error(http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (this *Message) readUrlencoded(v interface{}, params map[string]string) error {
	body, err := io.ReadAll(this.R.Body)
	if err != nil {
		return err
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return err
	}

	return UrlValues(values).read(v)
}

func (this *Message) readFormData(v interface{}, params map[string]string) error {
	form, err := multipart.NewReader(this.R.Body, params["boundary"]).ReadForm(-1)
	if err != nil {
		return err
	}

	if len(form.Value) > 0 {
		err = UrlValues(form.Value).read(v)
		if err != nil {
			return err
		}
	}

	if len(form.File) > 0 {
		err = FileHeaders(form.File).read(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Message) Reply(v interface{}) error {
	this.W.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(this.W).Encode(message.NewReply(v))
}

func (this *Message) Write(buf []byte) error {
	_, err := this.W.Write(buf)
	return err
}

func (this *Message) File(file string) {
	http.ServeFile(this.W, this.R, file)
}

func (this *Message) Error(status int, reason string) {
	http.Error(this.W, reason, status)
}

type MessageHandler func(*Message)

func (handler MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(&Message{R: r, W: w})
}

type PostHandler MessageHandler

func (handler PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var contentType = r.Header.Get("Content-Type")

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	switch mediaType {
	case "application/json":
	case "application/x-www-form-urlencoded":
	case "multipart/form-data":
	default:
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	MessageHandler(handler).ServeHTTP(w, r)
}

type GetHandler MessageHandler

func (handler GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	MessageHandler(handler).ServeHTTP(w, r)
}

func FileHandler() GetHandler {
	return func(msg *Message) {
		var path = msg.R.URL.Path
		fi, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(msg.W, msg.R)
				return
			}

			http.Error(msg.W, err.Error(), http.StatusInternalServerError)
			return
		}

		if fi.IsDir() {
			http.NotFound(msg.W, msg.R)
			return
		}

		http.ServeFile(msg.W, msg.R, path)
	}
}

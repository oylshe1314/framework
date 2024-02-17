package server

import (
	"bufio"
	"bytes"
	"fmt"
	"framework/errors"
	. "framework/http"
	"framework/log"
	"framework/message"
	"framework/util"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

type SslConfig struct {
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

type CorsConfig struct {
	AllowOrigin      string `json:"allowOrigin"`      //Access-Control-Allow-Origin 指示响应的资源是否可以被给定的来源共享。
	AllowCredentials string `json:"allowCredentials"` //Access-Control-Allow-Credentials 指示当请求的凭证标记为 true 时，是否可以公开对该请求响应。
	AllowHeaders     string `json:"allowHeaders"`     //Access-Control-Allow-Headers 用在对预检请求的响应中，指示实际的请求中可以使用哪些 HTTP 标头。
	AllowMethods     string `json:"allowMethods"`     //Access-Control-Allow-Methods 指定对预检请求的响应中，哪些 HTTP 方法允许访问请求的资源。
	ExposeHeaders    string `json:"exposeHeaders"`    //Access-Control-Expose-Headers 通过列出标头的名称，指示哪些标头可以作为响应的一部分公开。
	MaxAge           string `json:"maxAge"`           //Access-Control-Max-Age 指示预检请求的结果能被缓存多久。
	RequestHeaders   string `json:"requestHeaders"`   //Access-Control-Request-Headers 用于发起一个预检请求，告知服务器正式请求会使用哪些 HTTP 标头。
	RequestMethod    string `json:"requestMethod"`    //Access-Control-Request-Method 用于发起一个预检请求，告知服务器正式请求会使用哪一种 HTTP 请求方法。
}

type HttpServer struct {
	Listener

	detectable bool

	ssl  *SslConfig
	cors *CorsConfig

	hs *http.Server
	sm http.ServeMux

	running bool
	logger  log.Logger
}

func (this *HttpServer) WithDetectable(detectable bool) {
	this.detectable = detectable
}

func (this *HttpServer) WithSslConfig(sslConfig *SslConfig) {
	this.ssl = sslConfig
}

func (this *HttpServer) WithCorsConfig(corsConfig *CorsConfig) {
	this.cors = corsConfig
}

func (this *HttpServer) SetLogger(logger log.Logger) {
	this.logger = logger
}

func (this *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if util.Unix() >= expiration {
		http.Error(w, "The server was expired.", http.StatusServiceUnavailable)
		return
	}

	if this.cors != nil {
		if this.cors.AllowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", this.cors.AllowOrigin)
		}
		if this.cors.AllowCredentials != "" {
			w.Header().Set("Access-Control-Allow-Credentials", this.cors.AllowCredentials)
		}
		if this.cors.AllowHeaders != "" {
			w.Header().Set("Access-Control-Allow-Headers", this.cors.AllowHeaders)
		}
		if this.cors.AllowMethods != "" {
			w.Header().Set("Access-Control-Allow-Methods", this.cors.AllowMethods)
		}
		if this.cors.ExposeHeaders != "" {
			w.Header().Set("Access-Control-Expose-Headers", this.cors.ExposeHeaders)
		}
		if this.cors.MaxAge != "" {
			w.Header().Set("Access-Control-Max-Age", this.cors.MaxAge)
		}
		if this.cors.RequestHeaders != "" {
			w.Header().Set("Access-Control-Request-Headers", this.cors.RequestHeaders)
		}
		if this.cors.RequestMethod != "" {
			w.Header().Set("Access-Control-Request-Method", this.cors.RequestMethod)
		}
		if r.Method == http.MethodOptions {
			_, _ = w.Write(nil)
			return
		}
	}

	this.sm.ServeHTTP(w, r)
}

func (this *HttpServer) ServeFile(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if fi.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, path)
}

func (this *HttpServer) Handler(pattern string, handler MessageHandler) {
	this.sm.Handle(pattern, handler)
}

func (this *HttpServer) PostHandler(pattern string, handler PostHandler) {
	this.sm.Handle(pattern, handler)
}

func (this *HttpServer) GetHandler(pattern string, handler GetHandler) {
	this.sm.Handle(pattern, handler)
}

func (this *HttpServer) FileHandler(pattern string) {
	this.sm.Handle(pattern, FileHandler(this.ServeFile))
}

func (this *HttpServer) Init() (err error) {
	if this.logger == nil {
		this.logger = log.DefaultLogger
	}

	if this.ssl != nil {
		if len(this.ssl.CertFile) == 0 || len(this.ssl.KeyFile) == 0 {
			return errors.Error("'certFile' or 'keyFile' cannot be empty when ssl is enable")
		}
	}

	err = this.Listener.Init()
	if err != nil {
		return err
	}

	this.hs = &http.Server{Handler: this, ErrorLog: log.NewNativeLogger(this.logger, log.LevelError)}

	if this.detectable {
		this.GetHandler("/server/detect", this.detect)
	}
	return nil
}

func (this *HttpServer) detect(msg *Message) {
	var ack = &message.MsgServerDetectAck{}

	ack.ProgramHash = ProgramHash
	ack.DataHash = DataHash
	ack.ConfigHash = ConfigHash
	ack.Pid = os.Getpid()
	ack.Coroutine = runtime.NumGoroutine()

	if runtime.GOOS == "linux" {
		result, err := exec.Command(fmt.Sprintf("ps -aux | grep -w %d", ack.Pid)).Output()
		if err == nil {
			var pid = util.IntegerToString(ack.Pid)
			var re = regexp.MustCompile("\\s+")
			var r = bufio.NewReader(bytes.NewReader(result))
			for {
				line, err := r.ReadString('\n')
				if err == nil {
					var ss = re.Split(line, -1)
					if ss[1] != pid {
						continue
					}

					_ = util.StringToFloat2(ss[2], &ack.CPU)
					_ = util.StringToFloat2(ss[3], &ack.Memory)
					break
				}

				if err == io.EOF {
					break
				}
			}
		}
	}

	_ = msg.Reply(ack)
}

func (this *HttpServer) Serve() (err error) {

	err = this.Listener.Listen()
	if err != nil {
		return err
	}

	this.logger.Info("HttpServer is listening on ", this.Bind())

	this.running = true
	if this.ssl != nil {
		err = this.hs.ServeTLS(this.l, this.ssl.CertFile, this.ssl.KeyFile)
		if !this.running {
			return nil
		}
	} else {
		err = this.hs.Serve(this.l)
		if !this.running {
			return nil
		}
	}

	this.running = false
	return err
}

func (this *HttpServer) Close() (err error) {
	this.running = false
	_ = this.hs.Close()
	_ = this.Listener.Close()
	return nil
}

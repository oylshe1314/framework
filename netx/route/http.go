package route

import "github.com/gin-gonic/gin"

type HttpMux struct {
	*gin.Engine
}

func NewHttpMux() *HttpMux {
	return &HttpMux{Engine: gin.New()}
}

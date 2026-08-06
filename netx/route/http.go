package route

import "github.com/gin-gonic/gin"

type HttpMux struct {
	*gin.Engine
}

func NewHttpMux(basePath string, htmlPath string) *HttpMux {
	var engine = gin.New()
	if basePath != "" {
		engine.RouterGroup = *engine.Group(basePath)
	}

	if htmlPath != "" {
		engine.LoadHTMLGlob(htmlPath)
	}

	return &HttpMux{Engine: engine}
}

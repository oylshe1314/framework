package server

import (
	"fmt"
	"github.com/oylshe1314/framework/log"
	"io"
	"runtime"
)

var baseVersion = "1.0"
var buildVersion = "0"

func Version() string {
	return fmt.Sprintf("%s.%s", baseVersion, buildVersion)
}

func printVersion(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "Version: %s, Build on %s %s/%s\n", Version(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func logVersion(logger log.Logger) {
	logger.Infof("Version: %s, Build on %s %s/%s\n", Version(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

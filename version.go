package framework

import (
	_ "embed"
	"fmt"
	"io"
	"runtime"
)

var (
	//go:embed version
	frameworkVersion string

	programVersion string = "1.0.0"
)

func FrameworkVersion() string {
	return frameworkVersion
}

func ProgramVersion(version string) {
	programVersion = version
}

func PrintVersion(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "Version: %s with fromwork-%s, Built on %s %s/%s\n", programVersion, frameworkVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

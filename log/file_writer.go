package log

import (
	"fmt"
	"io"
	"time"

	frl "github.com/lestrrat-go/file-rotatelogs"
)

func newDailyWriter(option *Option) (io.Writer, error) {
	var options []frl.Option
	if len(option.Timezone) > 0 {
		location, err := time.LoadLocation(option.Timezone)
		if err != nil {
			return nil, err
		}
		options = append(options, frl.WithLocation(location))
	}

	return frl.New(fmt.Sprintf("%s/log-%%Y-%%m-%%d.log", option.Dir), options...)
}

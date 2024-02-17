package server

import (
	"framework/util"
	"time"
)

var expiration int64

const expireDatetime = "2024-06-01 00:00:00"

func init() {
	var err error
	expiration, err = util.ParseUnitx(time.DateTime, expireDatetime)
	if err != nil {
		panic(err)
	}
}

package server

import (
	"github.com/oylshe1314/framework/util"
	"time"
)

var expiration int64

const expireDatetime = "2026-01-01 00:00:00"

func init() {
	var err error
	expiration, err = util.ParseUnix(time.DateTime, expireDatetime)
	if err != nil {
		panic(err)
	}
}

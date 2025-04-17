package server

import (
	"github.com/oylshe1314/framework/util"
	"time"
)

var expiration int64

const expireDatetime = "2025-06-01 00:00:00"

func init() {
	var err error
	expiration, err = util.ParseUnitx(time.DateTime, expireDatetime)
	if err != nil {
		panic(err)
	}
}

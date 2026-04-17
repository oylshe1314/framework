package framework

import (
	"time"

	"github.com/oylshe1314/framework/util"
)

var expiration int64

const expireDatetime = "2027-01-01 00:00:00"

func init() {
	var err error
	expiration, err = util.ParseUnix(time.DateTime, expireDatetime)
	if err != nil {
		panic(err)
	}
}

func CheckExpired() {
	if util.NowUnix() >= expiration {
		panic("The server has expired")
	}
}

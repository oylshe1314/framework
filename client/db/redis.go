package db

import (
	"github.com/oylshe1314/framework/client/db/redis"
	"github.com/oylshe1314/framework/util"
)

type RedisClient struct {
	username string
	password string

	redis.Redis
	databaseClient
}

func (this *RedisClient) WithUsername(username string) {
	this.username = username
}

func (this *RedisClient) WithPassword(password string) {
	this.username = password
}

func (this *RedisClient) Init() (err error) {
	err = this.databaseClient.Init()
	if err != nil {
		return err
	}

	var db int
	if this.database != "" {
		err = util.StringToInteger2(this.database, &db)
		if err != nil {
			return err
		}
	}

	this.Redis = redis.OpenRedis(this.address, this.username, this.password, db)

	return
}

func (this *RedisClient) Close() (err error) {
	_ = this.databaseClient.Close()
	if this.Redis != nil {
		err = this.Redis.Close()
	}
	return
}

// Counter return the value before the increment
func (this *RedisClient) Counter(key string, inc uint64) (uint64, error) {
	if inc == 0 {
		return 0, nil
	}

	val, err := this.String(this.Context(), "incrby", key, util.IntegerToString(inc))
	if err != nil {
		return 0, err
	}

	counter, err := util.StringToInteger1[uint64](val)
	if err != nil {
		return 0, err
	}

	return counter + 1 - inc, nil
}

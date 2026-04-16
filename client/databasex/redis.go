package databasex

import "framework/option"

type RedisOption struct {
	Addresses []string `json:"addresses"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	DB        int      `json:"db"`
}

type RedisClient struct {
	option.Optional[RedisOption]
}

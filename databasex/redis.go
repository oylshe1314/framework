package databasex

import (
	"context"
	"time"

	"github.com/oylshe1314/framework/databasex/redis"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/option"
)

type RedisOption struct {
	Addresses       []string      `json:"addresses"`
	Protocol        int           `json:"protocol"`
	Username        string        `json:"username"`
	Password        string        `json:"password"`
	DB              int           `json:"DB"`
	ReadOnly        bool          `json:"readOnly"`
	DialTimeout     time.Duration `json:"dialTimeout"`
	ReadTimeout     time.Duration `json:"readTimeout"`
	WriteTimeout    time.Duration `json:"writeTimeout"`
	PoolSize        int           `json:"poolSize"`
	PoolTimeout     time.Duration `json:"poolTimeout"`
	MinIdleConns    int           `json:"minIdleConns"`
	MaxIdleConns    int           `json:"maxIdleConns"`
	MaxActiveConns  int           `json:"maxActiveConns"`
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime"`
	ConnMaxLifetime time.Duration `json:"connMaxLifetime"`
}

type RedisClient struct {
	option.Optional[RedisOption]

	redis.Redis
}

func (this *RedisClient) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("option is nil")
	}

	if len(this.GetOption().Addresses) == 0 {
		return errors.New("option 'addresses' is empty")
	}

	return this.open()
}

func (this *RedisClient) open() error {
	if len(this.GetOption().Addresses) > 1 {
		return this.openCluster()
	} else {
		return this.openRedis()
	}
}

func (this *RedisClient) openCluster() error {
	this.Redis = redis.OpenCluster(this.GetOption().Addresses, func(options *redis.ClusterOptions) {
		if this.GetOption().Protocol != 0 {
			options.Protocol = this.GetOption().Protocol
		}
		if this.GetOption().Username != "" {
			options.Username = this.GetOption().Username
		}
		if this.GetOption().Password != "" {
			options.Password = this.GetOption().Password
		}
		if this.GetOption().ReadOnly {
			options.ReadOnly = true
		}
		if this.GetOption().DialTimeout != 0 {
			options.DialTimeout = this.GetOption().DialTimeout
		}
		if this.GetOption().ReadTimeout != 0 {
			options.ReadTimeout = this.GetOption().ReadTimeout
		}
		if this.GetOption().WriteTimeout != 0 {
			options.WriteTimeout = this.GetOption().WriteTimeout
		}
		if this.GetOption().PoolSize != 0 {
			options.PoolSize = this.GetOption().PoolSize
		}
		if this.GetOption().PoolTimeout != 0 {
			options.PoolTimeout = this.GetOption().PoolTimeout
		}
		if this.GetOption().MinIdleConns != 0 {
			options.MinIdleConns = this.GetOption().MinIdleConns
		}
		if this.GetOption().MaxIdleConns != 0 {
			options.MaxIdleConns = this.GetOption().MaxIdleConns
		}
		if this.GetOption().MaxActiveConns != 0 {
			options.MaxActiveConns = this.GetOption().MaxActiveConns
		}
		if this.GetOption().ConnMaxIdleTime != 0 {
			options.ConnMaxIdleTime = this.GetOption().ConnMaxIdleTime
		}
		if this.GetOption().ConnMaxLifetime != 0 {
			options.ConnMaxLifetime = this.GetOption().ConnMaxLifetime
		}
	})
	return nil
}

func (this *RedisClient) openRedis() error {
	this.Redis = redis.OpenRedis(this.GetOption().Addresses[0], func(options *redis.Options) {
		if this.GetOption().Protocol != 0 {
			options.Protocol = this.GetOption().Protocol
		}
		if this.GetOption().Username != "" {
			options.Username = this.GetOption().Username
		}
		if this.GetOption().Password != "" {
			options.Password = this.GetOption().Password
		}
		if this.GetOption().DB != 0 {
			options.DB = this.GetOption().DB
		}
		if this.GetOption().DialTimeout != 0 {
			options.DialTimeout = this.GetOption().DialTimeout
		}
		if this.GetOption().ReadTimeout != 0 {
			options.ReadTimeout = this.GetOption().ReadTimeout
		}
		if this.GetOption().WriteTimeout != 0 {
			options.WriteTimeout = this.GetOption().WriteTimeout
		}
		if this.GetOption().PoolSize != 0 {
			options.PoolSize = this.GetOption().PoolSize
		}
		if this.GetOption().PoolTimeout != 0 {
			options.PoolTimeout = this.GetOption().PoolTimeout
		}
		if this.GetOption().MinIdleConns != 0 {
			options.MinIdleConns = this.GetOption().MinIdleConns
		}
		if this.GetOption().MaxIdleConns != 0 {
			options.MaxIdleConns = this.GetOption().MaxIdleConns
		}
		if this.GetOption().MaxActiveConns != 0 {
			options.MaxActiveConns = this.GetOption().MaxActiveConns
		}
		if this.GetOption().ConnMaxIdleTime != 0 {
			options.ConnMaxIdleTime = this.GetOption().ConnMaxIdleTime
		}
		if this.GetOption().ConnMaxLifetime != 0 {
			options.ConnMaxLifetime = this.GetOption().ConnMaxLifetime
		}
	})
	return nil
}

func (this *RedisClient) Close() error {
	if this.Redis != nil {
		return this.Redis.Close()
	}
	return nil
}

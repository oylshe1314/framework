package databasex

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/option"
)

type PostgresOption struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`

	ConnectTimeout  time.Duration `json:"connectTimeout"`
	MinConns        int32         `json:"minConns"`
	MaxConns        int32         `json:"maxConns"`
	MinIdleConns    int32         `json:"minIdleConns"`
	MaxConnLifetime time.Duration `json:"maxConnLifetime"`
	MaxConnIdleTime time.Duration `json:"maxConnIdleTime"`

	Options map[string]string `json:"options"`
}

type PostgresClient struct {
	option.Optional[PostgresOption]

	*pgxpool.Pool
}

func (this *PostgresClient) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("option is nil")
	}

	if len(this.GetOption().Address) == 0 {
		return errors.New("option 'address' is empty")
	}

	if len(this.GetOption().Username) == 0 {
		return errors.New("option 'username' is empty")
	}

	if len(this.GetOption().Password) == 0 {
		return errors.New("option 'password' is empty")
	}

	if len(this.GetOption().Database) == 0 {
		return errors.New("option 'database' is empty")
	}

	return this.connectPostgres(ctx)
}

func (this *PostgresClient) connectPostgres(ctx context.Context) error {
	var uri = fmt.Sprintf("postgres://%s:%s@%s/%s", this.GetOption().Username, this.GetOption().Password, this.GetOption().Address, this.GetOption().Username)

	var config, err = pgxpool.ParseConfig(uri)
	if err != nil {
		return err
	}

	if this.GetOption().ConnectTimeout != 0 {
		config.ConnConfig.ConnectTimeout = this.GetOption().ConnectTimeout
	}
	if this.GetOption().MinConns != 0 {
		config.MinConns = this.GetOption().MinConns
	}
	if this.GetOption().MaxConns != 0 {
		config.MaxConns = this.GetOption().MaxConns
	}
	if this.GetOption().MinIdleConns != 0 {
		config.MinIdleConns = this.GetOption().MinIdleConns
	}
	if this.GetOption().MaxConnLifetime != 0 {
		config.MaxConnLifetime = this.GetOption().MaxConnLifetime
	}
	if this.GetOption().MaxConnIdleTime != 0 {
		config.MaxConnIdleTime = this.GetOption().MaxConnIdleTime
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return err
	}

	this.Pool = pool
	return nil
}

func (this *PostgresClient) Close() error {
	if this.Pool != nil {
		this.Pool.Close()
	}
	return nil
}

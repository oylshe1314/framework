package databasex

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/option"
)

type MysqlOption struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`

	MaxOpenConns    int           `json:"maxOpenConns"`
	MaxIdleConns    int           `json:"maxIdleConns"`
	ConnMaxLifetime time.Duration `json:"connMaxLifetime"`
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime"`

	Options map[string]string `json:"options"`
}

type MysqlClient struct {
	option.Optional[MysqlOption]

	*sql.DB
}

func (this *MysqlClient) Init(ctx context.Context) error {
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

	return this.openMysql()
}

func (this *MysqlClient) openMysql() error {
	var uri = fmt.Sprintf("%s:%s@%s/%s", this.GetOption().Username, this.GetOption().Password, this.GetOption().Address, this.GetOption().Database)

	db, err := sql.Open("mysql", uri)
	if err != nil {
		return err
	}

	if this.GetOption().MaxOpenConns != 0 {
		db.SetMaxOpenConns(this.GetOption().MaxOpenConns)
	}
	if this.GetOption().MaxIdleConns != 0 {
		db.SetMaxIdleConns(this.GetOption().MaxIdleConns)
	}
	if this.GetOption().ConnMaxLifetime != 0 {
		db.SetConnMaxLifetime(this.GetOption().ConnMaxLifetime)
	}
	if this.GetOption().ConnMaxIdleTime != 0 {
		db.SetConnMaxIdleTime(this.GetOption().ConnMaxIdleTime)
	}

	this.DB = db
	return nil
}

func (this *MysqlClient) Close() error {
	if this.DB != nil {
		return this.DB.Close()
	}
	return nil
}

package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/util"
	"time"
)

type MysqlClient struct {
	*sql.DB
	databaseClient

	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime int
}

func (this *MysqlClient) WithMaxOpenConns(maxOpenConns int) {
	this.maxOpenConns = maxOpenConns
}

func (this *MysqlClient) WithMaxIdleConns(maxIdleConns int) {
	this.maxIdleConns = maxIdleConns
}

func (this *MysqlClient) WithConnMaxLifetime(connMaxLifetime int) {
	this.connMaxLifetime = connMaxLifetime
}

func (this *MysqlClient) Init() (err error) {
	err = this.databaseClient.Init()
	if err != nil {
		return err
	}

	this.DB, err = sql.Open("mysql", this.address+"/"+this.database+"?charset=utf8mb4")
	if err != nil {
		return err
	}

	if this.maxOpenConns > 0 {
		this.DB.SetMaxOpenConns(this.maxOpenConns)
	}

	if this.maxIdleConns > 0 {
		this.DB.SetMaxIdleConns(this.maxIdleConns)
	}

	if this.connMaxLifetime > 0 {
		this.DB.SetConnMaxLifetime(time.Duration(this.connMaxLifetime))
	}
	return nil
}

func (this *MysqlClient) Close() (err error) {
	_ = this.databaseClient.Close()
	if this.DB != nil {
		err = this.DB.Close()
	}
	return
}

// Counter
//
//	Table create SQL:
//	create table counter
//	(
//		`key`   varchar(255)    not null primary key,
//		`value` bigint unsigned not null
//	);
func (this *MysqlClient) Counter(key string, inc uint64) (uint64, error) {
	if inc == 0 {
		return 0, nil
	}

	var pair = &util.Pair[string, uint64]{Key: key}

	tx, err := this.Begin()
	if err != nil {
		return 0, err
	}

	row := tx.QueryRowContext(this.Context(), "select `value` from counter where `key`=?;", pair.Key)
	if err = row.Scan(&pair.Value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			pair.Value = 0
			_, err = tx.Exec("insert into counter (`key`, `value`) value (?, ?);", pair.Key, pair.Value+inc)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	} else {
		_, err = tx.Exec("update counter set `value`=? where `key`=?;", pair.Value+inc, pair.Key)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return pair.Value + 1, nil
}

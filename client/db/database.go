package db

import (
	"context"
	"framework/errors"
)

type databaseClient struct {
	address  string
	database string

	ctx    context.Context
	cancel context.CancelFunc
}

func (this *databaseClient) WithAddress(address string) {
	this.address = address
}

func (this *databaseClient) WithDatabase(database string) {
	this.database = database
}

func (this *databaseClient) Address() string {
	return this.address
}

func (this *databaseClient) Database() string {
	return this.database
}

func (this *databaseClient) Context() context.Context {
	return this.ctx
}

func (this *databaseClient) Init() (err error) {
	if len(this.address) == 0 {
		return errors.Error("'address' can not be empty")
	}

	if len(this.database) == 0 {
		return errors.Error("'database' can not be empty")
	}

	this.ctx, this.cancel = context.WithCancel(context.Background())

	return nil
}

func (this *databaseClient) Close() error {
	if this.cancel != nil {
		this.cancel()
	}
	return nil
}

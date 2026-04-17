package databasex

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/option"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoOption struct {
	Address  string `json:"address"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`

	ConnectTimeout  time.Duration `json:"connectTimeout"`
	ExecuteTimeout  time.Duration `json:"executeTimeout"`
	MinPoolSize     uint64        `json:"minPoolSize"`
	MaxPoolSize     uint64        `json:"maxPoolSize"`
	PoolTimeout     time.Duration `json:"poolTimeout"`
	MaxConnIdleTime time.Duration `json:"maxConnIdleTime"`

	Options map[string]string `json:"options"`
}

type MongoClient struct {
	option.Optional[MongoOption]

	client *mongo.Client

	*mongo.Database
}

func (this *MongoClient) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("option is nil")
	}

	if len(this.GetOption().Address) == 0 {
		return errors.New("option 'address' is empty")
	}

	if len(this.GetOption().Database) == 0 {
		return errors.New("option 'database' is empty")
	}

	return this.connectMongo()
}

func (this *MongoClient) connectMongo() error {

	var uri string
	if this.GetOption().Username != "" && this.GetOption().Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s/%s", this.GetOption().Username, this.GetOption().Password, this.GetOption().Address, this.GetOption().Database)
	} else {
		uri = fmt.Sprintf("mongodb://%s/%s", this.GetOption().Address, this.GetOption().Database)
	}

	if this.GetOption().Options != nil {
		var urlValues = url.Values{}
		for name, value := range this.GetOption().Options {
			urlValues.Set(name, value)
		}
		uri = uri + "?" + urlValues.Encode()
	}

	var clientOptions = options.Client().ApplyURI(uri).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	if this.GetOption().ConnectTimeout != 0 {
		clientOptions.SetConnectTimeout(this.GetOption().ConnectTimeout)
	}
	if this.GetOption().ExecuteTimeout != 0 {
		clientOptions.SetTimeout(this.GetOption().ExecuteTimeout)
	}
	if this.GetOption().MinPoolSize != 0 {
		clientOptions.SetMinPoolSize(this.GetOption().MinPoolSize)
	}
	if this.GetOption().MaxPoolSize != 0 {
		clientOptions.SetMaxPoolSize(this.GetOption().MaxPoolSize)
	}
	if this.GetOption().PoolTimeout != 0 {
		clientOptions.SetServerSelectionTimeout(this.GetOption().PoolTimeout)
	}
	if this.GetOption().MaxConnIdleTime != 0 {
		clientOptions.SetMaxConnIdleTime(this.GetOption().MaxConnIdleTime)
	}

	var err error
	this.client, err = mongo.Connect(clientOptions)
	if err != nil {
		return err
	}

	this.Database = this.client.Database(this.GetOption().Database)
	return nil
}

func (this *MongoClient) Close() error {
	if this.client != nil {
		return this.client.Disconnect(nil)
	}
	return nil
}

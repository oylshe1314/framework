package db

import (
	"framework/errors"
	"framework/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	databaseClient

	client   *mongo.Client
	database *mongo.Database
}

func (this *MongoClient) Init() (err error) {
	err = this.databaseClient.Init()
	if err != nil {
		return err
	}

	this.client, err = mongo.Connect(this.ctx, options.Client().ApplyURI(this.address))
	if err != nil {
		return err
	}

	this.database = this.client.Database(this.Database())
	return nil
}

func (this *MongoClient) Close() (err error) {
	_ = this.databaseClient.Close()
	if this.client != nil {
		err = this.client.Disconnect(nil)
	}
	return
}

func (this *MongoClient) Collection(name string) *mongo.Collection {
	return this.database.Collection(name)
}

func (this *MongoClient) Counter(key string, inc uint64) (uint64, error) {
	if inc == 0 {
		return 0, nil
	}

	var upsert = true
	var pair = &util.Pair[string, uint64]{Key: key}
	for {
		var err = this.Collection("counter").FindOneAndUpdate(
			this.Context(),
			bson.M{"_id": pair.Key},
			bson.M{"$inc": bson.M{"value": inc}},
			&options.FindOneAndUpdateOptions{Upsert: &upsert},
		).Decode(&pair)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				return 0, err
			}
		}
		return pair.Value + 1, nil
	}
}

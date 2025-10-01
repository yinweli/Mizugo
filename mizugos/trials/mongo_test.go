package trials

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestMongo(t *testing.T) {
	suite.Run(t, new(SuiteMongo))
}

type SuiteMongo struct {
	suite.Suite
}

func (this *SuiteMongo) TestMongoExist() {
	client := newMongo()
	_, _ = client.Database("mongo").Collection("exist").InsertOne(context.Background(), bson.M{"value": 1})
	this.True(MongoExist(client.Database("mongo"), "exist", "value", 1))
	this.False(MongoExist(client.Database("mongo"), "exist", "value", testdata.Unknown))
	_ = client.Database("mongo").Drop(context.Background())
}

func (this *SuiteMongo) TestMongoEqual() {
	client := newMongo()
	_, _ = client.Database("mongo").Collection("compare").InsertOne(context.Background(), bson.M{"value": 1})
	this.True(MongoEqual[testMongo](client.Database("mongo"), "compare", "value", 1, &testMongo{Value: 1}))
	this.False(MongoEqual[testMongo](client.Database("mongo"), "compare", "value", testdata.Unknown, &testMongo{Value: 1}))
	this.False(MongoEqual[testMongo](client.Database("mongo"), "compare", "value", 1, &testMongo{Value: 2}))
	_ = client.Database("mongo").Drop(context.Background())
}

func newMongo() *mongo.Client {
	option := options.Client().ApplyURI(testdata.MongoURI)
	client, _ := mongo.Connect(context.Background(), option)
	return client
}

type testMongo struct {
	Value int `bson:"value"`
}

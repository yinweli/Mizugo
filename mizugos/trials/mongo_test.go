package trials

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yinweli/Mizugo/testdata"
)

func TestMongo(t *testing.T) {
	suite.Run(t, new(SuiteMongo))
}

type SuiteMongo struct {
	suite.Suite
}

func (this *SuiteMongo) TestMongoExist() {
	dbname := "exist"
	table := "table"
	field := "field"
	key := "key"
	client := newMongo()
	_, _ = client.Database(dbname).Collection(table).UpdateOne(context.Background(),
		bson.M{field: key},
		bson.M{"$set": bson.M{field: key}},
		options.Update().SetUpsert(true))
	assert.True(this.T(), MongoExist(client.Database(dbname), table, field, key))
	assert.False(this.T(), MongoExist(client.Database(dbname), table, field, testdata.Unknown))
	_ = client.Database(dbname).Drop(context.Background())
}

func (this *SuiteMongo) TestMongoCompare() {
	dbname := "compare"
	table := "table"
	field := "value"
	client := newMongo()
	_, _ = client.Database(dbname).Collection(table).InsertOne(context.Background(), bson.M{field: "1"})
	assert.True(this.T(), MongoCompare[testMongo](client.Database(dbname), table, field, "1", &testMongo{Value: "1"}))
	assert.False(this.T(), MongoCompare[testMongo](client.Database(dbname), table, field, testdata.Unknown, &testMongo{Value: "1"}))
	_ = client.Database(dbname).Drop(context.Background())
}

func newMongo() *mongo.Client {
	option := options.Client().ApplyURI(testdata.MongoURI)
	client, _ := mongo.Connect(context.Background(), option)
	return client
}

type testMongo struct {
	Value string `bson:"value"`
}

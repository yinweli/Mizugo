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
		bson.D{{Key: field, Value: key}},
		bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: key}}}},
		options.Update().SetUpsert(true))
	assert.True(this.T(), MongoExist(client.Database(dbname), table, field, key))
	assert.False(this.T(), MongoExist(client.Database(dbname), table, field, testdata.Unknown))
	_ = client.Database(dbname).Drop(context.Background())
}

func (this *SuiteMongo) TestMongoCompare() {
	dbname := "compare"
	table := "table"
	field := "field"
	key := "999"
	client := newMongo()
	_, _ = client.Database(dbname).Collection(table).UpdateOne(context.Background(),
		bson.D{{Key: field, Value: key}},
		bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: key}}}},
		options.Update().SetUpsert(true))
	assert.True(this.T(), MongoCompare[testMongo](client.Database(dbname), table, field, key, &testMongo{Field: key}))
	assert.False(this.T(), MongoExist(client.Database(dbname), table, field, testdata.Unknown))
	_ = client.Database(dbname).Drop(context.Background())
}

func newMongo() *mongo.Client {
	option := options.Client().ApplyURI(testdata.MongoURI)
	client, _ := mongo.Connect(context.Background(), option)
	return client
}

type testMongo struct {
	Field string `bson:"field"`
}

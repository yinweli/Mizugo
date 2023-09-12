package testdata

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoExist 在mongo中資料是否存在
func MongoExist(database *mongo.Database, tableName, fieldName, key string) bool {
	table := database.Collection(tableName)

	if table == nil {
		return false
	} // if

	if table.FindOne(context.Background(), bson.D{{Key: fieldName, Value: key}}).Err() != nil {
		return false
	} // if

	return true
}

// MongoCompare 在mongo中比對資料是否相同
func MongoCompare[T any](database *mongo.Database, tableName, fieldName, key string, expected *T, cmpOpt ...cmp.Option) bool {
	table := database.Collection(tableName)

	if table == nil {
		return false
	} // if

	actual := new(T)

	if table.FindOne(context.Background(), bson.D{{Key: fieldName, Value: key}}).Decode(actual) != nil {
		return false
	} // if

	return cmp.Equal(expected, actual, cmpOpt...)
}

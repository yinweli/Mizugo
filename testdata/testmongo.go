package testdata

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClear 清除mongo, 直接刪除資料庫
func MongoClear(table *mongo.Database) {
	_ = table.Drop(context.Background())
}

// MongoCompare 在mongo中比對資料是否相同
func MongoCompare[T any](database *mongo.Database, tableName, fieldName, key string, expected *T, cmpOpt ...cmp.Option) bool {
	table := database.Collection(tableName)

	if table == nil {
		return false
	} // if

	actual := new(T)

	if table.FindOne(context.Background(), bson.D{{Key: fieldName, Value: key}}, options.FindOne()).Decode(actual) != nil {
		return false
	} // if

	return cmp.Equal(expected, actual, cmpOpt...)
}

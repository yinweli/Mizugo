package testdata

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClear 清除mongo, 直接刪除資料庫
func MongoClear(ctx context.Context, table *mongo.Database) {
	_ = table.Drop(ctx)
}

// MongoCompare 在mongo中比對資料是否相同
func MongoCompare[T any](ctx context.Context, database *mongo.Database, tableName, fieldName, key string, expected *T) bool {
	table := database.Collection(tableName)

	if table == nil {
		return false
	} // if

	actual := new(T)

	if table.FindOne(ctx, bson.D{{Key: fieldName, Value: key}}, options.FindOne()).Decode(actual) != nil {
		return false
	} // if

	return reflect.DeepEqual(expected, actual)
}

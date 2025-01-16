package trials

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoExist 在mongo中資料是否存在
func MongoExist(database *mongo.Database, table, field, key string) bool {
	collection := database.Collection(table)

	if collection == nil {
		return false
	} // if

	if collection.FindOne(context.Background(), bson.M{field: key}, options.FindOne()).Err() != nil {
		return false
	} // if

	return true
}

// MongoCompare 在mongo中比對資料是否相同
func MongoCompare[T any](database *mongo.Database, table, field, key string, expected *T, option ...cmp.Option) bool {
	collection := database.Collection(table)

	if collection == nil {
		return false
	} // if

	actual := new(T)

	if collection.FindOne(context.Background(), bson.M{field: key}, options.FindOne()).Decode(actual) != nil {
		return false
	} // if

	return cmp.Equal(expected, actual, option...)
}

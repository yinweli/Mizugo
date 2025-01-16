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

	if collection.FindOne(context.Background(), bson.D{{Key: field, Value: key}}, options.FindOne()).Err() != nil {
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

	if collection.FindOne(context.Background(), bson.D{{Key: field, Value: key}}, options.FindOne()).Decode(actual) != nil {
		return false
	} // if

	return cmp.Equal(expected, actual, option...)
}

// MongoCompareList 在mongo中比對列表是否相同, sort為排序欄位, asc = 1 為升序, -1為降序
func MongoCompareList[T any](database *mongo.Database, table, sort string, asc int, expected []*T, option ...cmp.Option) bool {
	collection := database.Collection(table)

	if collection == nil {
		return false
	} // if

	result, err := collection.Find(context.Background(), bson.D{}, options.Find().SetSort(bson.D{{Key: sort, Value: asc}}))

	if err != nil {
		return false
	} // if

	defer func() {
		_ = result.Close(context.Background())
	}()
	actual := []*T{}

	for result.Next(context.Background()) {
		a := new(T)

		if err = result.Decode(a); err != nil {
			return false
		} // if

		actual = append(actual, a)
	} // for

	return cmp.Equal(expected, actual, option...)
}

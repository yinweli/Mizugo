package trials

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoExist 檢查指定索引是否存在
func MongoExist(database *mongo.Database, table, field string, key any) bool {
	if err := database.Collection(table).FindOne(context.Background(), bson.M{field: key}, options.FindOne()).Err(); err != nil {
		fmt.Printf("mongo not exist: table %v: field %v: key %v: %v\n", table, field, key, err)
		return false
	} // if

	return true
}

// MongoEqual 比對資料是否符合預期
func MongoEqual[T any](database *mongo.Database, table, field string, key any, expected *T, option ...cmp.Option) bool {
	actual := new(T)

	if err := database.Collection(table).FindOne(context.Background(), bson.M{field: key}, options.FindOne()).Decode(actual); err != nil {
		fmt.Printf("mongo not equal: table %v: field %v: key %v: %v\n", table, field, key, err)
		return false
	} // if

	if cmp.Equal(expected, actual, option...) == false {
		fmt.Printf("mongo not equal: table %v: field %v: key %v\n", table, field, key)
		fmt.Println("  expected:")
		fmt.Printf("    %+v\n", expected)
		fmt.Println("  actual:")
		fmt.Printf("    %+v\n", actual)
		return false
	} // if

	return true
}

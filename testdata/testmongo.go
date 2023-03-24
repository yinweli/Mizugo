package testdata

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClear 清除mongo, 直接刪除資料庫
func MongoClear(ctx context.Context, table *mongo.Database) {
	_ = table.Drop(ctx)
}

// MongoFindOne 在mongo中搜尋單一結果
func MongoFindOne(ctx context.Context, database *mongo.Database, tableName, fieldName, key string, result any) bool {
	if table := database.Collection(tableName); table != nil {
		if table.FindOne(ctx, bson.D{{Key: fieldName, Value: key}}, options.FindOne()).Decode(result) == nil {
			return true
		} // if
	} // if

	return false
}

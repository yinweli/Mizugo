package querys

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos/redmos"
)

// MetaAuth 元資料: 認證資料
type MetaAuth struct {
}

func (this *MetaAuth) MajorKey(key any) string {
	return fmt.Sprintf("auth:%v", key)
}

func (this *MetaAuth) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *MetaAuth) MinorTable() string {
	return "auth"
}

func (this *MetaAuth) MinorField() string {
	return "account"
}

// Auth 認證資料
type Auth struct {
	Account string    `bson:"account"` // [主索引] 帳號
	Token   string    `bson:"token"`   // token
	Time    time.Time `bson:"time"`    // 更新時間
}

// NewAuthGet 建立取得認證資料行為
func NewAuthGet(key string, data *Auth) *redmos.Get[Auth] {
	return &redmos.Get[Auth]{
		Meta: &MetaAuth{},
		Key:  key,
		Data: data,
	}
}

// NewAuthSet 建立設定認證資料行為
func NewAuthSet(key string, data *Auth) *redmos.Set[Auth] {
	return &redmos.Set[Auth]{
		Meta: &MetaAuth{},
		Key:  key,
		Data: data,
	}
}

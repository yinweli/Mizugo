package querys

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/v2/mizugos/redmos"
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

// NewAuth 建立玩家資料
func NewAuth(account string) *Auth {
	return &Auth{
		Save:    redmos.NewSave(),
		Account: account,
	}
}

// Auth 認證資料
type Auth struct {
	*redmos.Save           // 儲存判斷資料
	Account      string    `bson:"account"` // [主索引] 帳號
	Token        string    `bson:"token"`   // token
	Time         time.Time `bson:"time"`    // 更新時間
}

// NewGetter 建立取得資料行為
func (this *Auth) NewGetter() redmos.Behavior {
	return &redmos.Get[Auth]{
		Meta:        &MetaAuth{},
		MajorEnable: true,
		Key:         this.Account,
		Data:        this,
	}
}

// NewSetter 建立設定資料行為
func (this *Auth) NewSetter() redmos.Behavior {
	return &redmos.Set[Auth]{
		Meta:        &MetaAuth{},
		MajorEnable: true,
		MinorEnable: true,
		Key:         this.Account,
		Data:        this,
	}
}

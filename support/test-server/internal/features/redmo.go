package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/redmos"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
)

// NewRedmo 建立資料庫資料
func NewRedmo() *Redmo {
	return &Redmo{
		name: "redmo",
	}
}

// Redmo 資料庫資料
type Redmo struct {
	name   string      // 資料庫名稱
	config RedmoConfig // 配置資料
}

// RedmoConfig 配置資料
type RedmoConfig struct {
	MajorURI redmos.RedisURI `yaml:"majorURI"` // 主要資料庫連接字串
	MinorURI redmos.MongoURI `yaml:"minorURI"` // 次要資料庫連接字串
}

// Initialize 初始化處理
func (this *Redmo) Initialize() error {
	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Redmomgr().AddMajor(defines.MajorName, this.config.MajorURI); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Redmomgr().AddMinor(defines.MinorName, this.config.MinorURI, defines.MongoDB); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Redmomgr().AddMixed(defines.MixedName, defines.MajorName, defines.MinorName); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Info(this.name).Caller(0).Message("initialize").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Redmo) Finalize() {
	mizugos.Redmomgr().Stop()
}

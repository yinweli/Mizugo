package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/redmos"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
)

const nameRedmo = "redmo" // 特性名稱

// NewRedmo 建立資料庫資料
func NewRedmo() *Redmo {
	return &Redmo{}
}

// Redmo 資料庫資料
type Redmo struct {
	config RedmoConfig // 配置資料
}

// RedmoConfig 配置資料
type RedmoConfig struct {
	MajorURI redmos.RedisURI `yaml:"majorURI"` // 主要資料庫連接字串
	MinorURI redmos.MongoURI `yaml:"minorURI"` // 次要資料庫連接字串
}

// Initialize 初始化處理
func (this *Redmo) Initialize() error {
	if err := mizugos.Configmgr().Unmarshal(nameRedmo, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", nameRedmo, err)
	} // if

	if _, err := mizugos.Redmomgr().AddMajor(defines.RedmoMajor, this.config.MajorURI, false); err != nil {
		return fmt.Errorf("%v initialize: %w", nameRedmo, err)
	} // if

	if _, err := mizugos.Redmomgr().AddMinor(defines.RedmoMinor, this.config.MinorURI, defines.MongoDB); err != nil {
		return fmt.Errorf("%v initialize: %w", nameRedmo, err)
	} // if

	if _, err := mizugos.Redmomgr().AddMixed(defines.RedmoMixed, defines.RedmoMajor, defines.RedmoMinor); err != nil {
		return fmt.Errorf("%v initialize: %w", nameRedmo, err)
	} // if

	System.Info(nameRedmo).Caller(0).Message("initialize").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Redmo) Finalize() {
	mizugos.Redmomgr().Finalize()
}

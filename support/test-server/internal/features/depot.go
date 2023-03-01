package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/depots"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
)

// NewDepot 建立資料庫資料
func NewDepot() *Depot {
	return &Depot{
		name: "depot",
	}
}

// Depot 資料庫資料
type Depot struct {
	name   string      // 資料庫名稱
	config DepotConfig // 配置資料
}

// DepotConfig 配置資料
type DepotConfig struct {
	MajorURI depots.RedisURI `yaml:"majorURI"` // 主要資料庫連接字串
	MinorURI depots.MongoURI `yaml:"minorURI"` // 次要資料庫連接字串
}

// Initialize 初始化處理
func (this *Depot) Initialize() error {
	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Depotmgr().AddMajor(defines.MajorName, this.config.MajorURI); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Depotmgr().AddMinor(defines.MinorName, this.config.MinorURI, defines.MongoDB); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Depotmgr().AddMixed(defines.MixedName, defines.MajorName, defines.MinorName); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Info(this.name).Caller(0).Message("initialize").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Depot) Finalize() {
	mizugos.Depotmgr().Stop()
}

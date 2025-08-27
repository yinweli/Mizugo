package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/v2/mizugos"
	"github.com/yinweli/Mizugo/v2/mizugos/redmos"
)

// RedmoInitialize 初始化資料庫
func RedmoInitialize() (err error) {
	config := &RedmoConfig{}

	if err = mizugos.Config.Unmarshal("redmo", config); err != nil {
		return fmt.Errorf("redmo initialize: %w", err)
	} // if

	if DBMajor, err = mizugos.Redmo.AddMajor("dbmajor", config.MajorURI); err != nil {
		return fmt.Errorf("redmo initialize: %w", err)
	} // if

	if DBMinor, err = mizugos.Redmo.AddMinor("dbminor", config.MinorURI, config.MinorDBName); err != nil {
		return fmt.Errorf("redmo initialize: %w", err)
	} // if

	if DBMixed, err = mizugos.Redmo.AddMixed("dbmixed", "dbmajor", "dbminor"); err != nil {
		return fmt.Errorf("redmo initialize: %w", err)
	} // if

	LogSystem.Get().Info("redmo").Message("initialize").EndFlush()
	return nil
}

// RedmoConfig 配置資料
type RedmoConfig struct {
	MajorURI    redmos.RedisURI `yaml:"majorURI"`    // 主要資料庫連接字串
	MinorURI    redmos.MongoURI `yaml:"minorURI"`    // 次要資料庫連接字串
	MinorDBName string          `yaml:"minorDBName"` // 次要資料庫資料庫名稱
}

var DBMajor *redmos.Major // 主要資料庫
var DBMinor *redmos.Minor // 次要資料庫
var DBMixed *redmos.Mixed // 混合資料庫

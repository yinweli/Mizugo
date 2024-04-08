package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugo/redmos"
)

// InitializeRedmo 初始化資料庫管理器
func InitializeRedmo() (err error) {
	name := "redmo"
	config := &RedmoConfig{}
	Redmo = redmos.NewRedmomgr()

	if err = Config.Unmarshal(name, config); err != nil {
		return fmt.Errorf("%v initialize: %w", name, err)
	} // if

	if DBMajor, err = Redmo.AddMajor("dbmajor", config.MajorURI); err != nil {
		return fmt.Errorf("%v initialize: %w", name, err)
	} // if

	if DBMinor, err = Redmo.AddMinor("dbminor", config.MinorURI, config.MinorDBName); err != nil {
		return fmt.Errorf("%v initialize: %w", name, err)
	} // if

	if DBMixed, err = Redmo.AddMixed("dbmixed", "dbmajor", "dbminor"); err != nil {
		return fmt.Errorf("%v initialize: %w", name, err)
	} // if

	LogSystem.Get().Info(name).Message("initialize").EndFlush()
	return nil
}

// FinalizeRedmo 結束資料庫管理器
func FinalizeRedmo() {
	if Redmo != nil {
		Redmo.Finalize()
	} // if
}

// RedmoConfig 配置資料
type RedmoConfig struct {
	MajorURI    redmos.RedisURI `yaml:"majorURI"`    // 主要資料庫連接字串
	MinorURI    redmos.MongoURI `yaml:"minorURI"`    // 次要資料庫連接字串
	MinorDBName string          `yaml:"minorDBName"` // 次要資料庫資料庫名稱
}

var Redmo *redmos.Redmomgr // 資料庫管理器
var DBMajor *redmos.Major  // 主要資料庫
var DBMinor *redmos.Minor  // 次要資料庫
var DBMixed *redmos.Mixed  // 混合資料庫

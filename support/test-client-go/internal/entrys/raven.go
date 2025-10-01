package entrys //nolint:dupl

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/v2/mizugos"
	"github.com/yinweli/Mizugo/v2/mizugos/cryptos"
	"github.com/yinweli/Mizugo/v2/mizugos/entitys"
	"github.com/yinweli/Mizugo/v2/mizugos/nets"
	"github.com/yinweli/Mizugo/v2/mizugos/procs"
	"github.com/yinweli/Mizugo/v2/support/test-client-go/internal/features"
	"github.com/yinweli/Mizugo/v2/support/test-client-go/internal/miscs"
	"github.com/yinweli/Mizugo/v2/support/test-client-go/internal/modules"
)

// RavenInitialize 初始化Raven入口
func RavenInitialize() (err error) {
	config := &RavenConfig{}

	if err = mizugos.Config.Unmarshal("raven", config); err != nil {
		return fmt.Errorf("raven initialize: %w", err)
	} // if

	raven.key = config.Key
	raven.delay = config.Delay
	raven.disconnect = config.Disconnect

	if config.Enable {
		miscs.GenerateConnection(config.Interval, config.Count, config.Batch, func() {
			mizugos.Network.AddConnectTCP(config.IP, config.Port, config.Timeout, raven.bind, raven.unbind, raven.wrong)
		})
	} // if

	features.LogSystem.Get().Info("raven").Message("initialize").EndFlush()
	return nil
}

// Raven Raven入口
type Raven struct {
	key        string        // 密鑰
	delay      time.Duration // 延遲時間
	disconnect bool          // 斷線旗標
}

// RavenConfig 配置資料
type RavenConfig struct {
	Enable     bool          `yaml:"enable"`     // 啟用旗標
	IP         string        `yaml:"ip"`         // 位址
	Port       string        `yaml:"port"`       // 埠號
	Key        string        `yaml:"key"`        // 密鑰
	Timeout    time.Duration `yaml:"timeout"`    // 逾時時間
	Interval   time.Duration `yaml:"interval"`   // 間隔時間
	Count      int           `yaml:"count"`      // 總連線數
	Batch      int           `yaml:"batch"`      // 批次連線數
	Delay      time.Duration `yaml:"delay"`      // 延遲時間
	Disconnect bool          `yaml:"disconnect"` // 斷線旗標
}

// bind 綁定處理
func (this *Raven) bind(session nets.Sessioner) bool {
	err := error(nil)
	entity := mizugos.Entity.Add()
	process := procs.NewRavenClient()

	if entity == nil {
		err = fmt.Errorf("bind: entity nil")
		goto Error
	} // if

	if err = entity.SetProcess(process); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err = entity.SetSession(session); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err = entity.AddModule(modules.NewRaven(this.delay, this.disconnect)); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err = entity.Initialize(this.wrong); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	session.SetCodec(process, cryptos.NewDesCBC(cryptos.PaddingPKCS7, this.key, this.key), cryptos.NewBase64())
	session.SetPublish(entity.PublishOnce)
	session.SetWrong(this.wrong)
	session.SetOwner(entity)
	features.LogSystem.Get().Info("raven").Message("bind").Caller(0).EndFlush()
	return true

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if

	session.Stop()
	features.LogSystem.Get().Error("raven").Caller(0).Error(err).EndFlush()
	return false
}

// unbind 解綁處理
func (this *Raven) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if
}

// wrong 錯誤處理
func (this *Raven) wrong(err error) {
	features.LogSystem.Get().Error("raven").Caller(1).Error(err).EndFlush()
}

var raven = &Raven{} // Raven入口

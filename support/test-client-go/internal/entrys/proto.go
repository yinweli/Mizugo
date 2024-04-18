package entrys //nolint:dupl

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/cryptos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/miscs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/modules"
)

// ProtoInitialize 初始化Proto入口
func ProtoInitialize() (err error) {
	config := &ProtoConfig{}

	if err = mizugos.Config.Unmarshal("proto", config); err != nil {
		return fmt.Errorf("proto initialize: %w", err)
	} // if

	proto.key = config.Key
	proto.delay = config.Delay
	proto.disconnect = config.Disconnect

	if config.Enable {
		miscs.GenerateConnection(config.Interval, config.Count, config.Batch, func() {
			mizugos.Network.AddConnectTCP(config.IP, config.Port, config.Timeout, proto.bind, proto.unbind, proto.wrong)
		})
	} // if

	features.LogSystem.Get().Info("proto").Message("initialize").EndFlush()
	return nil
}

// Proto Proto入口
type Proto struct {
	key        string        // 密鑰
	delay      time.Duration // 延遲時間
	disconnect bool          // 斷線旗標
}

// ProtoConfig 配置資料
type ProtoConfig struct {
	Enable     bool          `yaml:"enable"`     // 啟用旗標
	IP         string        `yaml:"ip"`         // 位址
	Port       string        `yaml:"port"`       // 埠號
	Key        string        `yaml:"key"`        // 密鑰
	Timeout    time.Duration `yaml:"timeout"`    // 超期時間
	Interval   time.Duration `yaml:"interval"`   // 間隔時間
	Count      int           `yaml:"count"`      // 總連線數
	Batch      int           `yaml:"batch"`      // 批次連線數
	Delay      time.Duration `yaml:"delay"`      // 延遲時間
	Disconnect bool          `yaml:"disconnect"` // 斷線旗標
}

// bind 綁定處理
func (this *Proto) bind(session nets.Sessioner) bool {
	err := error(nil)
	entity := mizugos.Entity.Add()
	process := procs.NewProto()
	desCBC := cryptos.NewDesCBC(cryptos.PaddingPKCS7, this.key, this.key)
	base64 := cryptos.NewBase64()

	session.SetPublish(entity.PublishOnce)
	session.SetWrong(this.wrong)
	session.SetCodec(process, desCBC, base64)

	if entity == nil {
		err = fmt.Errorf("bind: entity nil")
		goto Error
	} // if

	if err = entity.SetModulemap(entitys.NewModulemap()); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err = entity.SetEventmap(entitys.NewEventmap(defines.EventCapacity)); err != nil {
		err = fmt.Errorf("bind: %w", err)
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

	if err = entity.AddModule(modules.NewProto(this.delay, this.disconnect)); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err = entity.Initialize(this.wrong); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	session.SetOwner(entity)
	features.MeterConnect.Add(1)
	features.LogSystem.Get().Info("proto").Message("bind").Caller(0).EndFlush()
	return true

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if

	session.Stop()
	features.LogSystem.Get().Error("proto").Caller(0).Error(err).EndFlush()
	return false
}

// unbind 解綁處理
func (this *Proto) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
		features.MeterConnect.Add(-1)
	} // if
}

// wrong 錯誤處理
func (this *Proto) wrong(err error) {
	features.LogSystem.Get().Error("proto").Caller(1).Error(err).EndFlush()
}

var proto = &Proto{} // Proto入口

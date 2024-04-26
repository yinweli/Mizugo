package modules

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
	"github.com/yinweli/Mizugo/support/test-client-go/msgs"
)

// NewAuth 建立Auth模組
func NewAuth(delay time.Duration, account string, update int) *Auth {
	return &Auth{
		Module:  entitys.NewModule(defines.ModuleIDAuth),
		name:    "auth",
		delay:   delay,
		account: account,
		update:  update,
	}
}

// Auth Auth模組
type Auth struct {
	*entitys.Module               // 模組資料
	name            string        // 系統名稱
	delay           time.Duration // 延遲時間
	account         string        // 帳號
	update          int           // 更新次數
	token           string        // token
}

// Awake 喚醒處理
func (this *Auth) Awake() error {
	this.Entity().Subscribe(defines.EventBegin, this.eventBegin)
	this.Entity().AddMessage(int32(msgs.MsgID_LoginA), this.procMLoginA)
	this.Entity().AddMessage(int32(msgs.MsgID_UpdateA), this.procMUpdateA)
	return nil
}

// Start 啟動處理
func (this *Auth) Start() error {
	this.Entity().PublishDelay(defines.EventBegin, nil, this.delay)
	return nil
}

// event 開始事件
func (this *Auth) eventBegin(_ any) {
	this.sendMLoginQ()
}

// procMLoginA 處理回應登入
func (this *Auth) procMLoginA(message any) {
	_, msg, err := procs.JsonUnmarshal[msgs.MLoginA](message)

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("auth procMLoginA: %w", err)).EndFlush()
		return
	} // if

	if msgs.ErrID(msg.ErrID) != msgs.ErrID_Success {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("auth procMLoginA: %v", msg.ErrID)).EndFlush()
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.MeterAuth.Add(duration)
	features.LogSystem.Get().Info(this.name).KV("token", msg.Token).KV("duration", duration).Caller(0).EndFlush()

	this.token = msg.Token
	this.sendMUpdateQ()
}

// sendMLoginQ 傳送要求登入
func (this *Auth) sendMLoginQ() {
	msg, err := procs.JsonMarshal(int32(msgs.MsgID_LoginQ), &msgs.MLoginQ{
		Account: this.account,
		Time:    time.Now().UnixNano(),
	})

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(err).EndFlush()
		return
	} // if

	this.Entity().Send(msg)
}

// procMUpdateA 處理回應更新
func (this *Auth) procMUpdateA(message any) {
	_, msg, err := procs.JsonUnmarshal[msgs.MUpdateA](message)

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("auth procMUpdateA: %w", err)).EndFlush()
		return
	} // if

	if msgs.ErrID(msg.ErrID) != msgs.ErrID_Success {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("auth procMUpdateA: %v", msg.ErrID)).EndFlush()
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.MeterAuth.Add(duration)
	features.LogSystem.Get().Info(this.name).KV("token", msg.Token).KV("duration", duration).Caller(0).EndFlush()

	this.token = msg.Token
	this.sendMUpdateQ()
}

// sendMUpdateQ 傳送要求更新
func (this *Auth) sendMUpdateQ() {
	if this.update <= 0 {
		return
	} // if

	this.update--
	msg, err := procs.JsonMarshal(int32(msgs.MsgID_UpdateQ), &msgs.MUpdateQ{
		Account: this.account,
		Token:   this.token,
		Time:    time.Now().UnixNano(),
	})

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(err).EndFlush()
		return
	} // if

	this.Entity().Send(msg)
}

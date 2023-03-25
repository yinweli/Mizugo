package modules

import (
	"time"

	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/errs"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
	"github.com/yinweli/Mizugo/support/test-client-go/msgs"
)

const nameAuth = "module-auth" // 模組名稱

// NewAuth 建立Auth模組
func NewAuth(delay time.Duration, account string, update int) *Auth {
	return &Auth{
		Module:  entitys.NewModule(defines.ModuleIDAuth),
		delay:   delay,
		account: account,
		update:  update,
	}
}

// Auth Auth模組
type Auth struct {
	*entitys.Module               // 模組資料
	delay           time.Duration // 延遲時間
	account         string        // 帳號
	update          int           // 更新次數
	token           string        // token
}

// Awake 喚醒處理
func (this *Auth) Awake() error {
	this.Entity().Subscribe(defines.EventBegin, this.eventBegin)
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_LoginA), this.procMLoginA)
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_UpdateA), this.procMUpdateA)
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
		features.System.Warn(nameAuth).Caller(0).EndError(errs.Errore(msgs.ErrID_JsonUnmarshal, err))
		return
	} // if

	if msgs.ErrID(msg.ErrID) != msgs.ErrID_Success {
		features.System.Warn(nameAuth).Caller(0).EndError(errs.Errort(msg.ErrID))
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.Auth.Add(duration)
	features.System.Info(nameAuth).Caller(0).KV("duration", duration).KV("token", msg.Token).End()

	this.token = msg.Token
	this.sendMUpdateQ()
}

// sendMLoginQ 傳送要求登入
func (this *Auth) sendMLoginQ() {
	msg, err := procs.JsonMarshal(procs.MessageID(msgs.MsgID_LoginQ), &msgs.MLoginQ{
		Account: this.account,
		Time:    time.Now().UnixNano(),
	})

	if err != nil {
		features.System.Warn(nameAuth).Caller(0).EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}

// procMUpdateA 處理回應更新
func (this *Auth) procMUpdateA(message any) {
	_, msg, err := procs.JsonUnmarshal[msgs.MUpdateA](message)

	if err != nil {
		features.System.Warn(nameAuth).Caller(0).EndError(errs.Errore(msgs.ErrID_JsonUnmarshal, err))
		return
	} // if

	if msgs.ErrID(msg.ErrID) != msgs.ErrID_Success {
		features.System.Warn(nameAuth).Caller(0).EndError(errs.Errort(msg.ErrID))
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.Auth.Add(duration)
	features.System.Info(nameAuth).Caller(0).KV("duration", duration).KV("token", msg.Token).End()

	this.token = msg.Token
	this.sendMUpdateQ()
}

// sendMUpdateQ 傳送要求更新
func (this *Auth) sendMUpdateQ() {
	if this.update <= 0 {
		return
	} // if

	this.update--
	msg, err := procs.JsonMarshal(procs.MessageID(msgs.MsgID_UpdateQ), &msgs.MUpdateQ{
		Account: this.account,
		Token:   this.token,
		Time:    time.Now().UnixNano(),
	})

	if err != nil {
		features.System.Warn(nameAuth).Caller(0).EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}

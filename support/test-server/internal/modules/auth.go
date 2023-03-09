package modules

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/errs"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/redmos"
	"github.com/yinweli/Mizugo/support/test-server/internal/access"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/miscs"
	"github.com/yinweli/Mizugo/support/test-server/internal/msgs"
)

// NewAuth 建立Auth模組
func NewAuth() *Auth {
	return &Auth{
		Module: entitys.NewModule(defines.ModuleIDAuth),
		name:   "module auth",
	}
}

// Auth Auth模組
type Auth struct {
	*entitys.Module
	name     string        // 模組名稱
	database *redmos.Mixed // 資料庫物件
}

// Awake 喚醒處理
func (this *Auth) Awake() error {
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_LoginQ), this.procMLoginQ)
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_UpdateQ), this.procMUpdateQ)
	return nil
}

// Start 啟動處理
func (this *Auth) Start() error {
	if this.database = mizugos.Redmomgr().GetMixed(defines.RedmoMixed); this.database == nil {
		return fmt.Errorf("auth start: database nil")
	} // if

	return nil
}

// procMLoginQ 處理要求登入
func (this *Auth) procMLoginQ(message any) {
	rec := features.Login.Rec()
	defer rec()

	_, msg, err := procs.JsonUnmarshal[msgs.MLoginQ](message)

	if err != nil {
		this.sendMLoginA(nil, msgs.ErrID_JsonUnmarshal, "")
		mizugos.Warn(defines.LogSystem, this.name).Caller(0).EndError(errs.Errore(msgs.ErrID_JsonUnmarshal, err))
		return
	} // if

	submit, errID, err := miscs.HandleDatabase(defines.RedmoMixed, defines.MongoTable)

	if err != nil {
		this.sendMLoginA(msg, errID, "")
		mizugos.Error(defines.LogSystem, this.name).Caller(0).EndError(err)
		return
	} // if

	authGet := access.NewAuthGet(msg.Account, nil)

	if err = submit.Lock(msg.Account).Add(authGet).Exec(); err != nil {
		this.sendMLoginA(msg, msgs.ErrID_SubmitFailed, "")
		mizugos.Error(defines.LogSystem, this.name).Caller(0).EndError(errs.Errore(msgs.ErrID_SubmitFailed, err))
		return
	} // if

	authSet := access.NewAuthSet(msg.Account, &access.Auth{
		Account: msg.Account,
		Token:   uuid.NewString(),
		Time:    time.Now(),
	})

	if err = submit.Add(authSet).Unlock(msg.Account).Exec(); err != nil {
		this.sendMLoginA(msg, msgs.ErrID_SubmitFailed, "")
		mizugos.Error(defines.LogSystem, this.name).Caller(0).EndError(errs.Errore(msgs.ErrID_SubmitFailed, err))
		return
	} // if

	this.sendMLoginA(msg, msgs.ErrID_Success, authSet.Data.Token)
	mizugos.Info(defines.LogSystem, this.name).Caller(0).KV("account", authSet.Data.Account).KV("token", authSet.Data.Token).End()
}

// sendMLoginA 傳送回應登入
func (this *Auth) sendMLoginA(from *msgs.MLoginQ, errID msgs.ErrID, token string) {
	msg, err := procs.JsonMarshal(procs.MessageID(msgs.MsgID_LoginA), &msgs.MLoginA{
		From:  from,
		ErrID: int(errID),
		Token: token,
	})

	if err != nil {
		mizugos.Warn(defines.LogSystem, this.name).Caller(0).EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}

// procMUpdateQ 處理要求更新
func (this *Auth) procMUpdateQ(message any) {
	rec := features.Update.Rec()
	defer rec()

	_, msg, err := procs.JsonUnmarshal[msgs.MUpdateQ](message)

	if err != nil {
		this.sendMUpdateA(nil, msgs.ErrID_JsonUnmarshal, "")
		mizugos.Warn(defines.LogSystem, this.name).Caller(0).EndError(errs.Errore(msgs.ErrID_JsonUnmarshal, err))
		return
	} // if

	submit, errID, err := miscs.HandleDatabase(defines.RedmoMixed, defines.MongoTable)

	if err != nil {
		this.sendMUpdateA(msg, errID, "")
		mizugos.Error(defines.LogSystem, this.name).Caller(0).EndError(err)
		return
	} // if

	authGet := access.NewAuthGet(msg.Account, nil)

	if err = submit.Lock(msg.Account).Add(authGet).Exec(); err != nil {
		this.sendMUpdateA(msg, msgs.ErrID_SubmitFailed, "")
		mizugos.Error(defines.LogSystem, this.name).Caller(0).EndError(errs.Errore(msgs.ErrID_SubmitFailed, err))
		return
	} // if

	if authGet.Result == false {
		this.sendMUpdateA(msg, msgs.ErrID_AccountNotExist, "")
		mizugos.Warn(defines.LogSystem, this.name).Caller(0).EndError(errs.Errore(msgs.ErrID_AccountNotExist, err))
		return
	} // if

	if authGet.Data.Token != msg.Token {
		this.sendMUpdateA(msg, msgs.ErrID_TokenNotMatch, "")
		mizugos.Warn(defines.LogSystem, this.name).Caller(0).EndError(errs.Errore(msgs.ErrID_TokenNotMatch, err))
		return
	} // if

	authSet := access.NewAuthSet(msg.Account, &access.Auth{
		Account: msg.Account,
		Token:   uuid.NewString(),
		Time:    time.Now(),
	})

	if err = submit.Add(authSet).Unlock(msg.Account).Exec(); err != nil {
		this.sendMUpdateA(msg, msgs.ErrID_SubmitFailed, "")
		mizugos.Error(defines.LogSystem, this.name).Caller(0).EndError(errs.Errore(msgs.ErrID_SubmitFailed, err))
		return
	} // if

	this.sendMUpdateA(msg, msgs.ErrID_Success, authSet.Data.Token)
	mizugos.Info(defines.LogSystem, this.name).Caller(0).KV("account", authSet.Data.Account).KV("token", authSet.Data.Token).End()
}

// sendMUpdateA 傳送回應登入
func (this *Auth) sendMUpdateA(from *msgs.MUpdateQ, errID msgs.ErrID, token string) {
	msg, err := procs.JsonMarshal(procs.MessageID(msgs.MsgID_UpdateA), &msgs.MUpdateA{
		From:  from,
		ErrID: int(errID),
		Token: token,
	})

	if err != nil {
		mizugos.Warn(defines.LogSystem, this.name).Caller(0).EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}

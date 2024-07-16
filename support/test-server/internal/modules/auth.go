package modules

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/querys"
	"github.com/yinweli/Mizugo/support/test-server/msgs"
)

// NewAuth 建立Auth模組
func NewAuth() *Auth {
	return &Auth{
		Module: entitys.NewModule(defines.ModuleIDAuth),
		name:   "auth",
	}
}

// Auth Auth模組
type Auth struct {
	*entitys.Module
	name string // 系統名稱
}

// Awake 喚醒處理
func (this *Auth) Awake() error {
	this.Entity().AddMessage(int32(msgs.MsgID_LoginQ), this.procMLoginQ)
	this.Entity().AddMessage(int32(msgs.MsgID_UpdateQ), this.procMUpdateQ)
	return nil
}

// procMLoginQ 處理要求登入
func (this *Auth) procMLoginQ(message any) {
	rec := features.MeterLogin.Rec()
	defer rec()
	_, msg, err := procs.JsonUnmarshal[msgs.MLoginQ](message)

	if err != nil {
		this.sendMLoginA(nil, msgs.ErrID_JsonUnmarshal, "")
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("auth procMLoginQ: %w", err)).EndFlush()
		return
	} // if

	auth := querys.NewAuth(msg.Account)

	if err = features.DBMixed.Submit(context.Background()).Lock(msg.Account).Add(auth.NewGetter()).Exec(); err != nil {
		this.sendMLoginA(msg, msgs.ErrID_SubmitFailed, "")
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("auth procMLoginQ: %w", err)).EndFlush()
		return
	} // if

	auth.Token = uuid.NewString()
	auth.Time = time.Now()
	auth.SetSave()

	if err = features.DBMixed.Submit(context.Background()).Add(auth.NewSetter()).Unlock(msg.Account).Exec(); err != nil {
		this.sendMLoginA(msg, msgs.ErrID_SubmitFailed, "")
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("auth procMLoginQ: %w", err)).EndFlush()
		return
	} // if

	this.sendMLoginA(msg, msgs.ErrID_Success, auth.Token)
	features.LogSystem.Get().Info(this.name).KV("account", auth.Account).KV("token", auth.Token).Caller(0).EndFlush()
}

// sendMLoginA 傳送回應登入
func (this *Auth) sendMLoginA(from *msgs.MLoginQ, errID msgs.ErrID, token string) {
	msg, err := procs.JsonMarshal(int32(msgs.MsgID_LoginA), &msgs.MLoginA{
		From:  from,
		ErrID: int(errID),
		Token: token,
	})

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(err).EndFlush()
		return
	} // if

	this.Entity().Send(msg)
}

// procMUpdateQ 處理要求更新
func (this *Auth) procMUpdateQ(message any) {
	rec := features.MeterUpdate.Rec()
	defer rec()
	_, msg, err := procs.JsonUnmarshal[msgs.MUpdateQ](message)

	if err != nil {
		this.sendMUpdateA(nil, msgs.ErrID_JsonUnmarshal, "")
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("auth procMUpdateQ: %w", err)).EndFlush()
		return
	} // if

	auth := querys.NewAuth(msg.Account)

	if err = features.DBMixed.Submit(context.Background()).Lock(msg.Account).Add(auth.NewGetter()).Exec(); err != nil {
		this.sendMUpdateA(msg, msgs.ErrID_SubmitFailed, "")
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("auth procMUpdateQ: %w", err)).EndFlush()
		return
	} // if

	if auth.Token != msg.Token {
		this.sendMUpdateA(msg, msgs.ErrID_TokenNotMatch, "")
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("auth procMUpdateQ: %w", err)).EndFlush()
		return
	} // if

	auth.Token = uuid.NewString()
	auth.Time = time.Now()
	auth.SetSave()

	if err = features.DBMixed.Submit(context.Background()).Add(auth.NewSetter()).Unlock(msg.Account).Exec(); err != nil {
		this.sendMUpdateA(msg, msgs.ErrID_SubmitFailed, "")
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("auth procMUpdateQ: %w", err)).EndFlush()
		return
	} // if

	this.sendMUpdateA(msg, msgs.ErrID_Success, auth.Token)
	features.LogSystem.Get().Info(this.name).KV("account", auth.Account).KV("token", auth.Token).Caller(0).EndFlush()
}

// sendMUpdateA 傳送回應登入
func (this *Auth) sendMUpdateA(from *msgs.MUpdateQ, errID msgs.ErrID, token string) {
	msg, err := procs.JsonMarshal(int32(msgs.MsgID_UpdateA), &msgs.MUpdateA{
		From:  from,
		ErrID: int(errID),
		Token: token,
	})

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(err).EndFlush()
		return
	} // if

	this.Entity().Send(msg)
}

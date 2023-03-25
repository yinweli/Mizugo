package modules

import (
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/errs"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/msgs"
)

const nameProto = "module-proto" // 模組名稱

// NewProto 建立Proto模組
func NewProto(incr func() int64) *Proto {
	return &Proto{
		Module: entitys.NewModule(defines.ModuleIDProto),
		incr:   incr,
	}
}

// Proto Proto模組
type Proto struct {
	*entitys.Module
	incr func() int64 // 計數函式
}

// Awake 喚醒處理
func (this *Proto) Awake() error {
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_ProtoQ), this.procMProtoQ)
	return nil
}

// procMProtoQ 處理要求Proto
func (this *Proto) procMProtoQ(message any) {
	rec := features.Proto.Rec()
	defer rec()

	_, msg, err := procs.ProtoUnmarshal[msgs.MProtoQ](message)

	if err != nil {
		this.sendMProtoA(nil, msgs.ErrID_ProtoUnmarshal, 0)
		features.System.Warn(nameProto).Caller(0).EndError(errs.Errore(msgs.ErrID_ProtoUnmarshal, err))
		return
	} // if

	count := this.incr()
	this.sendMProtoA(msg, msgs.ErrID_Success, count)
	features.System.Info(nameProto).Caller(0).KV("count", count).End()
}

// sendMProtoA 傳送回應Proto
func (this *Proto) sendMProtoA(from *msgs.MProtoQ, errID msgs.ErrID, count int64) {
	msg, err := procs.ProtoMarshal(procs.MessageID(msgs.MsgID_ProtoA), &msgs.MProtoA{
		From:  from,
		ErrID: errID,
		Count: count,
	})

	if err != nil {
		features.System.Warn(nameProto).Caller(0).EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}

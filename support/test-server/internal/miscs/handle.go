package miscs

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/errs"
	"github.com/yinweli/Mizugo/mizugos/redmos"
	"github.com/yinweli/Mizugo/support/test-server/internal/msgs"
)

// HandleDatabase 資料庫處理工具
func HandleDatabase(mixedName, tableName string) (submit *redmos.Submit, id msgs.ErrID, err error) {
	var database *redmos.Mixed

	if database = mizugos.Redmomgr().GetMixed(mixedName); database == nil {
		return nil, msgs.ErrID_DatabaseNil, errs.Errort(msgs.ErrID_DatabaseNil)
	} // if

	if submit = database.Submit(ctxs.Root(), tableName); submit == nil {
		return nil, msgs.ErrID_SubmitNil, errs.Errort(msgs.ErrID_SubmitNil)
	} // if

	return submit, msgs.ErrID_Success, nil
}

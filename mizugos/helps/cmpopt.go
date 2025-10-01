package helps

import (
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProtoTimestampWithin 回傳一個 cmp.Option, 當兩個 google.protobuf.Timestamp 的時間差在 margin 以內時視為相等(margin < 0 會 panic)
func ProtoTimestampWithin(margin time.Duration) cmp.Option {
	if margin < 0 {
		panic("margin must be a non-negative number")
	} // if

	// 用來檢查 proto message 是否 timestamppb.Timestamp
	isTimestamp := func(v any) (protocmp.Message, bool) {
		p, ok := v.(protocmp.Message)
		return p, ok && p.Descriptor().FullName() == "google.protobuf.Timestamp"
	}
	// 用來轉換 proto message 為 timestamppb.Timestamp
	toTimestamp := func(v any) *timestamppb.Timestamp {
		if p, ok := isTimestamp(v); ok {
			if t, ok := p.Unwrap().(*timestamppb.Timestamp); ok {
				return t
			} // if
		} // if

		return nil
	}

	return cmp.FilterValues(
		func(l, r any) bool {
			_, lc := isTimestamp(l)
			_, rc := isTimestamp(r)
			return lc && rc
		},
		cmp.Comparer(func(l, r any) bool {
			ls, rs := toTimestamp(l), toTimestamp(r)

			if ls == nil || rs == nil {
				return false
			} // if

			lt, rt := ls.AsTime(), rs.AsTime()

			if lt.After(rt) {
				lt, rt = rt, lt
			} // if

			return rt.Sub(lt) <= margin
		}),
	)
}

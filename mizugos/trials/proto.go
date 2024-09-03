package trials

import (
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProtoEqual 比對訊息是否一致
func ProtoEqual(expected, actual proto.Message, option ...cmp.Option) bool {
	if cmp.Equal(expected, actual, append(option, protocmp.Transform())...) == false {
		fmt.Printf("expected: %v\n", protojson.Format(expected))
		fmt.Printf("actual: %v\n", protojson.Format(actual))
		return false
	} // if

	return true
}

// EquateApproxTimestamp 取得用於比較兩個 timestamppb.Timestamp 的比對選項
//
// margin 是允許的最大時間誤差, 必須是非負數, 否則會觸發 panic
//
// 當兩個時間戳之間的差異在指定的 margin 範圍內時, 視為相等
func EquateApproxTimestamp(margin time.Duration) cmp.Option {
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

			return lt.Add(margin).Before(rt) == false
		}),
	)
}

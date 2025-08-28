package redmos

import (
	"github.com/yinweli/Mizugo/v2/testdata"
)

// testMeta 元資料測試器
type testMeta struct {
}

func (this *testMeta) MajorKey(_ any) string {
	return ""
}

func (this *testMeta) MinorKey(_ any) string {
	return ""
}

func (this *testMeta) MinorTable() string {
	return testdata.Unknown
}

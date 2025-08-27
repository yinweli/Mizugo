package helps

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestReflect(t *testing.T) {
	suite.Run(t, new(SuiteReflect))
}

type SuiteReflect struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteReflect) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-reflect"))
}

func (this *SuiteReflect) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteReflect) TestReflect() {
	target, ok := ReflectFieldValue[string](reflect.ValueOf(testReflect{
		Export:   testdata.Unknown,
		unexport: testdata.Unknown,
	}), "Export")
	this.True(ok)
	this.Equal(testdata.Unknown, target)

	target, ok = ReflectFieldValue[string](reflect.ValueOf(&testReflect{
		Export:   testdata.Unknown,
		unexport: testdata.Unknown,
	}), "Export")
	this.True(ok)
	this.Equal(testdata.Unknown, target)

	_, ok = ReflectFieldValue[string](reflect.ValueOf(&testReflect{
		Export:   testdata.Unknown,
		unexport: testdata.Unknown,
	}), "unexport")
	this.False(ok)

	_, ok = ReflectFieldValue[string](reflect.ValueOf(&testReflect{
		Export:   testdata.Unknown,
		unexport: testdata.Unknown,
	}), testdata.Unknown)
	this.False(ok)

	_, ok = ReflectFieldValue[string](reflect.ValueOf(testdata.Unknown), testdata.Unknown)
	this.False(ok)

	_, ok = ReflectFieldValue[string](reflect.ValueOf(nil), testdata.Unknown)
	this.False(ok)
}

type testReflect struct {
	Export   string
	unexport string
}

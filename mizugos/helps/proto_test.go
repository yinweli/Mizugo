package helps

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/v2/mizugos/msgs"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestProto(t *testing.T) {
	suite.Run(t, new(SuiteProto))
}

type SuiteProto struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteProto) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-proto"))
}

func (this *SuiteProto) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteProto) TestToProtoAny() {
	target, err := ToProtoAny(&msgs.ProtoTest{}, &msgs.ProtoTest{}, &msgs.ProtoTest{})
	this.Nil(err)
	this.Len(target, 3)

	_, err = ToProtoAny(nil)
	this.NotNil(err)
}

func (this *SuiteProto) TestFromProtoAny() {
	message, _ := anypb.New(&msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	target, err := FromProtoAny[*msgs.ProtoTest](message)
	this.Nil(err)
	this.NotNil(target)
	this.Equal(testdata.Unknown, target.Data)

	_, err = FromProtoAny[*msgs.ProtoTest](nil)
	this.NotNil(err)

	_, err = FromProtoAny[*msgs.Proto](message)
	this.NotNil(err)
}

func (this *SuiteProto) TestProtoString() {
	this.NotEmpty(ProtoString(&msgs.ProtoTest{}))
	this.NotEmpty(ProtoString(nil))
}

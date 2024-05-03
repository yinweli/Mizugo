package helps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
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

func (this *SuiteProto) TestToProto() {
	result, err := ToProtoAny(&msgs.ProtoTest{}, &msgs.ProtoTest{}, &msgs.ProtoTest{})
	assert.Nil(this.T(), err)
	assert.Len(this.T(), result, 3)
	_, err = ToProtoAny(nil)
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestFromProtoAny() {
	input, _ := anypb.New(&msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	output, err := FromProtoAny[msgs.ProtoTest](input)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output)
	assert.Equal(this.T(), testdata.Unknown, output.Data)
	_, err = FromProtoAny[msgs.ProtoTest](nil)
	assert.NotNil(this.T(), err)
	_, err = FromProtoAny[int](input)
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestProtoString() {
	assert.NotNil(this.T(), ProtoString(&msgs.ProtoTest{}))
	assert.NotNil(this.T(), ProtoString(nil))
}

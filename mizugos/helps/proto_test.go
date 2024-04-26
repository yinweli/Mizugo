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

func (this *SuiteProto) TestProtoAny() {
	temp, _ := anypb.New(&msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	output, err := ProtoAny[msgs.ProtoTest](temp)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output)
	assert.Equal(this.T(), testdata.Unknown, output.Data)
	_, err = ProtoAny[msgs.ProtoTest](nil)
	assert.NotNil(this.T(), err)
	_, err = ProtoAny[int](temp)
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestProtoJson() {
	assert.NotNil(this.T(), ProtoJson(&msgs.ProtoTest{
		Data: testdata.Unknown,
	}))
	assert.NotNil(this.T(), ProtoJson(nil))
}

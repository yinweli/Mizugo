package trials

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestProto(t *testing.T) {
	suite.Run(t, new(SuiteProto))
}

type SuiteProto struct {
	suite.Suite
}

func (this *SuiteProto) TestProtoEqual() {
	assert.True(this.T(), ProtoEqual(&msgs.ProtoTest{Data: testdata.Unknown}, &msgs.ProtoTest{Data: testdata.Unknown}))
	assert.False(this.T(), ProtoEqual(&msgs.ProtoTest{Data: testdata.Unknown}, &msgs.ProtoTest{}))
}

func (this *SuiteProto) TestProtoContains() {
	assert.True(this.T(), ProtoContains([]proto.Message{&msgs.ProtoTest{}}, (*msgs.ProtoTest)(nil)))
	assert.False(this.T(), ProtoContains([]proto.Message{&msgs.ProtoTest{}}, (*msgs.Proto)(nil)))
}

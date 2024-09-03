package trials

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"

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
	now := time.Now()
	assert.True(this.T(), ProtoEqual(
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now)},
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now)}))
	assert.False(this.T(), ProtoEqual(
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now)},
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now.Add(time.Millisecond))}))
	assert.True(this.T(), ProtoEqual(
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now)},
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now.Add(time.Millisecond))},
		EquateApproxTimestamp(time.Second)))
	assert.False(this.T(), ProtoEqual(
		&msgs.ProtoTest{Data: testdata.Unknown},
		&msgs.ProtoTest{}))
}

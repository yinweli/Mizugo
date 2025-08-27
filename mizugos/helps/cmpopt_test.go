package helps

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yinweli/Mizugo/v2/mizugos/msgs"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestCmpOpt(t *testing.T) {
	suite.Run(t, new(SuiteCmpOpt))
}

type SuiteCmpOpt struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteCmpOpt) TestProtoTimestampWithin() {
	now := time.Now()
	this.True(trials.ProtoEqual(
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now)},
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now.Add(time.Millisecond))},
		ProtoTimestampWithin(time.Second),
	))
	this.False(trials.ProtoEqual(
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now)},
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now.Add(time.Millisecond))},
	))
}

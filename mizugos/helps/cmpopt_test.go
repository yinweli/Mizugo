package helps

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmpOpt(t *testing.T) {
	suite.Run(t, new(SuiteCmpOpt))
}

type SuiteCmpOpt struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteCmpOpt) TestEquateApproxProtoTimestamp() {
	now := time.Now()
	assert.True(this.T(), trials.ProtoEqual(
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now)},
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now.Add(time.Millisecond))},
		EquateApproxProtoTimestamp(time.Second)))
	assert.False(this.T(), trials.ProtoEqual(
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now)},
		&msgs.ProtoTest{Data: testdata.Unknown, Time: timestamppb.New(now.Add(time.Millisecond))}))
}

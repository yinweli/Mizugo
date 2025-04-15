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
	assert.False(this.T(), ProtoEqual(&msgs.ProtoTest{}, testdata.Unknown))
	assert.False(this.T(), ProtoEqual(testdata.Unknown, &msgs.ProtoTest{}))
}

func (this *SuiteProto) TestProtoListEqual() {
	assert.True(this.T(), ProtoListEqual(
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
			{Data: testdata.Unknown},
		},
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
			{Data: testdata.Unknown},
		}))
	assert.False(this.T(), ProtoListEqual(
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
			{Data: testdata.Unknown},
		},
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
		}))
	assert.False(this.T(), ProtoListEqual(
		[]string{
			testdata.Unknown,
		},
		[]string{
			"12345",
		}))
}

func (this *SuiteProto) TestProtoListContain() {
	assert.True(this.T(), ProtoListContain(&msgs.ProtoTest{Data: testdata.Unknown}, []*msgs.ProtoTest{
		nil,
		{Data: ""},
		{Data: "12345"},
		{Data: testdata.Unknown},
	}))
	assert.False(this.T(), ProtoListContain(&msgs.ProtoTest{Data: testdata.Unknown}, []*msgs.ProtoTest{
		nil,
		{Data: ""},
		{Data: "12345"},
	}))
	assert.False(this.T(), ProtoListContain(&msgs.ProtoTest{Data: testdata.Unknown}, []string{
		testdata.Unknown,
	}))
	assert.False(this.T(), ProtoListContain(testdata.Unknown, []*msgs.ProtoTest{
		nil,
		{Data: ""},
		{Data: "12345"},
	}))
}

func (this *SuiteProto) TestProtoListMatch() {
	assert.True(this.T(), ProtoListMatch(
		[]*msgs.ProtoTest{
			{Data: testdata.Unknown},
		},
		[]*msgs.ProtoTest{
			nil,
			{Data: ""},
			{Data: "12345"},
			{Data: testdata.Unknown},
		}))
	assert.False(this.T(), ProtoListMatch(
		[]*msgs.ProtoTest{
			{Data: testdata.Unknown},
		},
		[]*msgs.ProtoTest{
			nil,
			{Data: ""},
			{Data: "12345"},
		}))
}

func (this *SuiteProto) TestProtoListExist() {
	assert.True(this.T(), ProtoListExist((*msgs.ProtoTest)(nil), []proto.Message{&msgs.ProtoTest{}}))
	assert.False(this.T(), ProtoListExist((*msgs.Proto)(nil), []proto.Message{&msgs.ProtoTest{}}))
}

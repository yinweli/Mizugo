package trials

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/yinweli/Mizugo/v2/mizugos/msgs"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestProto(t *testing.T) {
	suite.Run(t, new(SuiteProto))
}

type SuiteProto struct {
	suite.Suite
}

func (this *SuiteProto) TestProtoEqual() {
	this.True(ProtoEqual(&msgs.ProtoTest{Data: testdata.Unknown}, &msgs.ProtoTest{Data: testdata.Unknown}))
	this.False(ProtoEqual(&msgs.ProtoTest{Data: testdata.Unknown}, &msgs.ProtoTest{}))
	this.False(ProtoEqual(&msgs.ProtoTest{}, testdata.Unknown))
	this.False(ProtoEqual(testdata.Unknown, &msgs.ProtoTest{}))
}

func (this *SuiteProto) TestProtoListEqual() {
	this.True(ProtoListEqual(
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
			{Data: testdata.Unknown},
		},
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
			{Data: testdata.Unknown},
		},
	))
	this.False(ProtoListEqual(
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
			{Data: testdata.Unknown},
		},
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
		},
	))
	this.False(ProtoListEqual(
		[]string{
			testdata.Unknown + "expected",
		},
		[]string{
			testdata.Unknown + "actual",
		},
	))
}

func (this *SuiteProto) TestProtoListMatch() {
	this.True(ProtoListMatch(
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
			{Data: testdata.Unknown},
		},
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
			{Data: testdata.Unknown},
		},
	))
	this.True(ProtoListMatch(
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
			{Data: testdata.Unknown},
		},
		[]*msgs.ProtoTest{
			{Data: testdata.Unknown},
			{Data: "12345"},
			{Data: ""},
		},
	))
	this.False(ProtoListMatch(
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
			{Data: testdata.Unknown},
		},
		[]*msgs.ProtoTest{
			{Data: ""},
			{Data: "12345"},
		},
	))
	this.False(ProtoListMatch(
		[]string{
			testdata.Unknown + "expected",
		},
		[]string{
			testdata.Unknown + "actual",
		},
	))
}

func (this *SuiteProto) TestProtoListHasData() {
	this.True(ProtoListHasData(&msgs.ProtoTest{Data: testdata.Unknown}, []*msgs.ProtoTest{
		{Data: ""},
		{Data: "12345"},
		{Data: testdata.Unknown},
	}))
	this.False(ProtoListHasData(&msgs.ProtoTest{Data: testdata.Unknown}, []*msgs.ProtoTest{
		{Data: ""},
		{Data: "12345"},
	}))
	this.False(ProtoListHasData(&msgs.ProtoTest{Data: testdata.Unknown}, []string{
		testdata.Unknown,
	}))
	this.False(ProtoListHasData(testdata.Unknown, []*msgs.ProtoTest{
		nil,
		{Data: ""},
		{Data: "12345"},
	}))
}

func (this *SuiteProto) TestProtoListHasType() {
	this.True(ProtoListHasType((*msgs.ProtoTest)(nil), []proto.Message{&msgs.ProtoTest{}}))
	this.False(ProtoListHasType((*msgs.Proto)(nil), []proto.Message{&msgs.ProtoTest{}}))
}

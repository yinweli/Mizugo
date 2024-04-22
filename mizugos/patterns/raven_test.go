package patterns

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestRaven(t *testing.T) {
	suite.Run(t, new(SuiteRaven))
}

type SuiteRaven struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteRaven) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-patterns-raven"))
}

func (this *SuiteRaven) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteRaven) TestRavenQ() {
	messageID := procs.MessageID(1)
	header := &msgs.RavenTest{
		Data: "header",
	}
	request := &msgs.RavenTest{
		Data: "request",
	}
	output1, err := RavenQBuilder(messageID, header, request)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output1)
	output2, err := RavenQParser[msgs.RavenTest, msgs.RavenTest](output1)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output2)
	assert.Equal(this.T(), messageID, output2.MessageID)
	assert.True(this.T(), proto.Equal(header, output2.Header))
	assert.True(this.T(), proto.Equal(request, output2.Request))
	_, err = RavenQParser[msgs.RavenTest, msgs.RavenTest](nil)
	assert.NotNil(this.T(), err)
	detail := output2.Detail()
	assert.NotNil(this.T(), detail)
	fmt.Println(detail)
}

func (this *SuiteRaven) TestRavenA() {
	messageID := procs.MessageID(1)
	errID := int32(2)
	header := &msgs.RavenTest{
		Data: "header",
	}
	request := &msgs.RavenTest{
		Data: "request",
	}
	respond1 := &msgs.RavenTest{
		Data: "respond1",
	}
	respond2 := &msgs.RavenTest{
		Data: "respond2",
	}
	output1, err := RavenABuilder(messageID, errID, header, request, respond1, respond2)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output1)
	output2, err := RavenAParser[msgs.RavenTest, msgs.RavenTest](output1)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output2)
	assert.Equal(this.T(), messageID, output2.MessageID)
	assert.Equal(this.T(), errID, output2.ErrID)
	assert.True(this.T(), proto.Equal(header, output2.Header))
	assert.True(this.T(), proto.Equal(request, output2.Request))
	assert.True(this.T(), proto.Equal(respond1, output2.GetRespond(respond1)))
	assert.Nil(this.T(), output2.GetRespond(nil))
	assert.True(this.T(), proto.Equal(respond2, output2.GetRespondAt(1)))
	assert.Nil(this.T(), output2.GetRespondAt(10))
	_, err = RavenAParser[msgs.RavenTest, msgs.RavenTest](nil)
	assert.NotNil(this.T(), err)
	detail := output2.Detail()
	assert.NotNil(this.T(), detail)
	fmt.Println(detail)
}

func (this *SuiteRaven) TestRavenTest() {
	messageID := procs.MessageID(1)
	errID := int32(2)
	header := &msgs.RavenTest{
		Data: "header",
	}
	request := &msgs.RavenTest{
		Data: "request",
	}
	respond1 := &msgs.RavenTest{
		Data: "respond1",
	}
	respond2 := &msgs.RavenTest{
		Data: "respond2",
	}
	ravenQ, _ := RavenQBuilder(messageID, header, request)
	ravenA, _ := RavenABuilder(messageID, errID, header, request, respond1, respond2)
	assert.True(this.T(), RavenTestMessageID(ravenQ, messageID))
	assert.True(this.T(), RavenTestMessageID(ravenA, messageID))
	assert.False(this.T(), RavenTestMessageID(ravenA, procs.MessageID(0)))
	assert.False(this.T(), RavenTestMessageID(nil, procs.MessageID(0)))
	assert.True(this.T(), RavenTestErrID(ravenA, errID))
	assert.False(this.T(), RavenTestErrID(nil, errID))
	assert.True(this.T(), RavenTestHeader(ravenQ, header))
	assert.True(this.T(), RavenTestHeader(ravenA, header))
	assert.False(this.T(), RavenTestHeader(ravenA, nil))
	assert.False(this.T(), RavenTestHeader(nil, header))
	assert.True(this.T(), RavenTestRequest(ravenQ, request))
	assert.True(this.T(), RavenTestRequest(ravenA, request))
	assert.False(this.T(), RavenTestRequest(ravenA, nil))
	assert.False(this.T(), RavenTestRequest(nil, request))
	assert.True(this.T(), RavenTestRespond(ravenA, respond1, respond2))
	assert.False(this.T(), RavenTestRespond(ravenA, respond1, respond2, respond2))
	assert.False(this.T(), RavenTestRespond(nil, respond1, respond2))
	assert.True(this.T(), RavenTestRespondType(ravenA, &msgs.RavenTest{}, &msgs.RavenTest{}))
	assert.False(this.T(), RavenTestRespondType(ravenA, &msgs.RavenTest{}, &msgs.RavenTest{}, &msgs.RavenTest{}))
	assert.False(this.T(), RavenTestRespondType(nil, &msgs.RavenTest{}, &msgs.RavenTest{}))
	assert.True(this.T(), RavenTestRespondLength(ravenA, 2))
	assert.False(this.T(), RavenTestRespondLength(ravenA, 3))
	assert.False(this.T(), RavenTestRespondLength(nil, 3))
}

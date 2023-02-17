package depots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/contexts"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMinor(t *testing.T) {
	suite.Run(t, new(SuiteMinor))
}

type SuiteMinor struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
	testdata.TestDB
	name string
}

func (this *SuiteMinor) SetupSuite() {
	this.Change("test-depots-minor")
	this.name = "minor"
}

func (this *SuiteMinor) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMinor) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMinor) TestNewMinor() {
	target, err := newMinor(contexts.Ctx(), testdata.MongoURI)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)
}

func (this *SuiteMinor) TestMinor() {
	target, err := newMinor(contexts.Ctx(), testdata.MongoURI)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target.Runner(this.name, this.name))
	assert.NotNil(this.T(), target.Client())
	target.stop(contexts.Ctx())
	assert.Nil(this.T(), target.Runner(this.name, this.name))
	assert.Nil(this.T(), target.Client())

	_, err = newMinor(contexts.Ctx(), testdata.MongoURIInvalid)
	assert.NotNil(this.T(), err)
}

func (this *SuiteMinor) TestClient() {
	this.Reset()

	target, err := newMinor(contexts.Ctx(), testdata.MongoURI)
	assert.Nil(this.T(), err)
	client := target.Client()
	assert.NotNil(this.T(), client)

	data := &minorTester{
		Key:  this.Key("client"),
		Data: utils.RandString(testdata.RandStringLength),
	}
	table := client.Database(this.name).Collection(this.name)
	assert.NotNil(this.T(), table)

	_, err = table.InsertOne(contexts.Ctx(), data)
	assert.Nil(this.T(), err)

	_, err = table.DeleteOne(contexts.Ctx(), data)
	assert.Nil(this.T(), err)

	this.MongoClear(contexts.Ctx(), table)
	target.stop(contexts.Ctx())
}

type minorTester struct {
	Key  string `bson:"key"`
	Data string `bson:"data"`
}

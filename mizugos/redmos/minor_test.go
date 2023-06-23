package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMinor(t *testing.T) {
	suite.Run(t, new(SuiteMinor))
}

type SuiteMinor struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteMinor) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-minor")
}

func (this *SuiteMinor) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteMinor) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMinor) TestMinor() {
	target, err := newMinor(testdata.MongoURI, "minor")
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.Submit())
	assert.NotNil(this.T(), target.Client())
	assert.NotNil(this.T(), target.Database())
	assert.Nil(this.T(), target.SwitchDB("minor"))
	target.DropDB()

	_, err = newMinor("", "minor")
	assert.NotNil(this.T(), err)

	_, err = newMinor(testdata.MongoURI, "")
	assert.NotNil(this.T(), err)

	assert.Nil(this.T(), target.Client().Ping(ctxs.Get().Ctx(), nil))

	target.stop()
	assert.Nil(this.T(), target.Submit())
	assert.Nil(this.T(), target.Client())
	assert.Nil(this.T(), target.Database())
	assert.NotNil(this.T(), target.SwitchDB("minor"))
	target.DropDB()

	_, err = newMinor(testdata.MongoURIInvalid, "minor")
	assert.NotNil(this.T(), err)
}

func (this *SuiteMinor) TestMinorSubmit() {
	target, _ := newMinor(testdata.MongoURI, "minor")
	submit := target.Submit()
	assert.NotNil(this.T(), submit)
	assert.NotNil(this.T(), submit.Table("minor"))
	assert.NotNil(this.T(), submit.Database())
	target.stop()
}

func BenchmarkMinorSet(b *testing.B) {
	type myData struct {
		Key  string `bson:"key"`
		Data string `bson:"data"`
	}

	name := "benchmark minor"
	data := &myData{
		Key:  utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
		Data: utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
	}
	target, _ := newMinor(testdata.MongoURI, name)
	submit := target.Submit()

	for i := 0; i < b.N; i++ {
		data.Data = utils.RandString(testdata.RandStringLength, testdata.RandStringLetter)
		_, _ = submit.Table(name).ReplaceOne(ctxs.Get().Ctx(), bson.D{{Key: "key", Value: data.Key}}, data, options.Replace().SetUpsert(true))
	} // for

	_ = submit.Drop(ctxs.Get().Ctx())
}

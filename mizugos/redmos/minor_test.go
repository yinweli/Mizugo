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
	name string
}

func (this *SuiteMinor) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-redmos-minor")
	this.name = "minor"
}

func (this *SuiteMinor) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
}

func (this *SuiteMinor) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMinor) TestNewMinor() {
	target, err := newMinor(ctxs.RootCtx(), testdata.MongoURI, this.name)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)

	_, err = newMinor(ctxs.RootCtx(), "", this.name)
	assert.NotNil(this.T(), err)

	_, err = newMinor(ctxs.RootCtx(), testdata.MongoURI, "")
	assert.NotNil(this.T(), err)
}

func (this *SuiteMinor) TestMinor() {
	target, err := newMinor(ctxs.RootCtx(), testdata.MongoURI, this.name)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target.Submit())
	assert.NotNil(this.T(), target.Client())
	assert.NotNil(this.T(), target.Database())
	assert.Nil(this.T(), target.Client().Ping(ctxs.RootCtx(), nil))
	assert.NotNil(this.T(), target.Database().Client())
	target.stop(ctxs.RootCtx())
	assert.Nil(this.T(), target.Submit())
	assert.Nil(this.T(), target.Client())
	assert.Nil(this.T(), target.Database())

	_, err = newMinor(ctxs.RootCtx(), testdata.MongoURIInvalid, this.name)
	assert.NotNil(this.T(), err)
}

func (this *SuiteMinor) TestMinorSubmit() {
	target, err := newMinor(ctxs.RootCtx(), testdata.MongoURI, this.name)
	assert.Nil(this.T(), err)
	submit := target.Submit()
	assert.NotNil(this.T(), submit)
	assert.NotNil(this.T(), submit.Table(this.name))
	assert.NotNil(this.T(), submit.Database())
	target.stop(ctxs.RootCtx())
}

func BenchmarkMinorSet(b *testing.B) {
	type myData struct {
		Key  string `bson:"key"`
		Data string `bson:"data"`
	}

	name := "benchmark minor"
	data := &myData{
		Key:  utils.RandString(testdata.RandStringLength),
		Data: utils.RandString(testdata.RandStringLength),
	}
	target, _ := newMinor(ctxs.RootCtx(), testdata.MongoURI, name)
	submit := target.Submit()

	for i := 0; i < b.N; i++ {
		data.Data = utils.RandString(testdata.RandStringLength)
		_, _ = submit.Table(name).ReplaceOne(ctxs.RootCtx(), bson.D{{Key: "key", Value: data.Key}}, data, options.Replace().SetUpsert(true))
	} // for

	_ = submit.Drop(ctxs.RootCtx())
}

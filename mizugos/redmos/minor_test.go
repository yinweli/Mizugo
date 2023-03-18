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
	testdata.TestEnv
	testdata.TestLeak
	testdata.TestDB
	name string
}

func (this *SuiteMinor) SetupSuite() {
	this.Change("test-redmos-minor")
	this.name = "minor"
}

func (this *SuiteMinor) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMinor) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMinor) TestNewMinor() {
	target, err := newMinor(ctxs.Root(), testdata.MongoURI, this.name)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)

	_, err = newMinor(ctxs.Root(), "", this.name)
	assert.NotNil(this.T(), err)

	_, err = newMinor(ctxs.Root(), testdata.MongoURI, "")
	assert.NotNil(this.T(), err)
}

func (this *SuiteMinor) TestMinor() {
	target, err := newMinor(ctxs.Root(), testdata.MongoURI, this.name)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target.Submit())
	client := target.Client()
	assert.NotNil(this.T(), client)
	assert.Nil(this.T(), client.Ping(ctxs.RootCtx(), nil))
	database := target.Database()
	assert.NotNil(this.T(), database)
	assert.NotNil(this.T(), database.Client())
	target.stop(ctxs.Root())
	assert.Nil(this.T(), target.Submit())
	assert.Nil(this.T(), target.Client())
	assert.Nil(this.T(), target.Database())

	_, err = newMinor(ctxs.Root(), testdata.MongoURIInvalid, this.name)
	assert.NotNil(this.T(), err)
}

func (this *SuiteMinor) TestMinorSubmit() {
	target, err := newMinor(ctxs.Root(), testdata.MongoURI, this.name)
	assert.Nil(this.T(), err)
	submit := target.Submit()
	assert.NotNil(this.T(), submit)
	assert.NotNil(this.T(), submit.Table(this.name))
	assert.NotNil(this.T(), submit.Database())
	target.stop(ctxs.Root())
}

func BenchmarkMinorSet(b *testing.B) {
	name := "benchmark minor"
	target, _ := newMinor(ctxs.Root(), testdata.MongoURI, name)
	submit := target.Submit()

	for i := 0; i < b.N; i++ {
		value := utils.RandString(testdata.RandStringLength)
		_, _ = submit.Table(name).ReplaceOne(
			ctxs.RootCtx(),
			bson.D{{Key: "key", Value: value}},
			&dataTester{Key: value, Data: value},
			options.Replace().SetUpsert(true))
	} // for

	_ = submit.Drop(ctxs.RootCtx())
}

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
}

func (this *SuiteMinor) TestMinor() {
	target, err := newMinor(ctxs.Root(), testdata.MongoURI, this.name)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target.Submit(this.name))
	assert.NotNil(this.T(), target.Client())
	target.stop(ctxs.Root())
	assert.Nil(this.T(), target.Submit(this.name))
	assert.Nil(this.T(), target.Client())

	_, err = newMinor(ctxs.Root(), testdata.MongoURIInvalid, this.name)
	assert.NotNil(this.T(), err)
}

func (this *SuiteMinor) TestClient() {
	target, err := newMinor(ctxs.Root(), testdata.MongoURI, this.name)
	assert.Nil(this.T(), err)
	client := target.Client()
	assert.NotNil(this.T(), client)

	data := &dataTester{
		Key:  this.Key("minor client"),
		Data: utils.RandString(testdata.RandStringLength),
	}
	table := client.Database(this.name).Collection(this.name)
	assert.NotNil(this.T(), table)

	_, err = table.InsertOne(ctxs.RootCtx(), data)
	assert.Nil(this.T(), err)

	_, err = table.DeleteOne(ctxs.RootCtx(), data)
	assert.Nil(this.T(), err)

	this.MongoClear(ctxs.RootCtx(), table)
	target.stop(ctxs.Root())
}

func BenchmarkMinorSet(b *testing.B) {
	name := "benchmark minor"
	target, _ := newMinor(ctxs.Root(), testdata.MongoURI, name)
	submit := target.Submit(name)

	for i := 0; i < b.N; i++ {
		value := utils.RandString(testdata.RandStringLength)
		_, _ = submit.ReplaceOne(
			ctxs.RootCtx(),
			bson.D{{Key: "key", Value: value}},
			&dataTester{Key: value, Data: value},
			options.Replace().SetUpsert(true))
	} // for

	_ = submit.Drop(ctxs.RootCtx())
}

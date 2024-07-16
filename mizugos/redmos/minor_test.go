package redmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMinor(t *testing.T) {
	suite.Run(t, new(SuiteMinor))
}

type SuiteMinor struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteMinor) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-minor"))
}

func (this *SuiteMinor) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteMinor) TestMinor() {
	target, err := newMinor(testdata.MongoURI, "minor")
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.Submit())
	assert.NotNil(this.T(), target.Client())
	assert.NotNil(this.T(), target.Database())
	assert.Nil(this.T(), target.SwitchDB("minor"))
	assert.NotNil(this.T(), target.SwitchDB(""))
	target.DropDB()
	assert.Nil(this.T(), target.Client().Ping(context.Background(), nil))
	target.stop()
	assert.Nil(this.T(), target.Submit())
	assert.Nil(this.T(), target.Client())
	assert.Nil(this.T(), target.Database())
	assert.NotNil(this.T(), target.SwitchDB("minor"))
	target.DropDB()

	_, err = newMinor("", "minor")
	assert.NotNil(this.T(), err)
	_, err = newMinor(testdata.MongoURI, "")
	assert.NotNil(this.T(), err)
	_, err = newMinor(testdata.MongoURIInvalid, "minor")
	assert.NotNil(this.T(), err)
}

func (this *SuiteMinor) TestMinorSubmit() {
	minor, _ := newMinor(testdata.MongoURI, "minor")
	target := minor.Submit()
	assert.NotNil(this.T(), target.Collection("minor"))
	assert.NotNil(this.T(), target.Operate("minor", mongo.NewReplaceOneModel()))
	target = minor.Submit()
	assert.Nil(this.T(), target.Exec(context.Background()))
}

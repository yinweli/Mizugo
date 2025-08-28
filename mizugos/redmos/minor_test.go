package redmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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
	this.Nil(err)
	this.NotNil(target)
	target.stop()
	_, err = newMinor("", "minor")
	this.NotNil(err)
	_, err = newMinor(testdata.MongoURI, "")
	this.NotNil(err)
	_, err = newMinor(testdata.MongoURIInvalid, "minor")
	this.NotNil(err)
}

func (this *SuiteMinor) TestSubmit() {
	target, _ := newMinor(testdata.MongoURI, "minor")
	submit := target.Submit()
	this.NotNil(submit.Collection("minor"))
	this.NotNil(submit.Operate("minor", mongo.NewReplaceOneModel()))
	submit = target.Submit()
	this.Nil(submit.Exec(context.Background()))
	target.stop()
	this.Nil(target.Submit())
}

func (this *SuiteMinor) TestClient() {
	target, _ := newMinor(testdata.MongoURI, "minor")
	client := target.Client()
	this.NotNil(client)
	this.Nil(target.Client().Ping(context.Background(), nil))
	target.stop()
}

func (this *SuiteMinor) TestDatabase() {
	target, _ := newMinor(testdata.MongoURI, "minor")
	this.NotNil(target.Database())
	target.stop()
}

func (this *SuiteMinor) TestSwitchDB() {
	target, _ := newMinor(testdata.MongoURI, "minor")
	this.Nil(target.SwitchDB("minor"))
	this.NotNil(target.SwitchDB(""))
	target.stop()
	this.NotNil(target.SwitchDB("minor"))
}

func (this *SuiteMinor) TestDropDB() {
	target, _ := newMinor(testdata.MongoURI, "minor")
	target.DropDB()
	target.stop()
}

func (this *SuiteMinor) TestMinorIndex() {
	target := MinorIndex(&testMeta{})
	this.NotNil(target)
	this.NotEmpty(target.Name)
	this.NotEmpty(target.Table)
	this.NotEmpty(target.Field)
}

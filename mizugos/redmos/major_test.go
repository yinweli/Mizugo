package redmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestMajor(t *testing.T) {
	suite.Run(t, new(SuiteMajor))
}

type SuiteMajor struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteMajor) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-major"))
}

func (this *SuiteMajor) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteMajor) TestMajor() {
	target, err := newMajor(testdata.RedisURI)
	this.Nil(err)
	this.NotNil(target)
	target.stop()
	_, err = newMajor("")
	this.NotNil(err)
	_, err = newMajor(testdata.RedisURIInvalid)
	this.NotNil(err)
}

func (this *SuiteMajor) TestSubmit() {
	target, _ := newMajor(testdata.RedisURI)
	submit := target.Submit()
	this.NotNil(submit)
	result := submit.Ping(context.Background())
	_, err := submit.Exec(context.Background())
	this.Nil(err)
	ping, err := result.Result()
	this.Nil(err)
	this.Equal("PONG", ping)
	target.stop()
	this.Nil(target.Submit())
}

func (this *SuiteMajor) TestClient() {
	target, _ := newMajor(testdata.RedisURI)
	client := target.Client()
	this.NotNil(client)
	ping, err := client.Ping(context.Background()).Result()
	this.Nil(err)
	this.Equal("PONG", ping)
	target.stop()
}

func (this *SuiteMajor) TestSwitchDB() {
	target, _ := newMajor(testdata.RedisURI)
	this.Nil(target.SwitchDB(1))
	this.NotNil(target.SwitchDB(999999))
	target.stop()
	this.NotNil(target.SwitchDB(1))
}

func (this *SuiteMajor) TestDropDB() {
	target, _ := newMajor(testdata.RedisURI)
	target.DropDB()
	target.stop()
}

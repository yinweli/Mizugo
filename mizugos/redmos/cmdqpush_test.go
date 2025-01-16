package redmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdQPush(t *testing.T) {
	suite.Run(t, new(SuiteCmdQPush))
}

type SuiteCmdQPush struct {
	suite.Suite
	trials.Catalog
	major *Major
	minor *Minor
}

func (this *SuiteCmdQPush) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-qpush"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "qpush")
}

func (this *SuiteCmdQPush) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdQPush) TestQPush() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	name := "qpush"
	dataAll1 := &dataQPush{Field: "redis+mongo", Value: helps.RandStringDefault()}
	dataAll2 := &dataQPush{Field: "redis+mongo", Value: helps.RandStringDefault()}
	dataRedis := &dataQPush{Field: "redis", Value: helps.RandStringDefault()}
	dataMongo := &dataQPush{Field: "mongo", Value: helps.RandStringDefault()}

	target := &QPush[dataQPush]{MajorEnable: true, MinorEnable: true, Name: name, Data: dataAll1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	target = &QPush[dataQPush]{MajorEnable: true, MinorEnable: true, Name: name, Data: dataAll2}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompareList[dataQPush](this.major.Client(), name, []*dataQPush{dataAll1, dataAll2}))
	assert.True(this.T(), trials.MongoCompareList[dataQPush](this.minor.Database(), name, QPushTime, 1, []*dataQPush{dataAll1, dataAll2}))
	this.major.DropDB() // 如果不把佇列清空會影響之後的測試
	this.minor.DropDB()

	target = &QPush[dataQPush]{MajorEnable: true, MinorEnable: false, Name: name, Data: dataRedis}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompareList[dataQPush](this.major.Client(), name, []*dataQPush{dataRedis}))
	assert.False(this.T(), trials.MongoCompareList[dataQPush](this.minor.Database(), name, QPushTime, 1, []*dataQPush{dataRedis}))
	this.major.DropDB() // 如果不把佇列清空會影響之後的測試
	this.minor.DropDB()

	target = &QPush[dataQPush]{MajorEnable: false, MinorEnable: true, Name: name, Data: dataMongo}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.False(this.T(), trials.RedisCompareList[dataQPush](this.major.Client(), name, []*dataQPush{dataMongo}))
	assert.True(this.T(), trials.MongoCompareList[dataQPush](this.minor.Database(), name, QPushTime, 1, []*dataQPush{dataMongo}))
	this.major.DropDB() // 如果不把佇列清空會影響之後的測試
	this.minor.DropDB()

	target = &QPush[dataQPush]{MajorEnable: true, MinorEnable: true, Name: "", Data: dataAll1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &QPush[dataQPush]{MajorEnable: true, MinorEnable: true, Name: name, Data: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

type dataQPush struct {
	Field string `bson:"field"`
	Value string `bson:"value"`
}

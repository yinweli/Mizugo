package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdExist(t *testing.T) {
	suite.Run(t, new(SuiteCmdExist))
}

type SuiteCmdExist struct {
	suite.Suite
	trials.Catalog
	meta  metaExist
	major *Major
	minor *Minor
}

func (this *SuiteCmdExist) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdexist"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdexist")
}

func (this *SuiteCmdExist) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdExist) TestExist() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	key := []string{"0001", "0002", "0003"}
	this.major.Client().Set(context.Background(), this.meta.MajorKey(key[0]), "value0", 0)
	this.major.Client().Set(context.Background(), this.meta.MajorKey(key[1]), "value1", 0)

	target := &Exist{Meta: &this.meta, Key: key}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	assert.Equal(this.T(), 2, target.Count)

	target = &Exist{Meta: nil, Key: key}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Exist{Meta: &this.meta, Key: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Exist{Meta: &this.meta, Key: []string{}}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

type metaExist struct {
}

func (this *metaExist) MajorKey(key any) string {
	return fmt.Sprintf("cmdexist:%v", key)
}

func (this *metaExist) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaExist) MinorTable() string {
	return "cmdexist"
}

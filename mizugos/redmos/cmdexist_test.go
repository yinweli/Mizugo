package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestCmdExist(t *testing.T) {
	suite.Run(t, new(SuiteCmdExist))
}

type SuiteCmdExist struct {
	suite.Suite
	trials.Catalog
	meta  testMetaExist
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
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.Equal(2, target.Count)

	target = &Exist{Meta: nil, Key: key}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &Exist{Meta: &this.meta, Key: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &Exist{Meta: &this.meta, Key: []string{}}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())
}

type testMetaExist struct {
}

func (this *testMetaExist) MajorKey(key any) string {
	return fmt.Sprintf("cmdexist:%v", key)
}

func (this *testMetaExist) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaExist) MinorTable() string {
	return "cmdexist"
}

package redmos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugo/ctxs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdExist(t *testing.T) {
	suite.Run(t, new(SuiteCmdExist))
}

type SuiteCmdExist struct {
	suite.Suite
	testdata.Env
	meta  metaExist
	major *Major
	minor *Minor
}

func (this *SuiteCmdExist) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-cmdexist")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdexist")
}

func (this *SuiteCmdExist) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdExist) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteCmdExist) TestExist() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	key := []string{"0001", "0002", "0003"}
	this.major.Client().Set(ctxs.Get().Ctx(), this.meta.MajorKey(key[0]), "value0", 0)
	this.major.Client().Set(ctxs.Get().Ctx(), this.meta.MajorKey(key[1]), "value1", 0)

	target := &Exist{Meta: &this.meta, Key: key}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), target.Complete())
	assert.Equal(this.T(), 2, target.Count)

	target = &Exist{Meta: nil, Key: key}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Exist{Meta: &this.meta, Key: nil}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Exist{Meta: &this.meta, Key: []string{}}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Exist{Meta: nil, Key: key}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Complete())
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

func (this *metaExist) MinorField() string {
	return "field"
}

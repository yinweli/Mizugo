package redmos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMxExist(t *testing.T) {
	suite.Run(t, new(SuiteMxExist))
}

type SuiteMxExist struct {
	suite.Suite
	testdata.Env
	meta  metaMxExist
	major *Major
	minor *Minor
}

func (this *SuiteMxExist) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mxexist")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "mxexist")
}

func (this *SuiteMxExist) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteMxExist) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxExist) TestExist() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	key := []string{"mxexist-0001", "mxexist-0002", "mxexist-0003"}
	this.major.Client().Set(ctxs.Get().Ctx(), this.meta.MajorKey(key[0]), "value0", 0)
	this.major.Client().Set(ctxs.Get().Ctx(), this.meta.MajorKey(key[1]), "value1", 0)

	exist := &Exist{Meta: &this.meta, Key: key}
	exist.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), exist.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), exist.Complete())
	assert.Equal(this.T(), 2, exist.Count)

	exist = &Exist{Meta: nil, Key: key}
	exist.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), exist.Prepare())

	exist = &Exist{Meta: &this.meta, Key: nil}
	exist.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), exist.Prepare())

	exist = &Exist{Meta: &this.meta, Key: []string{}}
	exist.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), exist.Prepare())
}

type metaMxExist struct {
}

func (this *metaMxExist) Enable() (major, minor bool) {
	return true, false
}

func (this *metaMxExist) MajorKey(key any) string {
	return fmt.Sprintf("mxexist:%v", key)
}

func (this *metaMxExist) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaMxExist) MinorTable() string {
	return "mxexist_table"
}

func (this *metaMxExist) MinorField() string {
	return "mxexist_key"
}

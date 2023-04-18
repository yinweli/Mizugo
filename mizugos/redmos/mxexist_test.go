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
	key   []string
	major *Major
	minor *Minor
}

func (this *SuiteMxExist) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mxexist")
	this.key = []string{"mxexist-0001", "mxexist-0002", "mxexist-0003"}
	this.major, _ = newMajor(ctxs.Get().Ctx(), testdata.RedisURI, true)
	this.minor, _ = newMinor(ctxs.Get().Ctx(), testdata.MongoURI, this.meta.MinorTable()) // 這裡偷懶把表格名稱當資料庫名稱來用
}

func (this *SuiteMxExist) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	testdata.RedisClear(this.major.Client(), this.major.UsedKey())
	testdata.MongoClear(this.minor.Database())
	this.major.stop()
	this.minor.stop(ctxs.Get().Ctx())
}

func (this *SuiteMxExist) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxExist) TestGet() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	exist := &Exist{Meta: &this.meta, Key: this.key}
	exist.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)

	this.major.Client().Set(ctxs.Get().Ctx(), this.meta.MajorKey(this.key[0]), "value0", 0)
	this.major.Client().Set(ctxs.Get().Ctx(), this.meta.MajorKey(this.key[1]), "value1", 0)

	assert.Nil(this.T(), exist.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), exist.Complete())
	assert.Equal(this.T(), 2, exist.Count)

	exist.Meta = nil
	exist.Key = this.key
	assert.NotNil(this.T(), exist.Prepare())

	exist.Meta = &this.meta
	exist.Key = nil
	assert.NotNil(this.T(), exist.Prepare())

	exist.Meta = &this.meta
	exist.Key = []string{}
	assert.NotNil(this.T(), exist.Prepare())
}

type metaMxExist struct {
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

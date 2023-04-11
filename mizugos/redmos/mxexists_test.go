package redmos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMxExists(t *testing.T) {
	suite.Run(t, new(SuiteMxExists))
}

type SuiteMxExists struct {
	suite.Suite
	testdata.Env
	meta  metaMxExists
	key   []string
	major *Major
	minor *Minor
}

func (this *SuiteMxExists) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mxexists")
	this.key = []string{"mxexists-0001", "mxexists-0002", "mxexists-0003"}
	this.major, _ = newMajor(ctxs.Get().Ctx(), testdata.RedisURI, true)
	this.minor, _ = newMinor(ctxs.Get().Ctx(), testdata.MongoURI, this.meta.MinorTable()) // 這裡偷懶把表格名稱當資料庫名稱來用
}

func (this *SuiteMxExists) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	testdata.RedisClear(ctxs.Get().Ctx(), this.major.Client(), this.major.UsedKey())
	testdata.MongoClear(ctxs.Get().Ctx(), this.minor.Database())
	this.major.stop()
	this.minor.stop(ctxs.Get().Ctx())
}

func (this *SuiteMxExists) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxExists) TestGet() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	exists := &Exists{Meta: &this.meta, Key: this.key}
	exists.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)

	this.major.Client().Set(ctxs.Get().Ctx(), this.meta.MajorKey(this.key[0]), "value0", 0)
	this.major.Client().Set(ctxs.Get().Ctx(), this.meta.MajorKey(this.key[1]), "value1", 0)

	assert.Nil(this.T(), exists.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), exists.Complete())
	assert.Equal(this.T(), 2, exists.Count)

	exists.Meta = nil
	exists.Key = this.key
	assert.NotNil(this.T(), exists.Prepare())

	exists.Meta = &this.meta
	exists.Key = nil
	assert.NotNil(this.T(), exists.Prepare())

	exists.Meta = &this.meta
	exists.Key = []string{}
	assert.NotNil(this.T(), exists.Prepare())
}

type metaMxExists struct {
}

func (this *metaMxExists) MajorKey(key any) string {
	return fmt.Sprintf("mxexists:%v", key)
}

func (this *metaMxExists) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaMxExists) MinorTable() string {
	return "mxexists_table"
}

func (this *metaMxExists) MinorField() string {
	return "mxexists_key"
}

package redmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestCmdIndexComplex(t *testing.T) {
	suite.Run(t, new(SuiteCmdIndexComplex))
}

type SuiteCmdIndexComplex struct {
	suite.Suite
	trials.Catalog
	major *Major
	minor *Minor
}

func (this *SuiteCmdIndexComplex) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdindexcomplex"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdindexcomplex")
}

func (this *SuiteCmdIndexComplex) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdIndexComplex) TestIndexComplex() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	keys := []SortField{{Field: "score", Order: -1}, {Field: "time", Order: 1}}

	target := &IndexComplex{Name: "rank_index", Table: "cmdindexcomplex", Key: keys}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())

	target = &IndexComplex{Name: "rank_index", Table: "cmdindexcomplex", Key: keys}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())

	target = &IndexComplex{Name: "", Table: "cmdindexcomplex", Key: keys}
	this.NotNil(target.Prepare())

	target = &IndexComplex{Name: "rank_index", Table: "", Key: keys}
	this.NotNil(target.Prepare())

	target = &IndexComplex{Name: "rank_index", Table: "cmdindexcomplex", Key: []SortField{}}
	this.NotNil(target.Prepare())

	target = &IndexComplex{Name: "rank_index", Table: "cmdindexcomplex", Key: []SortField{{Field: "", Order: -1}}}
	this.NotNil(target.Prepare())

	target = &IndexComplex{Name: "rank_index", Table: "cmdindexcomplex", Key: []SortField{{Field: "score", Order: 0}}}
	this.NotNil(target.Prepare())
}

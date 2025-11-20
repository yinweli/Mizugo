package helps

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestBasen(t *testing.T) {
	suite.Run(t, new(SuiteBasen))
}

type SuiteBasen struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteBasen) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-basen"))
}

func (this *SuiteBasen) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteBasen) TestBase58() {
	testcase := map[uint64]string{
		0:                  "0",
		1:                  "1",
		10:                 "a",
		58:                 "10",
		1000:               "he",
		3364:               "100",
		100000:             "vJ8",
		195112:             "1000",
		10000000:           "TeDN",
		1000000000:         "1wmf9i",
		100000000000:       "2CkE832",
		10000000000000:     "4wFBKTds",
		1000000000000000:   "7NUhDEjaQ",
		100000000000000000: "dsSaqW54EL",
	}

	for k, v := range testcase {
		this.Equal(v, ToBase58(k))
		target, err := FromBase58(v)
		this.Nil(err)
		this.Equal(k, target)
	} // for

	_, err := FromBase58("/")
	this.NotNil(err)

	this.True(LessBase58("1", "2"))
	this.False(LessBase58("2", "1"))

	for i := uint64(0); i < testdata.TestCount; i++ {
		for _, c := range []byte{'o', 'O', 'l', 'I'} {
			this.NotContains(ToBase58(i), string([]byte{c}))
		} // for
	} // for
}

func (this *SuiteBasen) TestBase80() {
	testcase := map[uint64]string{
		0:                  "0",
		1:                  "1",
		10:                 "a",
		58:                 "W",
		1000:               "cE",
		3364:               "G4",
		100000:             "fO0",
		195112:             "uC@",
		10000000:           "jGE0",
		1000000000:         "oxa00",
		100000000000:       "uFwE00",
		10000000000000:     "CbYO000",
		1000000000000000:   "LST.E000",
		100000000000000000: "XMtWa0000",
	}

	for k, v := range testcase {
		this.Equal(v, ToBase80(k))
		target, err := FromBase80(v)
		this.Nil(err)
		this.Equal(k, target)
	} // for

	_, err := FromBase80("/")
	this.NotNil(err)

	this.True(LessBase80("1", "2"))
	this.False(LessBase80("2", "1"))
}

func (this *SuiteBasen) TestBaseN() {
	model := "0123456789"
	rank := RankBaseN(model)
	testcase := map[uint64]string{
		0:         "0",
		1:         "1",
		22:        "22",
		333:       "333",
		4444:      "4444",
		55555:     "55555",
		666666:    "666666",
		7777777:   "7777777",
		88888888:  "88888888",
		999999999: "999999999",
	}

	for k, v := range testcase {
		this.Equal(v, ToBaseN(k, model))
		target, err := FromBaseN(v, model, rank)
		this.Nil(err)
		this.Equal(k, target)
	} // for

	_, err := FromBaseN("", model, rank)
	this.NotNil(err)
	_, err = FromBaseN("/", model, rank)
	this.NotNil(err)
	_, err = FromBaseN("18446744073709551616", model, rank)
	this.NotNil(err)

	this.True(LessBaseN("1", "2", model, rank))
	this.False(LessBaseN("2", "1", model, rank))
	this.True(LessBaseN("1", "10", model, rank))
	this.False(LessBaseN("10", "1", model, rank))
	this.False(LessBaseN("1", "1", model, rank))
	this.True(LessBaseN("01", "2", model, rank))
	this.False(LessBaseN("01", "1", model, rank))
	this.True(LessBaseN("1a", "1b", model, rank))
	this.False(LessBaseN("1b", "1a", model, rank))

	this.Panics(func() {
		_ = RankBaseN("0")
	})
	this.Panics(func() {
		_ = RankBaseN("00123456789")
	})
}

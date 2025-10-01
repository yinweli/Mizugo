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
		result, err := FromBase58(v)
		this.Nil(err)
		this.Equal(k, result)
	} // for

	for i := uint64(0); i < testdata.TestCount; i++ {
		this.NotContains(ToBase58(i), "oOlI")
	} // for

	_, err := FromBase58("{}")
	this.NotNil(err)
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
		result, err := FromBase80(v)
		this.Nil(err)
		this.Equal(k, result)
	} // for

	_, err := FromBase80("{}")
	this.NotNil(err)
}

func (this *SuiteBasen) TestBaseN() {
	model := "0123456789"
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
		this.Equal(v, ToBaseN(model, k))
		result, err := FromBaseN(model, v)
		this.Nil(err)
		this.Equal(k, result)
	} // for

	_, err := FromBaseN("", "{}")
	this.NotNil(err)
	_, err = FromBaseN(model, "")
	this.NotNil(err)
	_, err = FromBaseN(model, "{}")
	this.NotNil(err)
}

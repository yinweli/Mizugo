package helps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestBasen(t *testing.T) {
	suite.Run(t, new(SuiteBasen))
}

type SuiteBasen struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteBasen) SetupSuite() {
	this.Env = testdata.EnvSetup("test-helps-basen")
}

func (this *SuiteBasen) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteBasen) TearDownTest() {
	testdata.Leak(this.T(), true)
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
		assert.Equal(this.T(), v, ToBase58(k))
		result, err := FromBase58(v)
		assert.Nil(this.T(), err)
		assert.Equal(this.T(), k, result)
	} // for

	for i := uint64(0); i < testdata.TestCount; i++ {
		assert.NotContains(this.T(), ToBase58(i), "oOlI")
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
		assert.Equal(this.T(), v, ToBase80(k))
		result, err := FromBase80(v)
		assert.Nil(this.T(), err)
		assert.Equal(this.T(), k, result)
	} // for
}

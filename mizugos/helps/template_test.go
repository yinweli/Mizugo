package helps

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestTemplate(t *testing.T) {
	suite.Run(t, new(SuiteTemplate))
}

type SuiteTemplate struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteTemplate) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-template"))
}

func (this *SuiteTemplate) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteTemplate) TestWriteTmpl() {
	path := filepath.Join("test", "test.tmpl")
	blueprintReal := "{{$.Value}}"
	blueprintFake := "{{{$.Value}}"
	assert.Nil(this.T(), WriteTemplate(path, blueprintReal, map[string]string{"Value": "Value"}))
	assert.True(this.T(), trials.FileCompare(path, []byte("Value")))
	assert.NotNil(this.T(), WriteTemplate(path, blueprintFake, nil))
	assert.NotNil(this.T(), WriteTemplate(path, blueprintReal, "nothing!"))
}

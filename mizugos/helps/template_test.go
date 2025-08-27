package helps

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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
	blueprint1 := "{{$.Value}}"
	blueprint2 := "{{{$.Value}}"
	builder := &strings.Builder{}
	assert.Nil(this.T(), TemplateBuild(blueprint1, builder, map[string]string{"Value": "Value"}))
	assert.Equal(this.T(), "Value", builder.String())
	assert.NotNil(this.T(), TemplateBuild(blueprint1, builder, "nothing!"))
	assert.NotNil(this.T(), TemplateBuild(blueprint2, builder, nil))
}

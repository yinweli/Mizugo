package helps

import (
	"strings"
	"testing"

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

func (this *SuiteTemplate) TestTemplateBuild() {
	builder := &strings.Builder{}
	this.Nil(TemplateBuild("{{$.Value}}", builder, map[string]string{"Value": "Value"}))
	this.Equal("Value", builder.String())
	this.NotNil(TemplateBuild("{{$.Value}}", builder, "nothing"))
	this.NotNil(TemplateBuild("{{{$.Value}}", builder, "syntax error"))
}

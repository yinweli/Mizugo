package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestConfigmgr(t *testing.T) {
	suite.Run(t, new(SuiteConfigmgr))
}

type SuiteConfigmgr struct {
	suite.Suite
	testdata.TestEnv
	content     string
	fileValid   string
	fileInvalid string
	kInt        string
	kInt32      string
	kInt64      string
	kString     string
	kList       string
	kMap        string
	kObject     string
	vint        int
	vint32      int
	vint64      int
	vstring     string
	vlist       []string
	vmap        map[string]interface{}
	vobject     configTester
}

func (this *SuiteConfigmgr) SetupSuite() {
	this.Change("test-configs-configmgr")
	this.content = `
Int: 1
Int32: 32
Int64: 64
String: string
List: [ "s1", "s2", "s3" ]
Map:
  Int: 1
  Int32: 32
  Int64: 64
  String: string
Object:
  Int: 1
  Int32: 32
  Int64: 64
  String: string
  List: [ "s1", "s2", "s3" ]
  Map:
    Int: 1
    Int32: 32
    Int64: 64
    String: string
`
	this.fileValid = "valid.yaml"
	this.fileInvalid = "invalid.yaml"
	this.kInt = "Int"
	this.kInt32 = "Int32"
	this.kInt64 = "Int64"
	this.kString = "String"
	this.kList = "List"
	this.kMap = "Map"
	this.kObject = "Object"
	this.vint = 1
	this.vint32 = 32
	this.vint64 = 64
	this.vstring = "string"
	this.vlist = []string{"s1", "s2", "s3"}
	this.vmap = map[string]interface{}{
		"Int":    this.vint,
		"Int32":  this.vint32,
		"Int64":  this.vint64,
		"String": this.vstring,
	}
	this.vobject = configTester{
		Int:    this.vint,
		Int32:  int32(this.vint32),
		Int64:  int64(this.vint64),
		String: this.vstring,
		List:   this.vlist,
		Map:    this.vmap,
	}

	assert.Nil(this.T(), os.WriteFile(this.fileValid, []byte(this.content), os.ModePerm))
	assert.Nil(this.T(), os.WriteFile(this.fileInvalid, []byte("fake"), os.ModePerm))
}

func (this *SuiteConfigmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteConfigmgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteConfigmgr) TestNewConfigmgr() {
	assert.NotNil(this.T(), NewConfigmgr())
}

func (this *SuiteConfigmgr) TestReadFile() {
	target := NewConfigmgr()
	assert.Nil(this.T(), target.ReadFile(this.fileValid))
	assert.NotNil(this.T(), target.ReadFile(this.fileInvalid))
	assert.NotNil(this.T(), target.ReadFile("!?"))
}

func (this *SuiteConfigmgr) TestReadString() {
	target := NewConfigmgr()
	assert.Nil(this.T(), target.ReadString(this.content))
	assert.NotNil(this.T(), target.ReadString("!?"))
}

func (this *SuiteConfigmgr) TestGetInt() {
	target := NewConfigmgr()
	assert.Nil(this.T(), target.ReadFile(this.fileValid))

	result, err := target.GetInt(this.kInt)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), this.vint, result)

	result, err = target.GetInt(this.kInt32)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), this.vint32, result)

	result, err = target.GetInt(this.kInt64)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), this.vint64, result)

	_, err = target.GetInt("!?")
	assert.NotNil(this.T(), err)

	_, err = target.GetInt(this.kString)
	assert.NotNil(this.T(), err)
}

func (this *SuiteConfigmgr) TestGetString() {
	target := NewConfigmgr()
	assert.Nil(this.T(), target.ReadFile(this.fileValid))

	result, err := target.GetString(this.kString)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), this.vstring, result)

	_, err = target.GetString("!?")
	assert.NotNil(this.T(), err)

	_, err = target.GetString(this.kInt)
	assert.NotNil(this.T(), err)
}

func (this *SuiteConfigmgr) TestGetObject() {
	target := NewConfigmgr()
	assert.Nil(this.T(), target.ReadFile(this.fileValid))

	resultList := []string{}
	assert.Nil(this.T(), target.GetObject(this.kList, &resultList))
	assert.Equal(this.T(), this.vlist, resultList)

	resultMap := map[string]interface{}{}
	assert.Nil(this.T(), target.GetObject(this.kMap, &resultMap))
	assert.Equal(this.T(), this.vmap, resultMap)

	resultObject := configTester{}
	assert.Nil(this.T(), target.GetObject(this.kObject, &resultObject))
	assert.Equal(this.T(), this.vobject, resultObject)

	assert.NotNil(this.T(), target.GetObject("!?", &resultList))
	assert.NotNil(this.T(), target.GetObject(this.kMap, &resultList))
}

type configTester struct {
	Int    int
	Int32  int32
	Int64  int64
	String string
	List   []string
	Map    map[string]interface{}
}

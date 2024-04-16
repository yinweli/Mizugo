package loggers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestZap(t *testing.T) {
	suite.Run(t, new(SuiteZap))
}

type SuiteZap struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteZap) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-loggers-zap"))
}

func (this *SuiteZap) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteZap) TestZapLogger() {
	target := &ZapLogger{
		Name:       "zapLogger",
		Path:       "zapLogger",
		Json:       true,
		Console:    true,
		File:       true,
		Level:      LevelDebug,
		TimeLayout: "2006-01-02 15:04:05.000",
	}
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.Get())
	assert.Equal(this.T(), zapcore.DebugLevel, target.zapLevel(LevelDebug))
	assert.Equal(this.T(), zapcore.InfoLevel, target.zapLevel(LevelInfo))
	assert.Equal(this.T(), zapcore.WarnLevel, target.zapLevel(LevelWarn))
	assert.Equal(this.T(), zapcore.ErrorLevel, target.zapLevel(LevelError))
	assert.Equal(this.T(), zapcore.InvalidLevel, target.zapLevel("!?"))
	target.Finalize()

	target = &ZapLogger{
		Name:       "zapLogger",
		Path:       "zapLogger",
		Json:       false,
		Console:    false,
		File:       false,
		Level:      LevelDebug,
		TimeLayout: "2006-01-02 15:04:05.000",
		TimeZone:   "Asia/Taipei",
	}
	assert.Nil(this.T(), target.Initialize())
	target.Finalize()

	target = &ZapLogger{
		TimeZone: testdata.Unknown,
	}
	assert.NotNil(this.T(), target.Initialize())
}

func (this *SuiteZap) TestZapRetain() {
	logger := &ZapLogger{
		Name:       "zapStream",
		Path:       "zapStream",
		Json:       true,
		Console:    true,
		Level:      LevelDebug,
		TimeLayout: "2006-01-02 15:04:05.000",
	}
	_ = logger.Initialize()

	target := logger.Get()
	assert.Equal(this.T(), target, target.Clear())
	assert.Equal(this.T(), target, target.Debug("").End().Flush())
	assert.NotNil(this.T(), target.Debug(""))
	assert.NotNil(this.T(), target.Info(""))
	assert.NotNil(this.T(), target.Warn(""))
	assert.NotNil(this.T(), target.Error(""))
}

func (this *SuiteZap) TestZapStream() {
	logger := &ZapLogger{
		Name:       "zapStream",
		Path:       "zapStream",
		Json:       true,
		Console:    true,
		Level:      LevelDebug,
		TimeLayout: "2006-01-02 15:04:05.000",
	}
	_ = logger.Initialize()
	retain := logger.Get()

	target := retain.Debug("log")
	assert.Equal(this.T(), target, target.Message("message"))
	assert.Equal(this.T(), target, target.KV("key", "value"))
	assert.Equal(this.T(), target, target.Caller(0))
	assert.Equal(this.T(), target, target.Error(fmt.Errorf("error")))
	assert.Equal(this.T(), retain, target.End())
	target.EndFlush()

	logger.Finalize()
}

func (this *SuiteZap) TestZapStreamKV() {
	logger := &ZapLogger{
		Name:       "zapStream",
		Path:       "zapStream",
		Json:       true,
		Console:    true,
		Level:      LevelDebug,
		TimeLayout: "2006-01-02 15:04:05.000",
		TimeZone:   "Asia/Taipei",
		Jsonify:    true,
	}
	_ = logger.Initialize()
	retain := logger.Get()

	i8 := int8(0)
	ui8 := uint8(0)
	i8s := []int8{i8}
	ui8s := []uint8{ui8}
	i16 := int16(0)
	ui16 := uint16(0)
	i16s := []int16{i16}
	ui16s := []uint16{ui16}
	i32 := int32(0)
	ui32 := uint32(0)
	i32s := []int32{i32}
	ui32s := []uint32{ui32}
	i64 := int64(0)
	ui64 := uint64(0)
	i64s := []int64{i64}
	ui64s := []uint64{ui64}
	i := int(0)
	ui := uint(0)
	is := []int{i}
	uis := []uint{ui}
	f32 := float32(0)
	f32s := []float32{f32}
	f64 := float64(0)
	f64s := []float64{f64}
	c64 := complex64(0)
	c64s := []complex64{c64}
	c128 := complex128(0)
	c128s := []complex128{c128}
	s := "value"
	ss := []string{s}
	b := false
	bs := []bool{false}
	o := struct {
		Name  string
		Value int
	}{Name: "name", Value: 1}
	m := map[string]int{
		"0": 0,
		"1": 1,
	}

	target := retain.Debug("log")
	assert.Equal(this.T(), target, target.KV("int8", i8))
	assert.Equal(this.T(), target, target.KV("uint8", ui8))
	assert.Equal(this.T(), target, target.KV("*int8", &i8))
	assert.Equal(this.T(), target, target.KV("*uint8", &ui8))
	assert.Equal(this.T(), target, target.KV("int8[]", i8s))
	assert.Equal(this.T(), target, target.KV("uint8[]", ui8s))
	assert.Equal(this.T(), target, target.KV("int16", i16))
	assert.Equal(this.T(), target, target.KV("uint16", ui16))
	assert.Equal(this.T(), target, target.KV("*int16", &i16))
	assert.Equal(this.T(), target, target.KV("*uint16", &ui16))
	assert.Equal(this.T(), target, target.KV("[]int16", i16s))
	assert.Equal(this.T(), target, target.KV("[]uint16", ui16s))
	assert.Equal(this.T(), target, target.KV("int32", i32))
	assert.Equal(this.T(), target, target.KV("uint32", ui32))
	assert.Equal(this.T(), target, target.KV("*int32", &i32))
	assert.Equal(this.T(), target, target.KV("*uint32", &ui32))
	assert.Equal(this.T(), target, target.KV("[]int32", i32s))
	assert.Equal(this.T(), target, target.KV("[]uint32", ui32s))
	assert.Equal(this.T(), target, target.KV("int64", i64))
	assert.Equal(this.T(), target, target.KV("uint64", ui64))
	assert.Equal(this.T(), target, target.KV("*int64", &i64))
	assert.Equal(this.T(), target, target.KV("*uint64", &ui64))
	assert.Equal(this.T(), target, target.KV("[]int64", i64s))
	assert.Equal(this.T(), target, target.KV("[]uint64", ui64s))
	assert.Equal(this.T(), target, target.KV("int", i))
	assert.Equal(this.T(), target, target.KV("uint", ui))
	assert.Equal(this.T(), target, target.KV("*int", &i))
	assert.Equal(this.T(), target, target.KV("*uint", &ui))
	assert.Equal(this.T(), target, target.KV("[]int", is))
	assert.Equal(this.T(), target, target.KV("[]uint", uis))
	assert.Equal(this.T(), target, target.KV("float32", f32))
	assert.Equal(this.T(), target, target.KV("*float32", &f32))
	assert.Equal(this.T(), target, target.KV("[]float32", f32s))
	assert.Equal(this.T(), target, target.KV("float64", f64))
	assert.Equal(this.T(), target, target.KV("*float64", &f64))
	assert.Equal(this.T(), target, target.KV("[]float64", f64s))
	assert.Equal(this.T(), target, target.KV("complex64", c64))
	assert.Equal(this.T(), target, target.KV("*complex64", &c64))
	assert.Equal(this.T(), target, target.KV("[]complex64", c64s))
	assert.Equal(this.T(), target, target.KV("complex128", c128))
	assert.Equal(this.T(), target, target.KV("*complex128", &c128))
	assert.Equal(this.T(), target, target.KV("[]complex128", c128s))
	assert.Equal(this.T(), target, target.KV("string", s))
	assert.Equal(this.T(), target, target.KV("*string", &s))
	assert.Equal(this.T(), target, target.KV("[]string", ss))
	assert.Equal(this.T(), target, target.KV("bool", b))
	assert.Equal(this.T(), target, target.KV("*bool", &b))
	assert.Equal(this.T(), target, target.KV("[]bool", bs))
	assert.Equal(this.T(), target, target.KV("object", o))
	assert.Equal(this.T(), target, target.KV("*object", &o))
	assert.Equal(this.T(), target, target.KV("map", m))
	target.EndFlush()

	logger.Finalize()
}

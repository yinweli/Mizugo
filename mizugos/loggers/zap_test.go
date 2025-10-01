package loggers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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
	this.Nil(target.Initialize())
	this.NotNil(target.Get())
	this.Equal(zapcore.DebugLevel, target.zapLevel(LevelDebug))
	this.Equal(zapcore.InfoLevel, target.zapLevel(LevelInfo))
	this.Equal(zapcore.WarnLevel, target.zapLevel(LevelWarn))
	this.Equal(zapcore.ErrorLevel, target.zapLevel(LevelError))
	this.Equal(zapcore.InfoLevel, target.zapLevel("!?"))
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
	this.Nil(target.Initialize())
	target.Finalize()

	target = &ZapLogger{
		Name:     "zapLogger",
		Path:     "zapLogger",
		TimeZone: "Asia/Taipei",
	}
	this.Nil(target.Initialize())
	target.Finalize()

	target = &ZapLogger{
		TimeZone: testdata.Unknown,
	}
	this.NotNil(target.Initialize())
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
	this.Equal(target, target.Clear())
	this.Equal(target, target.Debug("").End().Flush())
	this.NotNil(target.Debug(""))
	this.NotNil(target.Info(""))
	this.NotNil(target.Warn(""))
	this.NotNil(target.Error(""))
	logger.Finalize()
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
	this.Equal(target, target.Message("message"))
	this.Equal(target, target.KV("key", "value"))
	this.Equal(target, target.Caller(0))
	this.Equal(target, target.Caller(0, true))
	this.Equal(target, target.Caller(0, false))
	this.Equal(target, target.Error(fmt.Errorf("error")))
	this.Equal(retain, target.End())
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

	retain := logger.Get()
	target := retain.Debug("log")
	this.Equal(target, target.KV("int8", i8))
	this.Equal(target, target.KV("uint8", ui8))
	this.Equal(target, target.KV("*int8", &i8))
	this.Equal(target, target.KV("*uint8", &ui8))
	this.Equal(target, target.KV("int8[]", i8s))
	this.Equal(target, target.KV("uint8[]", ui8s))
	this.Equal(target, target.KV("int16", i16))
	this.Equal(target, target.KV("uint16", ui16))
	this.Equal(target, target.KV("*int16", &i16))
	this.Equal(target, target.KV("*uint16", &ui16))
	this.Equal(target, target.KV("[]int16", i16s))
	this.Equal(target, target.KV("[]uint16", ui16s))
	this.Equal(target, target.KV("int32", i32))
	this.Equal(target, target.KV("uint32", ui32))
	this.Equal(target, target.KV("*int32", &i32))
	this.Equal(target, target.KV("*uint32", &ui32))
	this.Equal(target, target.KV("[]int32", i32s))
	this.Equal(target, target.KV("[]uint32", ui32s))
	this.Equal(target, target.KV("int64", i64))
	this.Equal(target, target.KV("uint64", ui64))
	this.Equal(target, target.KV("*int64", &i64))
	this.Equal(target, target.KV("*uint64", &ui64))
	this.Equal(target, target.KV("[]int64", i64s))
	this.Equal(target, target.KV("[]uint64", ui64s))
	this.Equal(target, target.KV("int", i))
	this.Equal(target, target.KV("uint", ui))
	this.Equal(target, target.KV("*int", &i))
	this.Equal(target, target.KV("*uint", &ui))
	this.Equal(target, target.KV("[]int", is))
	this.Equal(target, target.KV("[]uint", uis))
	this.Equal(target, target.KV("float32", f32))
	this.Equal(target, target.KV("*float32", &f32))
	this.Equal(target, target.KV("[]float32", f32s))
	this.Equal(target, target.KV("float64", f64))
	this.Equal(target, target.KV("*float64", &f64))
	this.Equal(target, target.KV("[]float64", f64s))
	this.Equal(target, target.KV("complex64", c64))
	this.Equal(target, target.KV("*complex64", &c64))
	this.Equal(target, target.KV("[]complex64", c64s))
	this.Equal(target, target.KV("complex128", c128))
	this.Equal(target, target.KV("*complex128", &c128))
	this.Equal(target, target.KV("[]complex128", c128s))
	this.Equal(target, target.KV("string", s))
	this.Equal(target, target.KV("*string", &s))
	this.Equal(target, target.KV("[]string", ss))
	this.Equal(target, target.KV("bool", b))
	this.Equal(target, target.KV("*bool", &b))
	this.Equal(target, target.KV("[]bool", bs))
	this.Equal(target, target.KV("object", o))
	this.Equal(target, target.KV("*object", &o))
	this.Equal(target, target.KV("map", m))
	target.EndFlush()
	logger.Finalize()
}

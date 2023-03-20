package logs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"

	"github.com/yinweli/Mizugo/testdata"
)

func TestZap(t *testing.T) {
	suite.Run(t, new(SuiteZap))
}

type SuiteZap struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteZap) SetupSuite() {
	this.Change("test-logs-zap")
}

func (this *SuiteZap) TearDownSuite() {
	this.Restore()
}

func (this *SuiteZap) TearDownTest() {
	this.GoLeak(this.T(), false) // 由於不清楚(或是沒辦法)優雅的關閉zap的執行緒, 所以只好把這裡的執行緒洩漏檢查關閉
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
	target.Finalize() // 初始化前執行, 這次應該不執行
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.Initialize()) // 故意啟動兩次, 這次應該失敗
	assert.NotNil(this.T(), target.Debug(""))
	assert.NotNil(this.T(), target.Info(""))
	assert.NotNil(this.T(), target.Warn(""))
	assert.NotNil(this.T(), target.Error(""))
	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行

	target = &ZapLogger{
		Name:       "zapLogger",
		Path:       "zapLogger",
		Json:       false,
		Console:    false,
		File:       false,
		Level:      LevelDebug,
		TimeLayout: "2006-01-02 15:04:05.000",
	}
	assert.Nil(this.T(), target.Initialize())
	target.Finalize()
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
	assert.Nil(this.T(), logger.Initialize())

	target := logger.Debug("log")
	assert.Equal(this.T(), target, target.Message("message"))
	assert.Equal(this.T(), target, target.Caller(0))
	assert.Equal(this.T(), target, target.Error(fmt.Errorf("error")))
	target.End()

	target = logger.Debug("log")
	assert.Equal(this.T(), target, target.Message("message"))
	target.EndError(fmt.Errorf("end error"))

	logger.Finalize()
}

func (this *SuiteZap) TestZapStreamKV() {
	logger := &ZapLogger{
		Name:       "zapStreamKV",
		Path:       "zapStreamKV",
		Json:       true,
		Console:    true,
		Level:      LevelDebug,
		TimeLayout: "2006-01-02 15:04:05.000",
	}
	assert.Nil(this.T(), logger.Initialize())

	key := "key"
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
	obj := struct {
		Name  string
		Value int
	}{Name: "name", Value: 1}

	target := logger.Debug("log")
	assert.Equal(this.T(), target, target.KV(key, i8))
	assert.Equal(this.T(), target, target.KV(key, ui8))
	assert.Equal(this.T(), target, target.KV(key, &i8))
	assert.Equal(this.T(), target, target.KV(key, &ui8))
	assert.Equal(this.T(), target, target.KV(key, i8s))
	assert.Equal(this.T(), target, target.KV(key, ui8s))
	assert.Equal(this.T(), target, target.KV(key, i16))
	assert.Equal(this.T(), target, target.KV(key, ui16))
	assert.Equal(this.T(), target, target.KV(key, &i16))
	assert.Equal(this.T(), target, target.KV(key, &ui16))
	assert.Equal(this.T(), target, target.KV(key, i16s))
	assert.Equal(this.T(), target, target.KV(key, ui16s))
	assert.Equal(this.T(), target, target.KV(key, i32))
	assert.Equal(this.T(), target, target.KV(key, ui32))
	assert.Equal(this.T(), target, target.KV(key, &i32))
	assert.Equal(this.T(), target, target.KV(key, &ui32))
	assert.Equal(this.T(), target, target.KV(key, i32s))
	assert.Equal(this.T(), target, target.KV(key, ui32s))
	assert.Equal(this.T(), target, target.KV(key, i64))
	assert.Equal(this.T(), target, target.KV(key, ui64))
	assert.Equal(this.T(), target, target.KV(key, &i64))
	assert.Equal(this.T(), target, target.KV(key, &ui64))
	assert.Equal(this.T(), target, target.KV(key, i64s))
	assert.Equal(this.T(), target, target.KV(key, ui64s))
	assert.Equal(this.T(), target, target.KV(key, i))
	assert.Equal(this.T(), target, target.KV(key, ui))
	assert.Equal(this.T(), target, target.KV(key, &i))
	assert.Equal(this.T(), target, target.KV(key, &ui))
	assert.Equal(this.T(), target, target.KV(key, is))
	assert.Equal(this.T(), target, target.KV(key, uis))
	assert.Equal(this.T(), target, target.KV(key, f32))
	assert.Equal(this.T(), target, target.KV(key, &f32))
	assert.Equal(this.T(), target, target.KV(key, f32s))
	assert.Equal(this.T(), target, target.KV(key, f64))
	assert.Equal(this.T(), target, target.KV(key, &f64))
	assert.Equal(this.T(), target, target.KV(key, f64s))
	assert.Equal(this.T(), target, target.KV(key, c64))
	assert.Equal(this.T(), target, target.KV(key, &c64))
	assert.Equal(this.T(), target, target.KV(key, c64s))
	assert.Equal(this.T(), target, target.KV(key, c128))
	assert.Equal(this.T(), target, target.KV(key, &c128))
	assert.Equal(this.T(), target, target.KV(key, c128s))
	assert.Equal(this.T(), target, target.KV(key, s))
	assert.Equal(this.T(), target, target.KV(key, &s))
	assert.Equal(this.T(), target, target.KV(key, ss))
	assert.Equal(this.T(), target, target.KV(key, b))
	assert.Equal(this.T(), target, target.KV(key, &b))
	assert.Equal(this.T(), target, target.KV(key, bs))
	assert.Equal(this.T(), target, target.KV(key, obj))
	target.End()

	logger.Finalize()
}

func (this *SuiteZap) TestZapLevel() {
	assert.Equal(this.T(), zapcore.DebugLevel, zapLevel(LevelDebug))
	assert.Equal(this.T(), zapcore.InfoLevel, zapLevel(LevelInfo))
	assert.Equal(this.T(), zapcore.WarnLevel, zapLevel(LevelWarn))
	assert.Equal(this.T(), zapcore.ErrorLevel, zapLevel(LevelError))
	assert.Equal(this.T(), zapcore.InvalidLevel, zapLevel("!?"))
}

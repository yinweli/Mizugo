package loggers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger 基於 uber/zap 的高效能日誌實作
//
// 使用方式:
//   - 建立 ZapLogger 並填寫公開欄位, 可從 YAML 設定檔載入
//   - 呼叫 Initialize 完成初始化
//   - 透過 Get 取得 Retain 來記錄日誌
//   - 程式結束時呼叫 Finalize 釋放資源
//
// 注意:
//   - ZapRetain 並非執行緒安全, 不可跨執行緒共用
//   - 每個執行緒應該建立並使用自己的 Retain
type ZapLogger struct {
	Name       string `yaml:"name"`       // 日誌檔案名稱
	Path       string `yaml:"path"`       // 日誌輸出路徑
	Json       bool   `yaml:"json"`       // 是否輸出為 JSON 格式(建議正式環境開啟)
	Console    bool   `yaml:"console"`    // 是否輸出到控制台
	File       bool   `yaml:"file"`       // 是否輸出到檔案
	Level      string `yaml:"level"`      // 日誌等級，僅記錄 >= 此等級的訊息; 可選 LevelDebug, LevelInfo, LevelWarn, LevelError
	TimeLayout string `yaml:"timeLayout"` // 時間格式(例如: 2006-01-02 15:04:05 或 ISO8601 / RFC3339 標準)
	TimeZone   string `yaml:"timeZone"`   // 時區字串(與 time.LoadLocation 相同, 預設 time.UTC)
	MaxSize    int    `yaml:"maxSize"`    // 日誌檔案最大容量(MB), 超過時會切割新檔案(預設 100MB)
	MaxTime    int    `yaml:"maxTime"`    // 日誌檔案保留天數, 超過時會刪除(預設不刪除)
	MaxBackups int    `yaml:"maxBackups"` // 保留的日誌檔案數量, 超過時刪除最舊檔案(預設不限制)
	Compress   bool   `yaml:"compress"`   // 是否壓縮舊日誌檔案(預設 false)
	Jsonify    bool   `yaml:"jsonify"`    // 記錄物件時是否以 JSON 格式序列化(預設 false)

	location *time.Location // 時區物件
	logger   *zap.Logger    // zap日誌物件

	// ISO8601標準: "2006-01-02T15:04:05.000Z0700"
	// RFC3339標準: "2006-01-02T15:04:05.000000Z07:00"
}

// Initialize 初始化處理
func (this *ZapLogger) Initialize() error {
	location, err := this.timeZone()

	if err != nil {
		return fmt.Errorf("zapLogger initialize: %w", err)
	} // if

	this.location = location
	this.logger = zap.New(
		zapcore.NewCore(
			this.encoder(),
			this.writeSyncer(),
			zap.NewAtomicLevelAt(this.zapLevel(this.Level)),
		),
		zap.AddCallerSkip(1))
	return nil
}

// Finalize 結束處理
func (this *ZapLogger) Finalize() {
	if this.logger != nil {
		_ = this.logger.Sync()
	} // if

	this.location = nil
	this.logger = nil
}

// Get 取得 Retain 儲存器
func (this *ZapLogger) Get() Retain {
	return &ZapRetain{
		location: this.location,
		logger:   this.logger,
		jsonify:  this.Jsonify,
	}
}

// timeZone 取得時區物件
func (this *ZapLogger) timeZone() (location *time.Location, err error) {
	if this.TimeZone == "" {
		return time.UTC, nil
	} // if

	if location, err = time.LoadLocation(this.TimeZone); err != nil {
		return nil, fmt.Errorf("zapLogger timeZone: %w", err)
	} // if

	return location, nil
}

// encoder 取得日誌編碼器
func (this *ZapLogger) encoder() zapcore.Encoder {
	layout := this.TimeLayout

	if layout == "" {
		layout = time.RFC3339Nano
	} // if

	encoder := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "label",
		CallerKey:      "line",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(layout),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	if this.Json {
		return zapcore.NewJSONEncoder(encoder)
	} else {
		return zapcore.NewConsoleEncoder(encoder)
	} // if
}

// writeSyncer 取得寫入同步器
func (this *ZapLogger) writeSyncer() zapcore.WriteSyncer {
	writeSyncer := []zapcore.WriteSyncer{}

	if this.Console {
		writeSyncer = append(writeSyncer, zapcore.AddSync(os.Stdout))
	} // if

	if this.File {
		writeSyncer = append(writeSyncer, zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(this.Path, this.Name),
			MaxSize:    this.MaxSize,
			MaxAge:     this.MaxTime,
			MaxBackups: this.MaxBackups,
			Compress:   this.Compress,
		}))
	} // if

	if len(writeSyncer) == 0 { // 如果都沒有設定, 使用 io.Discard 作為黑洞輸出, 避免錯誤
		writeSyncer = append(writeSyncer, zapcore.AddSync(io.Discard))
	} // if

	return zapcore.NewMultiWriteSyncer(writeSyncer...)
}

// zapLevel 日誌等級字串轉換為zap日誌等級
func (this *ZapLogger) zapLevel(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel

	case LevelInfo:
		return zapcore.InfoLevel

	case LevelWarn:
		return zapcore.WarnLevel

	case LevelError:
		return zapcore.ErrorLevel

	default:
		return zapcore.InfoLevel
	} // switch
}

// ZapRetain zap儲存器
type ZapRetain struct {
	location *time.Location // 時區物件
	logger   *zap.Logger    // 日誌物件
	jsonify  bool           // 物件記錄時是否以json字串記錄
	stream   []*ZapStream   // 記錄列表
}

// Clear 清空內部 Stream 列表
func (this *ZapRetain) Clear() Retain {
	this.stream = nil
	return this
}

// Flush 儲存並清空內部 Stream 列表
func (this *ZapRetain) Flush() Retain {
	if this.logger == nil {
		this.stream = nil
		return this
	} // if

	location := this.location

	if location == nil {
		location = time.UTC
	} // if

	for _, itor := range this.stream {
		if log := this.logger.Named(itor.label).Check(itor.level, itor.message); log != nil {
			log.Time = time.Now().In(location)
			log.Write(itor.field...)
		} // if
	} // for

	this.stream = nil
	return this
}

// Debug 建立除錯訊息的 Stream
func (this *ZapRetain) Debug(label string) Stream {
	return &ZapStream{
		retain:  this,
		level:   zapcore.DebugLevel,
		label:   label,
		jsonify: this.jsonify,
	}
}

// Info 建立一般訊息的 Stream
func (this *ZapRetain) Info(label string) Stream {
	return &ZapStream{
		retain:  this,
		level:   zapcore.InfoLevel,
		label:   label,
		jsonify: this.jsonify,
	}
}

// Warn 建立警告訊息的 Stream
func (this *ZapRetain) Warn(label string) Stream {
	return &ZapStream{
		retain:  this,
		level:   zapcore.WarnLevel,
		label:   label,
		jsonify: this.jsonify,
	}
}

// Error 建立錯誤訊息的 Stream
func (this *ZapRetain) Error(label string) Stream {
	return &ZapStream{
		retain:  this,
		level:   zapcore.ErrorLevel,
		label:   label,
		jsonify: this.jsonify,
	}
}

// ZapStream zap記錄器
type ZapStream struct {
	retain  *ZapRetain    // 儲存器
	level   zapcore.Level // 日誌等級
	label   string        // 日誌標籤
	jsonify bool          // 物件記錄時是否以json字串記錄
	message string        // 訊息字串
	field   []zap.Field   // 索引與數值列表
}

// Message 記錄文字訊息
func (this *ZapStream) Message(format string, a ...any) Stream {
	this.message = fmt.Sprintf(format, a...)
	return this
}

// KV 記錄鍵值訊息
func (this *ZapStream) KV(key string, value any) Stream {
	field := zap.Any(key, value)

	if this.jsonify && field.Type == zapcore.ReflectType {
		if bytes, err := json.Marshal(value); err == nil {
			field = zap.String(key, string(bytes))
		} else {
			field = zap.String(key, fmt.Sprintf("json_error:%v", err))
		} // if
	} // if

	this.field = append(this.field, field)
	return this
}

// Caller 記錄呼叫位置
func (this *ZapStream) Caller(skip int, simple ...bool) Stream {
	if pc, _, _, ok := runtime.Caller(skip + 1); ok { // 這裡把skip+1的原因是為了多跳過現在這層, 這樣外部使用時就可以指定0為呼叫起點, 比較直覺
		caller := filepath.Base(runtime.FuncForPC(pc).Name())

		if len(simple) > 0 && simple[0] {
			if last := strings.Index(caller, "."); last != -1 && last+1 < len(caller) {
				caller = strings.Trim(caller[last+1:], "()*")
			} // if
		} // if

		this.field = append(this.field, zap.String("caller", caller))
	} // if

	return this
}

// Error 記錄錯誤物件
func (this *ZapStream) Error(err error) Stream {
	this.field = append(this.field, zap.Error(err))
	return this
}

// End 結束記錄, 將記錄交回 Retain
func (this *ZapStream) End() Retain {
	this.retain.stream = append(this.retain.stream, this)
	return this.retain
}

// EndFlush 結束記錄, 將記錄交回 Retain 並立即儲存
func (this *ZapStream) EndFlush() {
	this.End().Flush()
}

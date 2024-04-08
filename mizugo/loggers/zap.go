package loggers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger zap日誌, uber實現的高效能日誌功能;
// 使用前必須填寫好 ZapLogger 中的公開成員, 可以選擇從yaml格式的配置檔案來填寫 ZapLogger 結構;
// 使用時要注意由於 ZapRetain 並非執行緒安全, 因此使用時不可以把 ZapRetain 儲存下來使用
type ZapLogger struct {
	// 關於時間布局字串
	// - 如果不需要時區資訊, 可以寫成"2006-01-02 15:04:05"
	// - 遵循ISO8601標準, 可以寫成"2006-01-02T15:04:05.000Z0700"
	// - 遵循RFC3339標準, 可以寫成"2006-01-02T15:04:05.000000Z07:00"

	Name       string `yaml:"name"`       // 日誌名稱, 會被用到日誌檔案名稱上
	Path       string `yaml:"path"`       // 日誌路徑, 指定日誌檔案的位置
	Json       bool   `yaml:"json"`       // 是否使用json格式日誌, 建議正式環境使用json格式日誌
	Console    bool   `yaml:"console"`    // 是否輸出到控制台
	File       bool   `yaml:"file"`       // 是否輸出到日誌檔案
	Level      string `yaml:"level"`      // 輸出日誌等級, 當記錄的日誌等級超過此值時才會儲存到檔案中; 有以下選擇: LevelDebug, LevelInfo, LevelWarn, LevelError
	TimeLayout string `yaml:"timeLayout"` // 時間布局字串
	TimeZone   string `yaml:"timeZone"`   // 時區字串, 採用與 time.LoadLocation 一樣的方式, 預設是UTC+0
	MaxSize    int    `yaml:"maxSize"`    // 日誌大小(MB), 當日誌檔案超過此大小時就會建立新檔案, 預設是100MB
	MaxTime    int    `yaml:"maxTime"`    // 日誌保留時間(日), 當日誌檔案儲存超過此時間時會被刪除, 預設不會刪除檔案
	MaxBackups int    `yaml:"maxBackups"` // 日誌保留數量, 當日誌檔案數量超過此數量時會刪除舊檔案, 預設不會刪除檔案
	Compress   bool   `yaml:"compress"`   // 是否壓縮日誌檔案, 預設不會壓縮
	Jsonify    bool   `yaml:"jsonify"`    // 物件記錄時是否以json字串記錄, 預設為關閉

	location *time.Location // 時區物件
	logger   *zap.Logger    // zap日誌物件
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

// Get 取得儲存器
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
		EncodeTime:     zapcore.TimeEncoderOfLayout(this.TimeLayout),
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
		return zapcore.InvalidLevel
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
	for _, itor := range this.stream {
		logger := this.logger.Named(itor.label)

		if logger == nil {
			continue
		} // if

		log := logger.Check(itor.level, itor.message)

		if log == nil {
			continue
		} // if

		log.Time = time.Now().In(this.location)
		log.Write(itor.field...)
	} // for

	this.stream = nil
	return this
}

// Debug 記錄除錯訊息, 用於記錄除錯訊息
func (this *ZapRetain) Debug(label string) Stream {
	return &ZapStream{
		retain:  this,
		level:   zapcore.DebugLevel,
		label:   label,
		jsonify: this.jsonify,
	}
}

// Info 記錄一般訊息, 用於記錄一般訊息
func (this *ZapRetain) Info(label string) Stream {
	return &ZapStream{
		retain:  this,
		level:   zapcore.InfoLevel,
		label:   label,
		jsonify: this.jsonify,
	}
}

// Warn 記錄警告訊息, 用於記錄邏輯錯誤
func (this *ZapRetain) Warn(label string) Stream {
	return &ZapStream{
		retain:  this,
		level:   zapcore.WarnLevel,
		label:   label,
		jsonify: this.jsonify,
	}
}

// Error 記錄錯誤訊息, 用於記錄嚴重錯誤
func (this *ZapRetain) Error(label string) Stream {
	return &ZapStream{
		retain:  this,
		level:   zapcore.ErrorLevel,
		label:   label,
		jsonify: this.jsonify,
	}
}

// ZapStream zap記錄
type ZapStream struct {
	retain  *ZapRetain    // 儲存器
	level   zapcore.Level // 日誌等級
	label   string        // 日誌標籤
	jsonify bool          // 物件記錄時是否以json字串記錄
	message string        // 訊息字串
	field   []zap.Field   // 索引與數值列表
}

// Message 記錄訊息
func (this *ZapStream) Message(format string, a ...any) Stream {
	this.message = fmt.Sprintf(format, a...)
	return this
}

// KV 記錄索引與數值
func (this *ZapStream) KV(key string, value any) Stream {
	field := zap.Any(key, value)

	if this.jsonify && field.Type == zapcore.ReflectType {
		bytes, _ := json.Marshal(value)
		field = zap.String(key, string(bytes))
	} // if

	this.field = append(this.field, field)
	return this
}

// Caller 記錄呼叫位置
func (this *ZapStream) Caller(skip int) Stream {
	if pc, _, _, ok := runtime.Caller(skip + 1); ok { // 這裡把skip+1的原因是為了多跳過現在這層, 這樣外部使用時就可以指定0為呼叫起點, 比較直覺
		this.field = append(this.field, zap.String("caller", filepath.Base(runtime.FuncForPC(pc).Name())))
	} // if

	return this
}

// Error 記錄錯誤
func (this *ZapStream) Error(err error) Stream {
	this.field = append(this.field, zap.Error(err))
	return this
}

// End 結束記錄
func (this *ZapStream) End() Retain {
	this.retain.stream = append(this.retain.stream, this)
	return this.retain
}

// EndFlush 結束記錄, 並把記錄加回到 Retain 中, 然後儲存記錄
func (this *ZapStream) EndFlush() {
	this.End().Flush()
}

package logs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

// zap日誌, uber實現的高效能日誌功能
// 使用前必須填寫好ZapLogger中的公開成員, 可以選擇從yaml格式的配置檔案來填寫ZapLogger結構

// ZapLogger zap日誌
type ZapLogger struct {
	Name       string `yaml:"name"`       // 日誌名稱, 會被用到日誌檔案名稱上
	Path       string `yaml:"path"`       // 日誌路徑, 指定日誌檔案的位置
	Json       bool   `yaml:"json"`       // 是否使用json格式日誌, 建議正式環境使用json格式日誌
	Console    bool   `yaml:"console"`    // 是否輸出到控制台
	File       bool   `yaml:"file"`       // 是否輸出到日誌檔案
	Level      Level  `yaml:"level"`      // 日誌等級
	MaxSize    int    `yaml:"maxSize"`    // 日誌大小(MB), 當日誌檔案超過此大小時就會建立新檔案, 預設是100MB
	MaxTime    int    `yaml:"maxTime"`    // 日誌保留時間(日), 當日誌檔案儲存超過此時間時會被刪除, 預設不會刪除檔案
	MaxBackups int    `yaml:"maxBackups"` // 日誌保留數量, 當日誌檔案數量超過此數量時會刪除舊檔案, 預設不會刪除檔案
	Compress   bool   `yaml:"compress"`   // 是否壓縮日誌檔案, 預設不會壓縮

	once   utils.SyncOnce // 單次執行物件
	logger *zap.Logger    // zap日誌物件
}

// Initialize 初始化處理
func (this *ZapLogger) Initialize() error {
	if this.once.Done() {
		return fmt.Errorf("zaplogger initialize: already initialize")
	} // if

	this.once.Do(func() {
		core := zapcore.NewCore(
			this.encoder(),
			this.writeSyncer(),
			zap.NewAtomicLevelAt(zapLevel(this.Level)),
		)
		this.logger = zap.New(core, zap.AddCallerSkip(1))
	})

	return nil
}

// Finalize 結束處理
func (this *ZapLogger) Finalize() {
	if this.once.Done() == false {
		return
	} // if

	if this.logger != nil {
		_ = this.logger.Sync()
	} // if
}

// New 建立日誌
func (this *ZapLogger) New(label string, level Level) Stream {
	return &ZapStream{
		logger: this.logger.Named(label),
		level:  zapLevel(level),
	}
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
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
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

// ZapStream zap記錄
type ZapStream struct {
	logger  *zap.Logger   // 日誌物件
	level   zapcore.Level // 日誌等級
	message string        // 訊息字串
	field   []zap.Field   // 記錄列表
}

// Message 記錄訊息
func (this *ZapStream) Message(format string, a ...any) Stream {
	this.message = fmt.Sprintf(format, a...)
	return this
}

// KV 記錄索引與數值
func (this *ZapStream) KV(key string, value any) Stream { //nolint
	switch v := value.(type) {
	case int8:
		this.field = append(this.field, zap.Int8(key, v))
	case uint8:
		this.field = append(this.field, zap.Uint8(key, v))
	case *int8:
		this.field = append(this.field, zap.Int8p(key, v))
	case *uint8:
		this.field = append(this.field, zap.Uint8p(key, v))
	case []int8:
		this.field = append(this.field, zap.Int8s(key, v))
	case []byte:
		this.field = append(this.field, zap.Binary(key, v))

	case int16:
		this.field = append(this.field, zap.Int16(key, v))
	case uint16:
		this.field = append(this.field, zap.Uint16(key, v))
	case *int16:
		this.field = append(this.field, zap.Int16p(key, v))
	case *uint16:
		this.field = append(this.field, zap.Uint16p(key, v))
	case []int16:
		this.field = append(this.field, zap.Int16s(key, v))
	case []uint16:
		this.field = append(this.field, zap.Uint16s(key, v))

	case int32:
		this.field = append(this.field, zap.Int32(key, v))
	case uint32:
		this.field = append(this.field, zap.Uint32(key, v))
	case *int32:
		this.field = append(this.field, zap.Int32p(key, v))
	case *uint32:
		this.field = append(this.field, zap.Uint32p(key, v))
	case []int32:
		this.field = append(this.field, zap.Int32s(key, v))
	case []uint32:
		this.field = append(this.field, zap.Uint32s(key, v))

	case int64:
		this.field = append(this.field, zap.Int64(key, v))
	case uint64:
		this.field = append(this.field, zap.Uint64(key, v))
	case *int64:
		this.field = append(this.field, zap.Int64p(key, v))
	case *uint64:
		this.field = append(this.field, zap.Uint64p(key, v))
	case []int64:
		this.field = append(this.field, zap.Int64s(key, v))
	case []uint64:
		this.field = append(this.field, zap.Uint64s(key, v))

	case int:
		this.field = append(this.field, zap.Int(key, v))
	case uint:
		this.field = append(this.field, zap.Uint(key, v))
	case *int:
		this.field = append(this.field, zap.Intp(key, v))
	case *uint:
		this.field = append(this.field, zap.Uintp(key, v))
	case []int:
		this.field = append(this.field, zap.Ints(key, v))
	case []uint:
		this.field = append(this.field, zap.Uints(key, v))

	case float32:
		this.field = append(this.field, zap.Float32(key, v))
	case *float32:
		this.field = append(this.field, zap.Float32p(key, v))
	case []float32:
		this.field = append(this.field, zap.Float32s(key, v))

	case float64:
		this.field = append(this.field, zap.Float64(key, v))
	case *float64:
		this.field = append(this.field, zap.Float64p(key, v))
	case []float64:
		this.field = append(this.field, zap.Float64s(key, v))

	case complex64:
		this.field = append(this.field, zap.Complex64(key, v))
	case *complex64:
		this.field = append(this.field, zap.Complex64p(key, v))
	case []complex64:
		this.field = append(this.field, zap.Complex64s(key, v))

	case complex128:
		this.field = append(this.field, zap.Complex128(key, v))
	case *complex128:
		this.field = append(this.field, zap.Complex128p(key, v))
	case []complex128:
		this.field = append(this.field, zap.Complex128s(key, v))

	case string:
		this.field = append(this.field, zap.String(key, v))
	case *string:
		this.field = append(this.field, zap.Stringp(key, v))
	case []string:
		this.field = append(this.field, zap.Strings(key, v))

	case bool:
		this.field = append(this.field, zap.Bool(key, v))
	case *bool:
		this.field = append(this.field, zap.Boolp(key, v))
	case []bool:
		this.field = append(this.field, zap.Bools(key, v))

	case uintptr:
		this.field = append(this.field, zap.Uintptr(key, v))
	case *uintptr:
		this.field = append(this.field, zap.Uintptrp(key, v))
	case []uintptr:
		this.field = append(this.field, zap.Uintptrs(key, v))

	default:
		this.field = append(this.field, zap.Reflect(key, v))
	} // switch

	return this
}

// Error 記錄錯誤
func (this *ZapStream) Error(err error) Stream {
	this.field = append(this.field, zap.Error(err))
	return this
}

// EndError 以錯誤結束記錄
func (this *ZapStream) EndError(err error) error {
	this.Error(err).End()
	return err
}

// End 結束記錄
func (this *ZapStream) End() {
	if l := this.logger.Check(this.level, this.message); l != nil {
		l.Write(this.field...)
	} // if
}

// zapLevel 日誌等級轉換為zap日誌等級
func zapLevel(level Level) zapcore.Level {
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

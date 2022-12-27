package configs

import (
	"bytes"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/spf13/viper"
)

// NewConfigmgr 建立配置管理器
func NewConfigmgr() *Configmgr {
	return &Configmgr{}
}

// Configmgr 配置管理器, 其實是對viper配置函式庫的包裝
type Configmgr struct {
	merge atomic.Bool // 合併旗標, 用來判斷從檔案讀取配置時, 要使用ReadInConfig還是MergeInConfig
}

// Reset 重置配置管理器
func (this *Configmgr) Reset() {
	viper.Reset()
}

// AddPath 新增配置路徑
func (this *Configmgr) AddPath(path ...string) {
	for _, itor := range path {
		viper.AddConfigPath(itor)
	} // for
}

// ReadFile 從檔案讀取配置, ext可選擇的項目可以參考viper.SupportedExts陣列
func (this *Configmgr) ReadFile(name, ext string) (err error) {
	viper.SetConfigName(name)
	viper.SetConfigType(ext)

	if this.merge.CompareAndSwap(false, true) {
		err = viper.ReadInConfig()
	} else {
		err = viper.MergeInConfig()
	} // if

	if err != nil {
		return fmt.Errorf("configmgr readfile: %v(%v): %w", name, ext, err)
	} // if

	return nil
}

// ReadString 從字串讀取配置
func (this *Configmgr) ReadString(value, ext string) (err error) {
	reader := bytes.NewBuffer([]byte(value))
	viper.SetConfigType(ext)

	if this.merge.CompareAndSwap(false, true) {
		err = viper.ReadConfig(reader)
	} else {
		err = viper.MergeConfig(reader)
	} // if

	if err != nil {
		return fmt.Errorf("configmgr readstring: %v: %w", ext, err)
	} // if

	return nil
}

// ReadBuffer 從讀取器讀取配置
func (this *Configmgr) ReadBuffer(reader io.Reader, ext string) (err error) {
	viper.SetConfigType(ext)

	if this.merge.CompareAndSwap(false, true) {
		err = viper.ReadConfig(reader)
	} else {
		err = viper.MergeConfig(reader)
	} // if

	if err != nil {
		return fmt.Errorf("configmgr readbuffer: %v: %w", ext, err)
	} // if

	return nil
}

// Get 取得配置
func (this *Configmgr) Get(key string) interface{} {
	return viper.Get(key)
}

// GetBool 取得布林值
func (this *Configmgr) GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetInt 取得整數
func (this *Configmgr) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetInt32 取得整數
func (this *Configmgr) GetInt32(key string) int32 {
	return viper.GetInt32(key)
}

// GetInt64 取得整數
func (this *Configmgr) GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

// GetUInt 取得整數
func (this *Configmgr) GetUInt(key string) uint {
	return viper.GetUint(key)
}

// GetUInt32 取得整數
func (this *Configmgr) GetUInt32(key string) uint32 {
	return viper.GetUint32(key)
}

// GetUInt64 取得整數
func (this *Configmgr) GetUInt64(key string) uint64 {
	return viper.GetUint64(key)
}

// GetFloat 取得浮點數
func (this *Configmgr) GetFloat(key string) float64 {
	return viper.GetFloat64(key)
}

// GetString 取得字串
func (this *Configmgr) GetString(key string) string {
	return viper.GetString(key)
}

// GetIntSlice 取得整數列表
func (this *Configmgr) GetIntSlice(key string) []int {
	return viper.GetIntSlice(key)
}

// GetStringSlice 取得字串列表
func (this *Configmgr) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// GetTime 取得時間
func (this *Configmgr) GetTime(key string) time.Time {
	return viper.GetTime(key)
}

// GetDuration 取得時間
func (this *Configmgr) GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

// GetSizeInBytes 取得位元長度
func (this *Configmgr) GetSizeInBytes(key string) uint {
	return viper.GetSizeInBytes(key)
}

// Unmarshal 反序列化為資料物件
func (this *Configmgr) Unmarshal(key string, obj interface{}) error {
	if viper.InConfig(key) == false {
		return fmt.Errorf("configmgr unmarshal: %v: not exist", key)
	} // if

	if err := viper.UnmarshalKey(key, obj); err != nil {
		return fmt.Errorf("configmgr unmarshal: %v: %w", key, err)
	} // if

	return nil
}

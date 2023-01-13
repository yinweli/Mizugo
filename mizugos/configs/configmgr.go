package configs

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// 配置管理器, 其實是對viper配置函式庫的包裝, 有以下幾種讀取配置的模式
// * 從檔案讀取配置
//   從設定好的路徑中讀取符合檔名與副檔名的檔案內容作為配置資料
//   首先用AddPath函式設置路徑, 可多次設置來指定多個路徑來源
//   接著用ReadFile函式設置檔案名稱與副檔名來嘗試讀取配置
// * 從字串讀取配置
//   從外部提供字串作為配置來源, 由於內部仍然用讀取檔案的方式來處理字串, 所以必須提供來源使用的檔案格式的副檔名
//   用ReadString函式設置字串與副檔名來嘗試讀取配置
// * 從讀取器讀取配置
//   從外部提供讀取器作為配置來源, 由於內部仍然用讀取檔案的方式來處理字串, 所以必須提供來源使用的檔案格式的副檔名
// * 支援的副檔名如下(可以參考viper.SupportedExts陣列)
//   - dotenv
//   - env
//   - hcl
//   - ini
//   - json
//   - prop
//   - properties
//   - props
//   - tfvars
//   - toml
//   - yaml
//   - yml

// NewConfigmgr 建立配置管理器
func NewConfigmgr() *Configmgr {
	return &Configmgr{}
}

// Configmgr 配置管理器
type Configmgr struct {
	once sync.Once // 單次執行鎖
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

	first := false // 首次旗標, 用來判斷是否首次讀取配置, 決定要用ReadInConfig還是MergeInConfig

	this.once.Do(func() {
		first = true
		err = viper.ReadInConfig()
	})

	if first == false {
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

	first := false // 首次旗標, 用來判斷是否首次讀取配置, 決定要用ReadConfig還是MergeConfig

	this.once.Do(func() {
		first = true
		err = viper.ReadConfig(reader)
	})

	if first == false {
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

	first := false // 首次旗標, 用來判斷是否首次讀取配置, 決定要用ReadConfig還是MergeConfig

	this.once.Do(func() {
		first = true
		err = viper.ReadConfig(reader)
	})

	if first == false {
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

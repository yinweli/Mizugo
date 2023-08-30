package configs

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/spf13/viper"
)

// NewConfigmgr 建立配置管理器
func NewConfigmgr() *Configmgr {
	return &Configmgr{
		config: viper.New(),
	}
}

// Configmgr 配置管理器, 內部使用viper實現功能, 有以下幾種讀取配置的模式
//   - 從檔案讀取配置: 從指定的路徑/檔名/副檔名讀取配置, 需要配合 AddPath 函式設置路徑(可多次設置來指定多個路徑來源)
//   - 從字串讀取配置: 從字串讀取配置, 由於內部仍然用讀取檔案的方式來處理字串, 所以必須提供來源使用的檔案格式的副檔名
//   - 從讀取器讀取配置: 從讀取器讀取配置, 由於內部仍然用讀取檔案的方式來處理字串, 所以必須提供來源使用的檔案格式的副檔名
//   - 從環境變數讀取配置: 從環境變數中讀取配置, 這需要事先決定好前綴字
//
// 以上讀取配置的模式內部都是用讀取檔案的方式來處理字串, 所以必須提供來源使用的檔案格式的副檔名, 支援的副檔名可以參考 viper.SupportedExts
//
// 當配置讀取完畢後, 需要從管理器中取得配置值時, 可以用索引字串來呼叫 Get... 系列函式來取得配置值;
// 或是用索引字串來呼叫 Unmarshal 來取得反序列化到結構的配置資料
//
// 當同時使用檔案/字串/讀取器以及環境變數時, viper有規定其使用順序, 上面的會覆蓋下面的配置:
//   - 環境變數
//   - 配置檔案
type Configmgr struct {
	config *viper.Viper // 配置資料
	read   bool         // 讀取旗標
}

// Reset 重置配置管理器
func (this *Configmgr) Reset() {
	viper.Reset() // 其實這個不執行也沒關係, 但是因為它會重置 viper.SupportedExts 跟 viper.SupportedRemoteProviders, 所以還是跑一下好了
	this.config = viper.New()
	this.read = false
}

// AddPath 新增配置路徑, 與 ReadFile 搭配使用, 可多次設置來指定多個路徑來源
func (this *Configmgr) AddPath(path ...string) {
	for _, itor := range path {
		this.config.AddConfigPath(itor)
	} // for
}

// ReadFile 從檔案讀取配置, 可以設定多個路徑來源, 支援的副檔名可以參考 viper.SupportedExts
func (this *Configmgr) ReadFile(name, ext string) (err error) {
	this.config.SetConfigType(ext)
	this.config.SetConfigName(name)

	if this.read == false {
		this.read = true
		err = this.config.ReadInConfig()
	} else {
		err = this.config.MergeInConfig()
	} // if

	if err != nil {
		return fmt.Errorf("configmgr readfile: %v(%v): %w", name, ext, err)
	} // if

	return nil
}

// ReadString 從字串讀取配置, 支援的副檔名可以參考 viper.SupportedExts
func (this *Configmgr) ReadString(input, ext string) (err error) {
	reader := bytes.NewBuffer([]byte(input))
	this.config.SetConfigType(ext)

	if this.read == false {
		this.read = true
		err = this.config.ReadConfig(reader)
	} else {
		err = this.config.MergeConfig(reader)
	} // if

	if err != nil {
		return fmt.Errorf("configmgr readstring: %v: %w", ext, err)
	} // if

	return nil
}

// ReadBuffer 從讀取器讀取配置, 支援的副檔名可以參考 viper.SupportedExts
func (this *Configmgr) ReadBuffer(reader io.Reader, ext string) (err error) {
	this.config.SetConfigType(ext)

	if this.read == false {
		this.read = true
		err = this.config.ReadConfig(reader)
	} else {
		err = this.config.MergeConfig(reader)
	} // if

	if err != nil {
		return fmt.Errorf("configmgr readbuffer: %v: %w", ext, err)
	} // if

	return nil
}

// ReadEnvironment 從環境變數讀取配置, 讀取時需要指定前綴字;
// 請注意在整個程式中只使用一個前綴字, 不然就得要在每次取得配置前, 都要執行 ReadEnvironment 來重新指定前綴字;
// 前綴字與第一層key之間將會以`_`連接, 第一層key與之後的key之間則以`.`連接
// 以下是個配置檔案與環境變數的映射範例, 假設我們將前綴字設定為 MYAPP
//
//	test:
//	  value1: a // 環境變數名稱為 MYAPP_TEST.VALUE1
//	  value2: b // 環境變數名稱為 MYAPP_TEST.VALUE2
//
// 在上面的例子中, `test`是第一層key, `value1`與`value2`則是之後的key;
// 環境變數在讀取時將會轉為全大寫來搜尋
func (this *Configmgr) ReadEnvironment(prefix string) (err error) {
	this.config.SetEnvPrefix(prefix)
	this.config.AutomaticEnv()
	return nil
}

// Get 取得配置
func (this *Configmgr) Get(key string) any {
	return this.config.Get(key)
}

// GetBool 取得布林值
func (this *Configmgr) GetBool(key string) bool {
	return this.config.GetBool(key)
}

// GetInt 取得整數
func (this *Configmgr) GetInt(key string) int {
	return this.config.GetInt(key)
}

// GetInt32 取得整數
func (this *Configmgr) GetInt32(key string) int32 {
	return this.config.GetInt32(key)
}

// GetInt64 取得整數
func (this *Configmgr) GetInt64(key string) int64 {
	return this.config.GetInt64(key)
}

// GetUInt 取得整數
func (this *Configmgr) GetUInt(key string) uint {
	return this.config.GetUint(key)
}

// GetUInt32 取得整數
func (this *Configmgr) GetUInt32(key string) uint32 {
	return this.config.GetUint32(key)
}

// GetUInt64 取得整數
func (this *Configmgr) GetUInt64(key string) uint64 {
	return this.config.GetUint64(key)
}

// GetFloat 取得浮點數
func (this *Configmgr) GetFloat(key string) float64 {
	return this.config.GetFloat64(key)
}

// GetString 取得字串
func (this *Configmgr) GetString(key string) string {
	return this.config.GetString(key)
}

// GetIntSlice 取得整數列表
func (this *Configmgr) GetIntSlice(key string) []int {
	return this.config.GetIntSlice(key)
}

// GetStringSlice 取得字串列表
func (this *Configmgr) GetStringSlice(key string) []string {
	return this.config.GetStringSlice(key)
}

// GetTime 取得時間
func (this *Configmgr) GetTime(key string) time.Time {
	return this.config.GetTime(key)
}

// GetDuration 取得時間
func (this *Configmgr) GetDuration(key string) time.Duration {
	return this.config.GetDuration(key)
}

// GetSizeInBytes 取得位元長度
func (this *Configmgr) GetSizeInBytes(key string) uint {
	return this.config.GetSizeInBytes(key)
}

// Unmarshal 反序列化為資料物件
func (this *Configmgr) Unmarshal(key string, obj any) error {
	if this.config.InConfig(key) == false {
		return fmt.Errorf("configmgr unmarshal: %v: not exist", key)
	} // if

	if err := this.config.UnmarshalKey(key, obj); err != nil {
		return fmt.Errorf("configmgr unmarshal: %v: %w", key, err)
	} // if

	return nil
}

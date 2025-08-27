package configs

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/viper"
)

// NewConfigmgr 建立配置管理器
func NewConfigmgr() *Configmgr {
	return &Configmgr{
		config: viper.New(),
	}
}

// Configmgr 配置管理器, 是一個基於 viper 的配置管理器, 支援多種載入來源:
//   - 檔案: 透過 AddPath 設定搜尋路徑, 指定檔名與副檔名載入
//   - 字串: 以字串內容載入, 內部仍走「讀檔流程」, 需提供對應的副檔名
//   - 讀取器: 自 io.Reader 載入, 同樣需提供副檔名
//   - 環境變數: 以指定前綴字讀取環境變數
//
// 由於上述載入方式在內部皆以「讀檔模式」處理, 必須提供正確的副檔名, 支援清單可參考 viper.SupportedExts
//
// 讀取完成後, 可使用 Unmarshal(key, obj) 將指定節點反序列化為結構體
//
// 多來源同時存在時, viper 的覆蓋優先序如下(上者覆蓋下者):
//   - 環境變數
//   - 設定檔案
type Configmgr struct {
	config *viper.Viper // 配置資料
	read   bool         // 讀取旗標
}

// Reset 重置配置管理器
func (this *Configmgr) Reset() {
	this.config = viper.New()
	this.read = false
}

// AddPath 新增配置路徑, 與 ReadFile 搭配使用, 可多次設置來指定多個路徑來源
func (this *Configmgr) AddPath(path ...string) {
	for _, itor := range path {
		this.config.AddConfigPath(itor)
	} // for
}

// ReadFile 從檔案讀取配置, 與 AddPath 搭配使用, 支援的副檔名可以參考 viper.SupportedExts
func (this *Configmgr) ReadFile(name, ext string) (err error) {
	if ext == "" {
		return fmt.Errorf("configmgr readfile: ext empty")
	} // if

	this.config.SetConfigType(ext)
	this.config.SetConfigName(name)

	if this.read == false {
		this.read = true
		err = this.config.ReadInConfig()
	} else {
		err = this.config.MergeInConfig()
	} // if

	if err != nil {
		return fmt.Errorf("configmgr readfile: %v.%v: %w", name, ext, err)
	} // if

	return nil
}

// ReadString 從字串讀取配置, 支援的副檔名可以參考 viper.SupportedExts
func (this *Configmgr) ReadString(input, ext string) (err error) {
	if ext == "" {
		return fmt.Errorf("configmgr readstring: ext empty")
	} // if

	reader := bytes.NewBufferString(input)
	this.config.SetConfigType(ext)

	if this.read == false {
		this.read = true
		err = this.config.ReadConfig(reader)
	} else {
		err = this.config.MergeConfig(reader)
	} // if

	if err != nil {
		return fmt.Errorf("configmgr readstring: %w", err)
	} // if

	return nil
}

// ReadBuffer 從讀取器讀取配置, 支援的副檔名可以參考 viper.SupportedExts
func (this *Configmgr) ReadBuffer(reader io.Reader, ext string) (err error) {
	if ext == "" {
		return fmt.Errorf("configmgr readbuffer: ext empty")
	} // if

	this.config.SetConfigType(ext)

	if this.read == false {
		this.read = true
		err = this.config.ReadConfig(reader)
	} else {
		err = this.config.MergeConfig(reader)
	} // if

	if err != nil {
		return fmt.Errorf("configmgr readbuffer: %w", err)
	} // if

	return nil
}

// ReadEnvironment 從環境變數載入配置, 需指定前綴字; 建議整個程式僅使用單一前綴字, 否則每次取值前都必須重新呼叫此方法設定前綴
//
// 命名規則:
//   - 前綴字與每一層 key 以 `_` 連接
//   - 環境變數名稱會自動轉為大寫
//
// 範例: 若 prefix = "MYAPP", 則以下 YAML 會映射為環境變數:
//
//	test:
//	  value1: a // MYAPP_TEST_VALUE1
//	  value2: b // MYAPP_TEST_VALUE2
func (this *Configmgr) ReadEnvironment(prefix string) (err error) {
	if prefix == "" {
		return fmt.Errorf("configmgr readenvironment: prefix empty")
	} // if

	this.config.SetEnvPrefix(prefix)
	this.config.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 為了使用方便, 讓每一層 key 的連接都使用 `_`
	this.config.AutomaticEnv()
	return nil
}

// Unmarshal 反序列化為資料物件
func (this *Configmgr) Unmarshal(key string, obj any) error {
	if this.config.IsSet(key) == false {
		return fmt.Errorf("configmgr unmarshal: %v: not exist", key)
	} // if

	if err := this.config.UnmarshalKey(key, obj); err != nil {
		return fmt.Errorf("configmgr unmarshal: %v: %w", key, err)
	} // if

	return nil
}

// Get 取得配置
func (this *Configmgr) Get(key string) any {
	return this.config.Get(key)
}

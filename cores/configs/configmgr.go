package configs

import (
	"fmt"
	"os"
	"sync"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

// NewConfigmgr 建立配置管理器
func NewConfigmgr() *Configmgr {
	return &Configmgr{
		data: map[string]interface{}{},
	}
}

// Configmgr 配置管理器
type Configmgr struct {
	data map[string]interface{} // 配置列表
	lock sync.Mutex             // 執行緒鎖
}

// ReadFile 從檔案讀取配置
func (this *Configmgr) ReadFile(filepath string) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	data, err := os.ReadFile(filepath)

	if err != nil {
		return fmt.Errorf("configmgr readfile: %v: %w", filepath, err)
	} // if

	if err = this.ReadData(data); err != nil {
		return fmt.Errorf("configmgr readfile: %v: %w", filepath, err)
	} // if

	return nil
}

// ReadString 從字串讀取配置
func (this *Configmgr) ReadString(str string) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.ReadData([]byte(str)); err != nil {
		return fmt.Errorf("configmgr readstring: %w", err)
	} // if

	return nil
}

// ReadData 從二進位資料讀取配置
func (this *Configmgr) ReadData(data []byte) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	config := map[string]interface{}{}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("configmgr readdata: %w", err)
	} // if

	for key, value := range config {
		this.data[key] = value
	} // for

	return nil
}

// Clear 清除配置
func (this *Configmgr) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.data = map[string]interface{}{}
}

// GetInt 取得數字
func (this *Configmgr) GetInt(key string) (result int, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	raw, ok := this.data[key]

	if ok == false {
		return 0, fmt.Errorf("configmgr getint: not exist")
	} // if

	value, ok := raw.(int)

	if ok == false {
		return 0, fmt.Errorf("configmgr getint: not int")
	} // if

	return value, nil
}

// GetString 取得字串
func (this *Configmgr) GetString(key string) (result string, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	raw, ok := this.data[key]

	if ok == false {
		return "", fmt.Errorf("configmgr getstring: not exist")
	} // if

	value, ok := raw.(string)

	if ok == false {
		return "", fmt.Errorf("configmgr getstring: not string")
	} // if

	return value, nil
}

// GetObject 取得物件
func (this *Configmgr) GetObject(key string, result interface{}) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	raw, ok := this.data[key]

	if ok == false {
		return fmt.Errorf("configmgr getobject: not exist")
	} // if

	if err := mapstructure.Decode(raw, &result); err != nil { // TODO: 這樣會成功嗎!? 要測試看看
		return fmt.Errorf("configmgr getobject: %w", err)
	} // if

	return nil
}

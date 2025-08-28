package redmos

import (
	"fmt"
	"sync"
)

// NewRedmomgr 建立資料庫管理器
func NewRedmomgr() *Redmomgr {
	return &Redmomgr{
		major: map[string]*Major{},
		minor: map[string]*Minor{},
		mixed: map[string]*Mixed{},
	}
}

// Redmomgr 資料庫管理器
//
// 用於管理雙層式資料庫架構, 包含:
//   - 主要資料庫(Major): 以 Redis 為基礎, 用於快取或高速查詢
//   - 次要資料庫(Minor): 以 Mongo 為基礎, 用於持久化或複雜查詢
//   - 混合資料庫(Mixed): 由 Major + Minor 組合而成, 整合兩者特性
//
// 使用流程:
//  1. 呼叫 AddMajor / AddMinor 新增資料庫
//  2. 呼叫 AddMixed 將已建立的 Major 與 Minor 綁定成混合資料庫
//  3. 呼叫 GetMajor / GetMinor / GetMixed 取得對應的資料庫物件並操作
//
// 注意事項:
//   - Mixed 必須同時依賴一個 Major 與一個 Minor, 缺少任一方則無法建立
//   - Add* 系列函式若遇到重複名稱會回傳錯誤
//   - Finalize 會關閉並清理所有已註冊的資料庫
type Redmomgr struct {
	major map[string]*Major // 主要資料庫列表(Redis)
	minor map[string]*Minor // 次要資料庫列表(Mongo)
	mixed map[string]*Mixed // 混合資料庫列表(Redis+Mongo)
	lock  sync.RWMutex      // 執行緒鎖
}

// AddMajor 新增主要資料庫
func (this *Redmomgr) AddMajor(majorName string, uri RedisURI) (major *Major, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.major[majorName]; ok {
		return nil, fmt.Errorf("redmomgr addMajor: duplicate database")
	} // if

	major, err = newMajor(uri)

	if err != nil {
		return nil, fmt.Errorf("redmomgr addMajor: %w", err)
	} // if

	this.major[majorName] = major
	return major, nil
}

// GetMajor 取得主要資料庫
func (this *Redmomgr) GetMajor(majorName string) *Major {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if major, ok := this.major[majorName]; ok {
		return major
	} // if

	return nil
}

// AddMinor 新增次要資料庫
func (this *Redmomgr) AddMinor(minorName string, uri MongoURI, dbName string) (minor *Minor, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.minor[minorName]; ok {
		return nil, fmt.Errorf("redmomgr addMinor: duplicate database")
	} // if

	minor, err = newMinor(uri, dbName)

	if err != nil {
		return nil, fmt.Errorf("redmomgr addMinor: %w", err)
	} // if

	this.minor[minorName] = minor
	return minor, nil
}

// GetMinor 取得次要資料庫
func (this *Redmomgr) GetMinor(minorName string) *Minor {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if minor, ok := this.minor[minorName]; ok {
		return minor
	} // if

	return nil
}

// AddMixed 新增混合資料庫
func (this *Redmomgr) AddMixed(mixedName, majorName, minorName string) (mixed *Mixed, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.mixed[mixedName]; ok {
		return nil, fmt.Errorf("redmomgr addMixed: duplicate database")
	} // if

	major, ok := this.major[majorName]

	if ok == false {
		return nil, fmt.Errorf("redmomgr addMixed: major not exist")
	} // if

	minor, ok := this.minor[minorName]

	if ok == false {
		return nil, fmt.Errorf("redmomgr addMixed: minor not exist")
	} // if

	mixed = newMixed(major, minor)
	this.mixed[mixedName] = mixed
	return mixed, nil
}

// GetMixed 取得混合資料庫
func (this *Redmomgr) GetMixed(mixedName string) *Mixed {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if mixed, ok := this.mixed[mixedName]; ok {
		return mixed
	} // if

	return nil
}

// Finalize 結束處理
func (this *Redmomgr) Finalize() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range this.major {
		itor.stop()
	} // if

	for _, itor := range this.minor {
		itor.stop()
	} // if

	this.major = map[string]*Major{}
	this.minor = map[string]*Minor{}
	this.mixed = map[string]*Mixed{}
}

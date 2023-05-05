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

// Redmomgr 資料庫管理器, 用於管理雙層式資料庫架構
//   - 主要資料庫: 用redis實作
//   - 次要資料庫: 用mongo實作
//
// 當要新增資料庫時, 需要遵循以下流程:
//   - 新增主要/次要資料庫
//   - 新增混合資料庫: 這時會去取得先前新增的主要/次要資料庫, 並且將其`綁定`到混合資料庫中;
//     要注意的是, 混合資料庫必定是一個主要資料庫加上一個次要資料庫的組合, 若是缺少了任何一方則會失敗
//
// 若要執行資料庫操作時, 呼叫 Get... 系列函式來取得資料庫物件
type Redmomgr struct {
	major map[string]*Major // 主要資料庫列表
	minor map[string]*Minor // 次要資料庫列表
	mixed map[string]*Mixed // 混合資料庫列表
	lock  sync.RWMutex      // 執行緒鎖
}

// AddMajor 新增主要資料庫, 需要提供 RedisURI 來指定要連接的資料庫以及連接選項
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

// AddMinor 新增次要資料庫, 需要提供 MongoURI 來指定要連接的資料庫以及連接選項;
// 另外需要指定mongo資料庫名稱, 簡化後面取得執行器的流程, 但也因此限制次要資料庫不能在多個mongo資料庫間切換
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

// AddMixed 新增混合資料庫, 必須確保 majorName 與 minorName 必須是先前建立好的資料庫, 否則會失敗
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

	mixed = &Mixed{
		major: major,
		minor: minor,
	}
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

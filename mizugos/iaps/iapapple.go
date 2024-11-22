package iaps

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/awa/go-iap/appstore/api"
)

// NewIAPApple 建立Apple驗證器
func NewIAPApple(config *IAPAppleConfig) *IAPApple {
	return &IAPApple{
		config: config,
	}
}

// IAPApple Apple驗證器
type IAPApple struct {
	config *IAPAppleConfig  // 驗證設定
	client *api.StoreClient // 驗證客戶端
	verify chan *iapApple   // 驗證通道
	signal sync.WaitGroup   // 通知信號
}

// IAPAppleConfig Apple驗證設定資料
type IAPAppleConfig struct {
	Key      string        `yaml:"key"`      // 密鑰字串
	KeyID    string        `yaml:"keyID"`    // 密鑰ID
	Bundle   string        `yaml:"bundle"`   // 軟體包名稱
	Issuer   string        `yaml:"issuer"`   // 發行人名稱
	Sandbox  bool          `yaml:"sandbox"`  // 沙盒旗標
	Capacity int           `yaml:"capacity"` // 通道容量
	Retry    int           `yaml:"retry"`    // 重試次數
	Timeout  time.Duration `yaml:"timeout"`  // 驗證逾時時間
	Interval time.Duration `yaml:"interval"` // 驗證間隔時間
}

// IAPAppleResult Apple驗證結果資料
type IAPAppleResult struct {
	Err  error     // 驗證結果, 若為nil表示驗證成功, 否則失敗
	Time time.Time // 購買時間
}

// iapApple Apple驗證資料
type iapApple struct {
	productID   string              // 產品編號
	certificate string              // 購買憑證
	retry       int                 // 重試次數
	retryErr    error               // 重試錯誤
	result      chan IAPAppleResult // 結果通道
}

// Initialize 初始化處理
func (this *IAPApple) Initialize() error {
	this.client = api.NewStoreClient(&api.StoreConfig{
		KeyContent: []byte(this.config.Key),
		KeyID:      this.config.KeyID,
		BundleID:   this.config.Bundle,
		Issuer:     this.config.Issuer,
		Sandbox:    this.config.Sandbox,
	})
	this.verify = make(chan *iapApple, this.config.Capacity+1) // 避免使用者將通道容量設為0導致卡住
	this.signal.Add(1)
	go this.execute(this.verify)
	return nil
}

// Finalize 結束處理
func (this *IAPApple) Finalize() {
	close(this.verify)
	this.verify = nil
	this.signal.Wait()
}

// Verify 驗證憑證
func (this *IAPApple) Verify(productID, certificate string) IAPAppleResult {
	if this.verify == nil {
		return this.fail(fmt.Errorf("iapApple verify: close"))
	} // if

	result := &iapApple{
		productID:   productID,
		certificate: certificate,
		result:      make(chan IAPAppleResult),
	}
	this.verify <- result
	return <-result.result
}

// execute 執行驗證
func (this *IAPApple) execute(verify chan *iapApple) {
	for itor := range verify {
		// 由於驗證api有速率限制, 所以需要等待後才能繼續下一個驗證
		time.Sleep(this.config.Interval)

		// 如果重試超過限制, 還是只能當作錯誤
		if itor.retry > 0 && itor.retry >= this.config.Retry {
			itor.result <- this.fail(itor.retryErr)
			continue
		} // if

		ctx, cancel := context.WithTimeout(context.Background(), this.config.Timeout)
		respond, err := this.client.GetTransactionInfo(ctx, itor.certificate)
		cancel() // 避免cancel洩漏

		// 由於偶爾會出現驗證資料都填寫正確, 但是驗證api卻回應錯誤, 所以只好在這裡重複嘗試
		if err != nil {
			itor.retry++
			itor.retryErr = fmt.Errorf("iapApple execute: %w", err)
			verify <- itor
			continue
		} // if

		result, err := this.client.ParseSignedTransaction(respond.SignedTransactionInfo)

		if err != nil {
			itor.result <- this.fail(fmt.Errorf("iapApple execute: %w", err))
			continue
		} // if

		if result.ProductID != itor.productID {
			itor.result <- this.fail(fmt.Errorf("iapApple execute: productID"))
			continue
		} // if

		if result.TransactionID != itor.certificate {
			itor.result <- this.fail(fmt.Errorf("iapApple execute: certificate"))
			continue
		} // if

		itor.result <- this.succ(result.PurchaseDate)
	} // for

	this.signal.Done()
}

// succ 取得成功物件
func (this *IAPApple) succ(millisecond int64) IAPAppleResult {
	secs := millisecond / 1000
	nano := (millisecond % 1000) * int64(time.Millisecond)
	return IAPAppleResult{
		Time: time.Unix(secs, nano),
	}
}

// fail 取得失敗物件
func (this *IAPApple) fail(err error) IAPAppleResult {
	return IAPAppleResult{
		Err: err,
	}
}

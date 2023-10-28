package iaps

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/awa/go-iap/appstore/api"
)

// NewIAPApple 建立apple驗證器
func NewIAPApple(config *IAPAppleConfig) *IAPApple {
	return &IAPApple{
		config: config,
	}
}

// IAPApple apple驗證器
type IAPApple struct {
	config *IAPAppleConfig  // 驗證設定
	client *api.StoreClient // 驗證客戶端
	verify chan *iapApple   // 驗證通道
	signal sync.WaitGroup   // 通知信號
}

// IAPAppleConfig apple驗證設定資料
type IAPAppleConfig struct {
	Key      string        `yaml:"key"`      // 密鑰字串
	KeyID    string        `yaml:"keyID"`    // 密鑰ID
	Bundle   string        `yaml:"bundle"`   // 軟體包名稱
	Issuer   string        `yaml:"issuer"`   // 發行人名稱
	Sandbox  bool          `yaml:"sandbox"`  // 沙盒旗標
	WaitTime time.Duration `yaml:"waitTime"` // 等待時間
	Capacity int           `yaml:"capacity"` // 通道容量
	Retry    int           `yaml:"retry"`    // 重試次數
}

// iapApple apple驗證資料
type iapApple struct {
	productID   string     // 產品編號
	certificate string     // 購買憑證
	retry       int        // 重試次數
	retryErr    error      // 重試錯誤
	result      chan error // 結果通道
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
	this.verify = make(chan *iapApple, this.config.Capacity)
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
func (this *IAPApple) Verify(productID, certificate string) error {
	if this.verify == nil {
		return fmt.Errorf("iapApple verify: close")
	} // if

	result := &iapApple{
		productID:   productID,
		certificate: certificate,
		result:      make(chan error),
	}
	this.verify <- result
	return <-result.result
}

// execute 執行驗證
func (this *IAPApple) execute(verify chan *iapApple) {
	for itor := range verify {
		time.Sleep(this.config.WaitTime) // 由於驗證api有速率限制, 所以需要等待後才能繼續下一個驗證

		if itor.retry >= this.config.Retry { // 如果重試超過限制, 還是只能當作錯誤
			itor.result <- itor.retryErr
			continue
		} // if

		// 由於測試時 ctxs.Get() 常會被奇怪的關閉, 所以這裡使用正常的ctx
		respond, err := this.client.GetTransactionInfo(context.Background(), itor.certificate)

		if err != nil { // 由於偶爾會出現驗證資料都填寫正確, 但是驗證api卻回應錯誤, 所以只好在這裡重複嘗試
			itor.retry++
			itor.retryErr = fmt.Errorf("iapApple execute: %w", err)
			verify <- itor
			continue
		} // if

		result, err := this.client.ParseSignedTransaction(respond.SignedTransactionInfo)

		if err != nil {
			itor.result <- fmt.Errorf("iapApple execute: %w", err)
			continue
		} // if

		if result.ProductID != itor.productID {
			itor.result <- fmt.Errorf("iapApple execute: productID")
			continue
		} // if

		if result.TransactionID != itor.certificate {
			itor.result <- fmt.Errorf("iapApple execute: certificate")
			continue
		} // if

		itor.result <- nil
	} // for

	this.signal.Done()
}

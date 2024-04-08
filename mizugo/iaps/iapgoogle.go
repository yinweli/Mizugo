package iaps

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/awa/go-iap/playstore"
)

// NewIAPGoogle 建立google驗證器
func NewIAPGoogle(config *IAPGoogleConfig) *IAPGoogle {
	return &IAPGoogle{
		config: config,
	}
}

// IAPGoogle google驗證器
type IAPGoogle struct {
	config *IAPGoogleConfig  // 驗證設定
	client *playstore.Client // 驗證客戶端
	verify chan *iapGoogle   // 驗證通道
	signal sync.WaitGroup    // 通知信號
}

// IAPGoogleConfig google驗證設定資料
type IAPGoogleConfig struct {
	Key      string        `yaml:"key"`      // 密鑰字串
	Bundle   string        `yaml:"bundle"`   // 軟體包名稱
	WaitTime time.Duration `yaml:"waitTime"` // 等待時間
	Capacity int           `yaml:"capacity"` // 通道容量
}

// iapGoogle google驗證資料
type iapGoogle struct {
	productID   string     // 產品編號
	certificate string     // 購買憑證
	result      chan error // 結果通道
}

// Initialize 初始化處理
func (this *IAPGoogle) Initialize() error {
	client, err := playstore.New([]byte(this.config.Key))

	if err != nil {
		return fmt.Errorf("iapApple initialize: %w", err)
	} // if

	this.client = client
	this.verify = make(chan *iapGoogle, this.config.Capacity)
	this.signal.Add(1)
	go this.execute(this.verify)
	return nil
}

// Finalize 結束處理
func (this *IAPGoogle) Finalize() {
	close(this.verify)
	this.verify = nil
	this.signal.Wait()
}

// Verify 驗證憑證
func (this *IAPGoogle) Verify(productID, certificate string) error {
	if this.verify == nil {
		return fmt.Errorf("iapGoogle verify: close")
	} // if

	result := &iapGoogle{
		productID:   productID,
		certificate: certificate,
		result:      make(chan error),
	}
	this.verify <- result
	return <-result.result
}

// execute 執行驗證
func (this *IAPGoogle) execute(verify chan *iapGoogle) {
	for itor := range verify {
		time.Sleep(this.config.WaitTime) // 由於驗證api有速率限制, 所以需要等待後才能繼續下一個驗證

		// 由於測試時 ctxs.Get() 常會被奇怪的關閉, 所以這裡使用正常的ctx
		_, err := this.client.VerifyProduct(context.Background(), this.config.Bundle, itor.productID, itor.certificate)

		if err != nil {
			itor.result <- fmt.Errorf("iapGoogle execute: %w", err)
			continue
		} // if

		itor.result <- nil
	} // for

	this.signal.Done()
}

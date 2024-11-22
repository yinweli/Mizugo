package iaps

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/awa/go-iap/playstore"
)

// NewIAPGoogle 建立Google驗證器
func NewIAPGoogle(config *IAPGoogleConfig) *IAPGoogle {
	return &IAPGoogle{
		config: config,
	}
}

// IAPGoogle Google驗證器
type IAPGoogle struct {
	config *IAPGoogleConfig  // 驗證設定
	client *playstore.Client // 驗證客戶端
	verify chan *iapGoogle   // 驗證通道
	signal sync.WaitGroup    // 通知信號
}

// IAPGoogleConfig Google驗證設定資料
type IAPGoogleConfig struct {
	Key      string        `yaml:"key"`      // 密鑰字串
	Bundle   string        `yaml:"bundle"`   // 軟體包名稱
	Capacity int           `yaml:"capacity"` // 通道容量
	Timeout  time.Duration `yaml:"timeout"`  // 驗證逾時時間
	Interval time.Duration `yaml:"interval"` // 驗證間隔時間
}

// IAPGoogleResult Google驗證結果資料
type IAPGoogleResult struct {
	Err  error     // 驗證結果, 若為nil表示驗證成功, 否則失敗
	Time time.Time // 購買時間
}

// iapGoogle Google驗證資料
type iapGoogle struct {
	productID   string               // 產品編號
	certificate string               // 購買憑證
	result      chan IAPGoogleResult // 結果通道
}

// Initialize 初始化處理
func (this *IAPGoogle) Initialize() error {
	client, err := playstore.New([]byte(this.config.Key))

	if err != nil {
		return fmt.Errorf("iapGoogle initialize: %w", err)
	} // if

	this.client = client
	this.verify = make(chan *iapGoogle, this.config.Capacity+1) // 避免使用者將通道容量設為0導致卡住
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
func (this *IAPGoogle) Verify(productID, certificate string) IAPGoogleResult {
	if this.verify == nil {
		return this.fail(fmt.Errorf("iapGoogle verify: close"))
	} // if

	result := &iapGoogle{
		productID:   productID,
		certificate: certificate,
		result:      make(chan IAPGoogleResult),
	}
	this.verify <- result
	return <-result.result
}

// execute 執行驗證
func (this *IAPGoogle) execute(verify chan *iapGoogle) {
	for itor := range verify {
		// 由於驗證api有速率限制, 所以需要等待後才能繼續下一個驗證
		time.Sleep(this.config.Interval)

		ctx, cancel := context.WithTimeout(context.Background(), this.config.Timeout)
		result, err := this.client.VerifyProduct(ctx, this.config.Bundle, itor.productID, itor.certificate)
		cancel() // 避免cancel洩漏

		if err != nil {
			itor.result <- this.fail(fmt.Errorf("iapGoogle execute: %w", err))
			continue
		} // if

		itor.result <- this.succ(result.PurchaseTimeMillis)
	} // for

	this.signal.Done()
}

// succ 取得成功物件
func (this *IAPGoogle) succ(millisecond int64) IAPGoogleResult {
	secs := millisecond / 1000
	nano := (millisecond % 1000) * int64(time.Millisecond)
	return IAPGoogleResult{
		Time: time.Unix(secs, nano),
	}
}

// fail 取得失敗物件
func (this *IAPGoogle) fail(err error) IAPGoogleResult {
	return IAPGoogleResult{
		Err: err,
	}
}

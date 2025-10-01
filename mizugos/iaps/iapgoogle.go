package iaps

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/awa/go-iap/playstore"
	"google.golang.org/api/androidpublisher/v3"
)

// NewIAPGoogle 建立 Google IAP 驗證器
func NewIAPGoogle(config *IAPGoogleConfig) *IAPGoogle {
	return &IAPGoogle{
		config: config,
	}
}

// IAPGoogle Google IAP 驗證器
type IAPGoogle struct {
	config *IAPGoogleConfig   // 驗證設定
	client IAPGoogleClient    // 驗證客戶端
	verify chan *iapGoogle    // 驗證通道
	signal sync.WaitGroup     // 通知信號
	ctx    context.Context    // 關閉物件
	cancel context.CancelFunc // 關閉函式
}

// IAPGoogleConfig Google IAP 驗證設定資料
type IAPGoogleConfig struct {
	Key    string `yaml:"key"`    // 服務帳號私鑰(JSON 格式)
	Bundle string `yaml:"bundle"` // App 的套件名稱(PackageName)
}

// IAPGoogleResult Google IAP 驗證結果資料
type IAPGoogleResult struct {
	Err  error     // 驗證結果, 若為nil表示驗證成功, 否則失敗
	Time time.Time // 購買時間
}

// IAPGoogleClient Google IAP 驗證客戶端介面
type IAPGoogleClient interface {
	VerifyProduct(context.Context, string, string, string) (*androidpublisher.ProductPurchase, error)
}

// iapGoogle Google IAP 驗證資料
type iapGoogle struct {
	productID string               // 產品編號
	receipt   string               // 購買憑證(PurchaseToken)
	result    chan IAPGoogleResult // 驗證結果通道
}

// Initialize 初始化處理
func (this *IAPGoogle) Initialize(client ...IAPGoogleClient) (err error) {
	if len(client) > 0 && client[0] != nil {
		this.client = client[0]
	} else {
		if this.client, err = playstore.New([]byte(this.config.Key)); err != nil {
			return fmt.Errorf("iapGoogle initialize: %w", err)
		} // if
	} // if

	this.verify = make(chan *iapGoogle, capacity)
	this.signal.Add(1)
	this.ctx, this.cancel = context.WithCancel(context.Background())
	go this.execute(this.verify)
	return nil
}

// Finalize 結束處理
func (this *IAPGoogle) Finalize() {
	if this.cancel != nil {
		this.cancel()
	} // if

	this.signal.Wait()
	this.verify = nil
}

// Verify 驗證憑證
//   - productID: Google 內的產品編號
//   - receipt: 購買憑證
func (this *IAPGoogle) Verify(productID, receipt string) IAPGoogleResult {
	if this.verify == nil {
		return this.fail(fmt.Errorf("iapGoogle verify: close"))
	} // if

	result := &iapGoogle{
		productID: productID,
		receipt:   receipt,
		result:    make(chan IAPGoogleResult, 1),
	}

	// 嘗試送出驗證請求, 若通道塞滿則等待直到逾時
	select {
	case <-this.ctx.Done():
		return this.fail(fmt.Errorf("iapGoogle verify: shutdown"))

	case <-time.After(timeoutMax):
		return this.fail(fmt.Errorf("iapGoogle verify: timeout send"))

	case this.verify <- result:
	} // select

	// 等待驗證結果回傳, 若逾時則失敗
	select {
	case <-time.After(timeoutMax):
		return this.fail(fmt.Errorf("iapGoogle verify: timeout recv"))

	case r := <-result.result:
		return r
	} // select
}

// execute 執行驗證
func (this *IAPGoogle) execute(verify chan *iapGoogle) {
	defer this.signal.Done()

	for {
		select {
		case <-this.ctx.Done():
			return

		case itor := <-verify:
			time.Sleep(interval) // 由於驗證 API 有速率限制, 所以需要等待後才能繼續下一個驗證
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			result, err := this.client.VerifyProduct(ctx, this.config.Bundle, itor.productID, itor.receipt)
			cancel() // 避免cancel洩漏

			if err != nil {
				channelTry(itor.result, this.fail(fmt.Errorf("iapGoogle execute: %w", err)))
				continue
			} // if

			channelTry(itor.result, this.succ(result.PurchaseTimeMillis))
		} // select
	} // for
}

// succ 建立成功的驗證結果
func (this *IAPGoogle) succ(millisecond int64) IAPGoogleResult {
	return IAPGoogleResult{
		Time: time.Unix(
			millisecond/1000, //nolint:mnd
			(millisecond%1000)*int64(time.Millisecond),
		),
	}
}

// fail 建立失敗的驗證結果
func (this *IAPGoogle) fail(err error) IAPGoogleResult {
	return IAPGoogleResult{
		Err: err,
	}
}

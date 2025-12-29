package iaps

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/awa/go-iap/appstore/api"
)

// NewIAPApple 建立 Apple IAP 驗證器
func NewIAPApple(config *IAPAppleConfig) *IAPApple {
	return &IAPApple{
		config: config,
	}
}

// IAPApple Apple IAP 驗證器
type IAPApple struct {
	config *IAPAppleConfig    // 驗證設定
	client IAPAppleClient     // 驗證客戶端
	verify chan *iapApple     // 驗證通道
	signal sync.WaitGroup     // 通知信號
	ctx    context.Context    // 關閉物件
	cancel context.CancelFunc // 關閉函式
}

// IAPAppleConfig Apple IAP 驗證設定資料
type IAPAppleConfig struct {
	Key     string `yaml:"key"`     // 密鑰字串
	KeyID   string `yaml:"keyID"`   // 密鑰ID
	Bundle  string `yaml:"bundle"`  // 軟體包名稱(BundleID)
	Issuer  string `yaml:"issuer"`  // 發行人名稱
	Sandbox bool   `yaml:"sandbox"` // 是否使用沙盒環境
}

// IAPAppleClient Apple IAP 驗證客戶端介面
type IAPAppleClient interface {
	GetTransactionInfo(context.Context, string) (*api.TransactionInfoResponse, error)
	ParseSignedTransaction(string) (*api.JWSTransaction, error)
}

// iapApple Apple IAP 驗證資料
type iapApple struct {
	productID string         // 產品編號
	receipt   string         // 購買憑證(對應 Apple TransactionID)
	retry     int            // 重試次數
	retryErr  error          // 重試錯誤
	result    chan IAPResult // 驗證結果通道
}

// Initialize 初始化處理
func (this *IAPApple) Initialize(client ...IAPAppleClient) (err error) {
	if len(client) > 0 && client[0] != nil {
		this.client = client[0]
	} else {
		this.client = api.NewStoreClient(&api.StoreConfig{
			KeyContent: []byte(this.config.Key),
			KeyID:      this.config.KeyID,
			BundleID:   this.config.Bundle,
			Issuer:     this.config.Issuer,
			Sandbox:    this.config.Sandbox,
		})
	} // if

	this.verify = make(chan *iapApple, capacity)
	this.signal.Add(1)
	this.ctx, this.cancel = context.WithCancel(context.Background())
	go this.execute(this.verify)
	return nil
}

// Finalize 結束處理
func (this *IAPApple) Finalize() {
	if this.cancel != nil {
		this.cancel()
	} // if

	this.signal.Wait()
	this.verify = nil
}

// Verify 驗證憑證
//   - productID: Apple 內的產品編號
//   - receipt: 購買憑證(TransactionID)
func (this *IAPApple) Verify(productID, receipt string) IAPResult {
	if this.verify == nil {
		return this.fail(fmt.Errorf("iapApple verify: close"))
	} // if

	result := &iapApple{
		productID: productID,
		receipt:   receipt,
		result:    make(chan IAPResult, 1),
	}

	// 嘗試送出驗證請求, 若通道塞滿則等待直到逾時
	select {
	case <-this.ctx.Done():
		return this.fail(fmt.Errorf("iapApple verify: shutdown"))

	case <-time.After(timeoutMax):
		return this.fail(fmt.Errorf("iapApple verify: timeout send"))

	case this.verify <- result:
	} // select

	// 等待驗證結果回傳, 若逾時則失敗
	select {
	case <-time.After(timeoutMax):
		return this.fail(fmt.Errorf("iapApple verify: timeout recv"))

	case r := <-result.result:
		return r
	} // select
}

// execute 執行驗證
func (this *IAPApple) execute(verify chan *iapApple) {
	defer this.signal.Done()

	for {
		select {
		case <-this.ctx.Done():
			return

		case itor := <-verify:
			time.Sleep(interval) // 由於驗證 API 有速率限制, 所以需要等待後才能繼續下一個驗證
			ctx, cancel := context.WithTimeout(this.ctx, timeout)
			respond, err := this.client.GetTransactionInfo(ctx, itor.receipt)
			cancel() // 避免 cancel 洩漏

			// 驗證 API 有時會不明原因錯誤, 這裡採用重試策略
			if err != nil {
				itor.retry++
				itor.retryErr = fmt.Errorf("iapApple execute: %w", err)

				if itor.retry >= retry { // 超過最大重試次數則結束, 否則繼續重試
					channelTry(itor.result, this.fail(itor.retryErr))
				} else {
					select {
					case <-this.ctx.Done():
						channelTry(itor.result, this.fail(fmt.Errorf("iapApple execute: shutdown")))

					case verify <- itor:
					default:
						channelTry(itor.result, this.fail(fmt.Errorf("iapApple execute: queue full")))
					} // select
				} // if

				continue
			} // if

			result, err := this.client.ParseSignedTransaction(respond.SignedTransactionInfo)

			if err != nil {
				channelTry(itor.result, this.fail(fmt.Errorf("iapApple execute: %w", err)))
				continue
			} // if

			if result.ProductID != itor.productID {
				channelTry(itor.result, this.fail(fmt.Errorf("iapApple execute: productID")))
				continue
			} // if

			if result.TransactionID != itor.receipt {
				channelTry(itor.result, this.fail(fmt.Errorf("iapApple execute: receipt")))
				continue
			} // if

			channelTry(itor.result, this.succ(result.PurchaseDate))
		} // select
	} // for
}

// succ 建立成功的驗證結果
func (this *IAPApple) succ(millisecond int64) IAPResult {
	return IAPResult{
		Time: time.Unix(
			millisecond/1000, //nolint:mnd
			(millisecond%1000)*int64(time.Millisecond),
		),
	}
}

// fail 建立失敗的驗證結果
func (this *IAPApple) fail(err error) IAPResult {
	return IAPResult{
		Err: err,
	}
}

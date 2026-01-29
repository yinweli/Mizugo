package iaps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
)

// NewIAPOneStore 建立 OneStore IAP 驗證器
func NewIAPOneStore(config *IAPOneStoreConfig) *IAPOneStore {
	return &IAPOneStore{
		config: config,
	}
}

// IAPOneStore OneStore IAP 驗證器
type IAPOneStore struct {
	config      *IAPOneStoreConfig // 驗證設定
	client      IAPOneStoreClient  // 驗證客戶端
	verify      chan *iapOneStore  // 驗證通道
	signal      sync.WaitGroup     // 通知信號
	ctx         context.Context    // 關閉物件
	cancel      context.CancelFunc // 關閉函式
	token       string             // 權杖字串
	tokenExpire time.Time          // 權杖逾期時間
}

// IAPOneStoreConfig OneStore IAP 驗證設定資料
type IAPOneStoreConfig struct {
	Global       bool   `yaml:"global"`       // true表示全球市場, false表示僅限韓國
	ClientID     string `yaml:"clientID"`     // 用戶端 ID, 同時也就是 PackageName
	ClientSecret string `yaml:"clientSecret"` // 用戶端密鑰
	Sandbox      bool   `yaml:"sandbox"`      // 是否使用沙盒環境
}

// IAPOneStoreClient OneStore IAP 驗證客戶端介面
type IAPOneStoreClient interface {
	Do(*http.Request) (*http.Response, error)
}

// iapOneStore OneStore IAP 驗證資料
type iapOneStore struct {
	productID string         // 產品編號
	receipt   string         // 購買憑證
	retry     int            // 重試次數
	retryErr  error          // 重試錯誤
	result    chan IAPResult // 驗證結果通道
}

// iapOneStoreToken OneStore IAP 獲取權杖資料
type iapOneStoreToken struct {
	Token  string `json:"access_token"` // 權杖字串
	Expire int    `json:"expires_in"`   // 權杖有效時間(秒)
}

// iapOneStoreVerify OneStore IAP 驗證結果資料
type iapOneStoreVerify struct {
	DeveloperPayload string `json:"developerPayload"` // 開發公司提供的支付固有標示符
	PurchaseID       string `json:"purchaseId"`       // 購買 ID
	PurchaseTime     int64  `json:"purchaseTime"`     // 購買時間(毫秒)
	AcknowledgeState int    `json:"acknowledgeState"` // 確認狀態(0: 未確認, 1: 確認)
	ConsumptionState int    `json:"consumptionState"` // 使用狀態(0: 未使用, 1: 使用)
	PurchaseState    int    `json:"purchaseState"`    // 購買狀態(0: 購買完成, 1: 取消完成)
	Quantity         int    `json:"quantity"`         // 購買數量
}

// iapOneStoreError OneStore IAP 錯誤訊息資料
type iapOneStoreError struct {
	Error struct {
		Code    string `json:"code"`    // 錯誤編號
		Message string `json:"message"` // 錯誤訊息
	} `json:"error"`
}

// Initialize 初始化處理
func (this *IAPOneStore) Initialize(client ...IAPOneStoreClient) error {
	if len(client) > 0 && client[0] != nil {
		this.client = client[0]
	} else {
		this.client = &http.Client{}
	} // if

	this.verify = make(chan *iapOneStore, capacity)
	this.signal.Add(1)
	this.ctx, this.cancel = context.WithCancel(context.Background())
	go this.execute(this.verify)
	return nil
}

// Finalize 結束處理
func (this *IAPOneStore) Finalize() {
	if this.cancel != nil {
		this.cancel()
	} // if

	this.signal.Wait()
	this.verify = nil
	this.token = ""
	this.tokenExpire = time.Time{}
}

// Verify 驗證憑證
//   - productID: OneStore 內的產品編號
//   - receipt: 購買憑證
func (this *IAPOneStore) Verify(productID, receipt string) IAPResult {
	if this.verify == nil {
		return this.fail(fmt.Errorf("iapOneStore verify: close"))
	} // if

	result := &iapOneStore{
		productID: productID,
		receipt:   receipt,
		result:    make(chan IAPResult, 1),
	}

	// 嘗試送出驗證請求, 若通道塞滿則等待直到逾時
	select {
	case <-this.ctx.Done():
		return this.fail(fmt.Errorf("iapOneStore verify: shutdown"))

	case <-time.After(timeoutMax):
		return this.fail(fmt.Errorf("iapOneStore verify: timeout send"))

	case this.verify <- result:
	} // select

	// 等待驗證結果回傳, 若逾時則失敗
	select {
	case <-time.After(timeoutMax):
		return this.fail(fmt.Errorf("iapOneStore verify: timeout recv"))

	case r := <-result.result:
		return r
	} // select
}

// Client 取得驗證客戶端
func (this *IAPOneStore) Client() IAPOneStoreClient {
	return this.client
}

// execute 執行驗證
func (this *IAPOneStore) execute(verify chan *iapOneStore) {
	defer this.signal.Done()

	for {
		select {
		case <-this.ctx.Done():
			return

		case itor := <-verify:
			time.Sleep(interval) // 由於驗證 API 有速率限制, 所以需要等待後才能繼續下一個驗證
			ctx, cancel := context.WithTimeout(this.ctx, timeout)
			err := this.doToken(ctx)
			cancel() // 避免 cancel 洩漏

			if err != nil {
				itor.retry++
				itor.retryErr = fmt.Errorf("iapOneStore execute: %w", err)

				if itor.retry >= retry { // 超過最大重試次數則結束, 否則繼續重試
					channelTry(itor.result, this.fail(itor.retryErr))
				} else {
					select {
					case <-this.ctx.Done():
						channelTry(itor.result, this.fail(fmt.Errorf("iapOneStore execute: shutdown")))

					case verify <- itor:
					default:
						channelTry(itor.result, this.fail(fmt.Errorf("iapOneStore execute: queue full")))
					} // select
				} // if

				continue
			} // if

			ctx, cancel = context.WithTimeout(this.ctx, timeout)
			result := this.doPurchase(ctx, itor.productID, itor.receipt)
			cancel() // 避免 cancel 洩漏
			channelTry(itor.result, result)
		} // select
	} // for
}

// doToken 執行獲取權杖
func (this *IAPOneStore) doToken(ctx context.Context) error {
	now := helps.Time()

	if this.token != "" && this.tokenExpire.After(now) {
		return nil
	} // if

	api := this.getAPI(this.config.Sandbox) + "/v7/oauth/token"
	post := url.Values{}
	post.Set("grant_type", "client_credentials")
	post.Set("client_id", this.config.ClientID)
	post.Set("client_secret", this.config.ClientSecret)
	request, err := http.NewRequestWithContext(ctx, "POST", api, bytes.NewBufferString(post.Encode()))

	if err != nil {
		return fmt.Errorf("iapOneStore token: %w", err)
	} // if

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("x-market-code", this.getMKT(this.config.Global))
	respond, err := this.client.Do(request)

	if err != nil {
		return fmt.Errorf("iapOneStore token: %w", err)
	} // if

	defer func() {
		_ = respond.Body.Close()
	}()

	body, err := io.ReadAll(respond.Body)

	if err != nil {
		return fmt.Errorf("iapOneStore token: %w", err)
	} // if

	if respond.StatusCode != http.StatusOK {
		result := &iapOneStoreError{}

		if err = json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("iapOneStore token: %w", err)
		} // if

		this.token = ""
		this.tokenExpire = now
		return fmt.Errorf("iapOneStore token: [%v] %v", result.Error.Code, result.Error.Message)
	} // if

	result := &iapOneStoreToken{}

	if err = json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("iapOneStore token: %w", err)
	} // if

	this.token = result.Token
	this.tokenExpire = this.getExpire(now, result.Expire)
	return nil
}

// doPurchase 執行驗證查詢
func (this *IAPOneStore) doPurchase(ctx context.Context, productID, receipt string) IAPResult {
	api := this.getAPI(this.config.Sandbox) + fmt.Sprintf("/v7/apps/%v/purchases/inapp/products/%v/%v", this.config.ClientID, productID, receipt)
	request, err := http.NewRequestWithContext(ctx, "GET", api, http.NoBody)

	if err != nil {
		return this.fail(fmt.Errorf("iapOneStore verify: %w", err))
	} // if

	request.Header.Set("Authorization", "Bearer "+this.token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-market-code", this.getMKT(this.config.Global))
	respond, err := this.client.Do(request)

	if err != nil {
		return this.fail(fmt.Errorf("iapOneStore verify: %w", err))
	} // if

	defer func() {
		_ = respond.Body.Close()
	}()

	body, err := io.ReadAll(respond.Body)

	if err != nil {
		return this.fail(fmt.Errorf("iapOneStore verify: %w", err))
	} // if

	if respond.StatusCode != http.StatusOK {
		result := &iapOneStoreError{}

		if err = json.Unmarshal(body, result); err != nil {
			return this.fail(fmt.Errorf("iapOneStore verify: %w", err))
		} // if

		return this.fail(fmt.Errorf("iapOneStore verify: [%v] %v", result.Error.Code, result.Error.Message))
	} // if

	result := &iapOneStoreVerify{}

	if err = json.Unmarshal(body, result); err != nil {
		return this.fail(fmt.Errorf("iapOneStore verify: %w", err))
	} // if

	if result.PurchaseState != 0 {
		return this.fail(fmt.Errorf("iapOneStore verify: unpurchased"))
	} // if

	return this.succ(result.PurchaseTime)
}

// getAPI 取得 API 網址
func (this *IAPOneStore) getAPI(sandbox bool) string {
	if sandbox {
		return "https://sbpp.onestore.co.kr" // 開發環境
	} // if

	return "https://iap-apis.onestore.net" // 正式環境
}

// getMKT 取得市場分類代碼
func (this *IAPOneStore) getMKT(global bool) string {
	if global {
		return "MKT_GLB" // 全球市場, 伺服器API回應提供的時間標準為 UTC+0
	} // if

	return "MKT_ONE" // 韓國市場, 伺服器API回應提供的時間標準為 UTC+9
}

// getExpire 取得權杖逾期時間
func (this *IAPOneStore) getExpire(now time.Time, expire int) time.Time {
	const factor = 0.85 // 故意讓權杖預期時間短一些, 方便更新權杖
	return now.Add(time.Duration(float64(expire)*factor) * helps.TimeSecond)
}

// succ 建立成功的驗證結果
func (this *IAPOneStore) succ(millisecond int64) IAPResult {
	return IAPResult{
		Time: time.Unix(
			millisecond/1000, //nolint:mnd
			(millisecond%1000)*int64(time.Millisecond),
		),
	}
}

// fail 建立失敗的驗證結果
func (this *IAPOneStore) fail(err error) IAPResult {
	return IAPResult{
		Err: err,
	}
}

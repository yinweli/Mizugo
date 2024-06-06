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

	"github.com/yinweli/Mizugo/mizugos/helps"
)

// NewIAPOneStore 建立OneStore驗證器
func NewIAPOneStore(config *IAPOneStoreConfig) *IAPOneStore {
	return &IAPOneStore{
		config: config,
	}
}

// IAPOneStore OneStore驗證器
type IAPOneStore struct {
	config      *IAPOneStoreConfig // 驗證設定
	client      *http.Client       // 驗證客戶端
	verify      chan *iapOneStore  // 驗證通道
	signal      sync.WaitGroup     // 通知信號
	token       string             // 權杖字串
	tokenExpire time.Time          // 權杖逾期時間
}

// IAPOneStoreConfig OneStore驗證設定資料
type IAPOneStoreConfig struct {
	Global       bool          `yaml:"global"`       // true表示全球市場, false表示僅限韓國
	ClientID     string        `yaml:"clientID"`     // ClientID, 同時也就是 PackageName
	ClientSecret string        `yaml:"clientSecret"` // ClientSecret
	Sandbox      bool          `yaml:"sandbox"`      // 沙盒旗標
	WaitTime     time.Duration `yaml:"waitTime"`     // 等待時間
	Capacity     int           `yaml:"capacity"`     // 通道容量
	Retry        int           `yaml:"retry"`        // 重試次數
}

// iapOneStore OneStore驗證資料
type iapOneStore struct {
	productID   string     // 產品編號
	certificate string     // 購買憑證
	retry       int        // 重試次數
	retryErr    error      // 重試錯誤
	result      chan error // 結果通道
}

// iapOneStoreError OneStore錯誤訊息資料
type iapOneStoreError struct {
	Error struct {
		Code    string `json:"code"`    // 錯誤編號
		Message string `json:"message"` // 錯誤訊息
	} `json:"error"`
}

// iapOneStoreToken OneStore獲取權杖資料
type iapOneStoreToken struct {
	Token  string `json:"access_token"` // 權杖字串
	Expire int    `json:"expires_in"`   // 權杖有效時間(秒)
}

// iapOneStoreVerify OneStore驗證結果資料
type iapOneStoreVerify struct {
	DeveloperPayload string `json:"developerPayload"` // 開發公司提供的支付固有標示符
	PurchaseID       string `json:"purchaseId"`       // 購買ID
	PurchaseTime     int64  `json:"purchaseTime"`     // 購買時間(毫秒)
	AcknowledgeState int    `json:"acknowledgeState"` // 確認狀態(0: 未確認, 1: 確認)
	ConsumptionState int    `json:"consumptionState"` // 使用狀態(0: 未使用, 1: 使用)
	PurchaseState    int    `json:"purchaseState"`    // 購買狀態(0: 購買完成, 1: 取消完成)
	Quantity         int    `json:"quantity"`         // 購買數量
}

// Initialize 初始化處理
func (this *IAPOneStore) Initialize() error {
	this.client = &http.Client{}
	this.verify = make(chan *iapOneStore, this.config.Capacity+1) // 避免使用者將通道容量設為0導致卡住
	this.signal.Add(1)
	go this.execute(this.verify)
	return nil
}

// Finalize 結束處理
func (this *IAPOneStore) Finalize() {
	close(this.verify)
	this.verify = nil
	this.signal.Wait()
}

// Verify 驗證憑證
func (this *IAPOneStore) Verify(productID, certificate string) error {
	if this.verify == nil {
		return fmt.Errorf("iapOneStore verify: close")
	} // if

	result := &iapOneStore{
		productID:   productID,
		certificate: certificate,
		result:      make(chan error),
	}
	this.verify <- result
	return <-result.result
}

// execute 執行驗證
func (this *IAPOneStore) execute(verify chan *iapOneStore) {
	for itor := range verify {
		// 由於驗證api有速率限制, 所以需要等待後才能繼續下一個驗證
		time.Sleep(this.config.WaitTime)

		// 如果重試超過限制, 還是只能當作錯誤
		if itor.retry > 0 && itor.retry >= this.config.Retry {
			itor.result <- itor.retryErr
			continue
		} // if

		if err := this.executeToken(); err != nil {
			itor.retry++
			itor.retryErr = fmt.Errorf("iapOneStore execute: %w", err)
			verify <- itor
			continue
		} // if

		if err := this.executePurchase(itor.productID, itor.certificate); err != nil {
			itor.result <- fmt.Errorf("iapOneStore execute: %w", err)
			continue
		} // if

		itor.result <- nil
	} // for

	this.signal.Done()
}

// executeToken 執行獲取權杖
func (this *IAPOneStore) executeToken() error {
	now := helps.Time()

	if this.token != "" && this.tokenExpire.After(now) {
		return nil
	} // if

	api := this.getAPI(this.config.Sandbox) + "/v7/oauth/token"
	post := url.Values{}
	post.Set("grant_type", "client_credentials")
	post.Set("client_id", this.config.ClientID)
	post.Set("client_secret", this.config.ClientSecret)
	request, err := http.NewRequestWithContext(context.Background(), "POST", api, bytes.NewBufferString(post.Encode()))

	if err != nil {
		return fmt.Errorf("token: %w", err)
	} // if

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("x-market-code", this.getMKT(this.config.Global))
	respond, err := this.client.Do(request)

	if err != nil {
		return fmt.Errorf("token: %w", err)
	} // if

	defer func() {
		_ = respond.Body.Close()
	}()

	body, err := io.ReadAll(respond.Body)

	if err != nil {
		return fmt.Errorf("token: %w", err)
	} // if

	if respond.StatusCode == http.StatusOK {
		result := &iapOneStoreToken{}

		if err = json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("token: %w", err)
		} // if

		this.token = result.Token
		this.tokenExpire = this.getExpire(now, result.Expire)
		return nil
	} else {
		result := &iapOneStoreError{}

		if err = json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("token: %w", err)
		} // if

		this.token = ""
		this.tokenExpire = now
		return fmt.Errorf("token: [%v] %v", result.Error.Code, result.Error.Message)
	} // if
}

// executePurchase 執行驗證查詢
func (this *IAPOneStore) executePurchase(productID, certificate string) error {
	api := this.getAPI(this.config.Sandbox) + fmt.Sprintf("/v7/apps/%v/purchases/inapp/products/%v/%v", this.config.ClientID, productID, certificate)
	request, err := http.NewRequestWithContext(context.Background(), "GET", api, http.NoBody)

	if err != nil {
		return fmt.Errorf("purchase: %w", err)
	} // if

	request.Header.Set("Authorization", "Bearer "+this.token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-market-code", this.getMKT(this.config.Global))
	respond, err := this.client.Do(request)

	if err != nil {
		return fmt.Errorf("purchase: %w", err)
	} // if

	defer func() {
		_ = respond.Body.Close()
	}()

	body, err := io.ReadAll(respond.Body)

	if err != nil {
		return fmt.Errorf("purchase: %w", err)
	} // if

	if respond.StatusCode == http.StatusOK {
		result := &iapOneStoreVerify{}

		if err = json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("purchase: %w", err)
		} // if

		if result.PurchaseState != 0 {
			return fmt.Errorf("purchase: unpurchased")
		} // if

		return nil
	} else {
		result := &iapOneStoreError{}

		if err = json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("purchase: %w", err)
		} // if

		return fmt.Errorf("purchase: [%v] %v", result.Error.Code, result.Error.Message)
	} // if
}

// getAPI 取得API網址
func (this *IAPOneStore) getAPI(sandbox bool) string {
	if sandbox {
		return "https://sbpp.onestore.co.kr" // 開發環境
	} else {
		return "https://iap-apis.onestore.net" // 正式環境
	} // if
}

// getMKT 取得市場分類代碼
func (this *IAPOneStore) getMKT(global bool) string {
	if global {
		return "MKT_GLB" // 全球市場, 伺服器API回應提供的時間標準為 UTC+0
	} else {
		return "MKT_ONE" // 韓國市場, 伺服器API回應提供的時間標準為 UTC+9
	} // if
}

// getExpire 取得權杖逾期時間
func (this *IAPOneStore) getExpire(now time.Time, expire int) time.Time {
	const factor = 0.85 // 故意讓權杖預期時間短一些, 方便更新權杖
	return now.Add(time.Duration(float64(expire)*factor) * helps.TimeSecond)
}

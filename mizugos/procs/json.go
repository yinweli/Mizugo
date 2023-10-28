package procs

import (
	"encoding/json"
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/cryptos"
	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/msgs"
)

// NewJson 建立json處理器
func NewJson() *Json {
	return &Json{
		Procmgr: NewProcmgr(),
	}
}

// Json json處理器, 封包結構使用JsonMsg, 可以選擇是否啟用base64編碼或是des-cbc加密
//   - 訊息定義: support/proto/mizugo/msg-go/msgs-json/jsonmsg.go
//   - 訊息定義: support/proto/mizugo/msg-cs/msgs-json/Jsonmsg.cs
//   - 封包編碼: json編碼成位元陣列, (可選)des-cbc加密, (可選)base64編碼
//   - 封包解碼: (可選)base64解碼, (可選)des-cbc解密, json解碼成訊息結構
type Json struct {
	*Procmgr        // 管理器
	base64   bool   // 是否啟用base64
	desCBC   bool   // 是否啟用des-cbc加密
	desKey   []byte // des密鑰
	desIV    []byte // des初始向量
}

// Encode 封包編碼
func (this *Json) Encode(input any) (output []byte, err error) {
	message, err := helps.CastPointer[msgs.JsonMsg](input)

	if err != nil {
		return nil, fmt.Errorf("json encode: %w", err)
	} // if

	output, err = json.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("json encode: %w", err)
	} // if

	if this.desCBC {
		if output, err = cryptos.DesCBCEncrypt(cryptos.PaddingPKCS7, this.desKey, this.desIV, output); err != nil {
			return nil, fmt.Errorf("json encode: %w", err)
		} // if
	} // if

	if this.base64 {
		output = cryptos.Base64Encode(output)
	} // if

	return output, nil
}

// Decode 封包解碼
func (this *Json) Decode(input []byte) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("json decode: input nil")
	} // if

	if this.base64 {
		input, err = cryptos.Base64Decode(input)

		if err != nil {
			return nil, fmt.Errorf("json decode: %w", err)
		} // if
	} // if

	if this.desCBC {
		input, err = cryptos.DesCBCDecrypt(cryptos.PaddingPKCS7, this.desKey, this.desIV, input)

		if err != nil {
			return nil, fmt.Errorf("json decode: %w", err)
		} // if
	} // if

	message := &msgs.JsonMsg{}

	if err = json.Unmarshal(input, message); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *Json) Process(input any) error {
	message, err := helps.CastPointer[msgs.JsonMsg](input)

	if err != nil {
		return fmt.Errorf("json process: %w", err)
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("json process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// Base64 設定是否啟用base64
func (this *Json) Base64(enable bool) *Json {
	this.base64 = enable
	return this
}

// DesCBC 是否啟用des-cbc加密
func (this *Json) DesCBC(enable bool, key, iv string) *Json {
	this.desCBC = enable
	this.desKey = []byte(key)
	this.desIV = []byte(iv)
	return this
}

// JsonMarshal json訊息序列化
func JsonMarshal(messageID MessageID, input any) (output *msgs.JsonMsg, err error) {
	if input == nil {
		return nil, fmt.Errorf("json marshal: input nil")
	} // if

	message, err := json.Marshal(input)

	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	} // if

	return &msgs.JsonMsg{
		MessageID: messageID,
		Message:   message,
	}, nil
}

// JsonUnmarshal json訊息反序列化
func JsonUnmarshal[T any](input any) (messageID MessageID, output *T, err error) {
	if input == nil {
		return 0, nil, fmt.Errorf("json unmarshal: input nil")
	} // if

	message, err := helps.CastPointer[msgs.JsonMsg](input)

	if err != nil {
		return 0, nil, fmt.Errorf("json unmarshal: %w", err)
	} // if

	output = new(T)

	if err = json.Unmarshal(message.Message, output); err != nil {
		return 0, nil, fmt.Errorf("json unmarshal: %w", err)
	} // if

	return message.MessageID, output, nil
}

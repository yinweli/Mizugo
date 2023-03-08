package procs

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/cryptos"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/utils"
)

// NewProto 建立proto處理器
func NewProto() *Proto {
	return &Proto{
		Procmgr: NewProcmgr(),
	}
}

// Proto proto處理器, 封包結構使用ProtoMsg, 可以選擇是否啟用base64編碼或是des-cbc加密
//   - 訊息定義: support/proto/mizugo/protomsg.proto
//   - 封包編碼: protobuf編碼成位元陣列, (可選)des-cbc加密, (可選)base64編碼
//   - 封包解碼: (可選)base64解碼, (可選)des-cbc解密, protobuf解碼成訊息結構
type Proto struct {
	*Procmgr        // 管理器
	base64   bool   // 是否啟用base64
	desCBC   bool   // 是否啟用des-cbc加密
	desKey   []byte // des密鑰
	desIV    []byte // des初始向量
}

// Encode 封包編碼
func (this *Proto) Encode(input any) (output []byte, err error) {
	message, err := utils.CastPointer[msgs.ProtoMsg](input)

	if err != nil {
		return nil, fmt.Errorf("proto encode: %w", err)
	} // if

	output, err = proto.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("proto encode: %w", err)
	} // if

	if this.desCBC {
		if output, err = cryptos.DesCBCEncrypt(cryptos.PaddingPKCS7, this.desKey, this.desIV, output); err != nil {
			return nil, fmt.Errorf("proto encode: %w", err)
		} // if
	} // if

	if this.base64 {
		output = cryptos.Base64Encode(output)
	} // if

	return output, nil
}

// Decode 封包解碼
func (this *Proto) Decode(input []byte) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("proto decode: input nil")
	} // if

	if this.base64 {
		input, err = cryptos.Base64Decode(input)

		if err != nil {
			return nil, fmt.Errorf("proto decode: %w", err)
		} // if
	} // if

	if this.desCBC {
		input, err = cryptos.DesCBCDecrypt(cryptos.PaddingPKCS7, this.desKey, this.desIV, input)

		if err != nil {
			return nil, fmt.Errorf("proto decode: %w", err)
		} // if
	} // if

	message := &msgs.ProtoMsg{}

	if err = proto.Unmarshal(input, message); err != nil {
		return nil, fmt.Errorf("proto decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *Proto) Process(input any) error {
	message, err := utils.CastPointer[msgs.ProtoMsg](input)

	if err != nil {
		return fmt.Errorf("proto process: %w", err)
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("proto process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// Base64 設定是否啟用base64
func (this *Proto) Base64(enable bool) *Proto {
	this.base64 = enable
	return this
}

// DesCBC 是否啟用des-cbc加密
func (this *Proto) DesCBC(enable bool, key, iv string) *Proto {
	this.desCBC = enable
	this.desKey = []byte(key)
	this.desIV = []byte(iv)
	return this
}

// ProtoMarshal 序列化proto訊息
func ProtoMarshal(messageID MessageID, input proto.Message) (output *msgs.ProtoMsg, err error) {
	if input == nil {
		return nil, fmt.Errorf("proto marshal: input nil")
	} // if

	message, err := anypb.New(input)

	if err != nil {
		return nil, fmt.Errorf("proto marshal: %w", err)
	} // if

	return &msgs.ProtoMsg{
		MessageID: messageID,
		Message:   message,
	}, nil
}

// ProtoUnmarshal 反序列化proto訊息
func ProtoUnmarshal[T any](input any) (messageID MessageID, output *T, err error) {
	if input == nil {
		return 0, nil, fmt.Errorf("proto unmarshal: input nil")
	} // if

	message, err := utils.CastPointer[msgs.ProtoMsg](input)

	if err != nil {
		return 0, nil, fmt.Errorf("proto unmarshal: %w", err)
	} // if

	temp, err := message.Message.UnmarshalNew()

	if err != nil {
		return 0, nil, fmt.Errorf("proto unmarshal: %w", err)
	} // if

	output, ok := temp.(any).(*T)

	if ok == false {
		return 0, nil, fmt.Errorf("proto unmarshal: cast failed")
	} // if

	return message.MessageID, output, nil
}

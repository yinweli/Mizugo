using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using System;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    /// <summary>
    /// plist處理器, 封包結構使用PListMsg
    /// 由於使用到des加密, 安全性較高, 適合用來傳送一般封包, 使用時需要設定密鑰
    /// 由於採用複數訊息設計, 因此封包內可以填入多個訊息來跟伺服器溝通(json/proto處理器則使用訊息結構與伺服器溝通)
    /// 訊息內容: support/proto/mizugo/plistmsg.proto
    /// 封包編碼: protobuf編碼成位元陣列, 再通過des加密
    /// 封包解碼: des解密, 再通過protobuf解碼成訊息結構
    /// </summary>
    public class PListProc : Procmgr
    {
        public override byte[] Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not PListMsg message)
                throw new InvalidMessageException("encode");

            var protoBytes = message.ToByteArray();
        }

        public override object Decode(byte[] input)
        {
            throw new NotImplementedException();
        }

        public override void Process(object input)
        {
            throw new NotImplementedException();
        }

        /// <summary>
        /// proto訊息序列化
        /// </summary>
        /// <param name="messageID">訊息編號</param>
        /// <param name="message">訊息物件</param>
        /// <returns>訊息物件</returns>
        public static ProtoMsg Marshal(MessageID messageID, IMessage message)
        {
            if (message == null)
                throw new ArgumentNullException("input");

            return new ProtoMsg { MessageID = messageID, Message = Any.Pack(message) };
        }

        /// <summary>
        /// proto訊息反序列化
        /// </summary>
        /// <typeparam name="T">訊息類型</typeparam>
        /// <param name="input">輸入物件</param>
        /// <param name="messageID">訊息編號</param>
        /// <param name="message">訊息物件</param>
        public static void Unmarshal<T>(object input, out MessageID messageID, out T message)
            where T : IMessage, new()
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not ProtoMsg protoMsg)
                throw new InvalidMessageException("unmarshal");

            messageID = protoMsg.MessageID;
            message = protoMsg.Message.Unpack<T>();
        }
    }
}

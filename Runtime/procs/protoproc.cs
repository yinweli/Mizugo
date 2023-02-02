using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using System;
using System.Text;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    /// <summary>
    /// proto處理器, 封包結構使用ProtoMsg
    /// 沒有使用加密技術, 所以安全性很低, 僅用於傳送簡單訊息或是傳送密鑰使用
    /// 訊息內容: support/proto/mizugo/protomsg.proto
    /// 封包編碼: protobuf編碼成位元陣列, 再通過base64編碼
    /// 封包解碼: base64解碼, 再通過protobuf解碼成訊息結構
    /// </summary>
    public class ProtoProc : Procmgr
    {
        public override byte[] Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not ProtoMsg message)
                throw new InvalidMessageException("encode");

            var protoBytes = message.ToByteArray();
            var base64 = Convert.ToBase64String(protoBytes);
            var base64Bytes = Encoding.UTF8.GetBytes(base64);

            return base64Bytes;
        }

        public override object Decode(byte[] input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            var base64 = Encoding.UTF8.GetString(input);
            var base64Bytes = Convert.FromBase64String(base64);
            var message = ProtoMsg.Parser.ParseFrom(base64Bytes);

            return message;
        }

        public override void Process(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not ProtoMsg message)
                throw new InvalidMessageException("process");

            var process = Get(message.MessageID);

            if (process == null)
                throw new UnProcessException(message.MessageID);

            process(message);
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

using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using System;
using System.Security.Cryptography;
using System.Text;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    /// <summary>
    /// plist處理器, 封包結構使用PListMsg
    /// 採用des-cbc加密, 安全性較高, 適合用來傳送一般封包, 使用時需要設定密鑰以及初始向量
    /// 採用複數訊息設計, 因此封包內可以填入多個訊息來跟伺服器溝通
    /// 封包結構: support/proto/mizugo/plistmsg.proto
    /// 封包編碼: protobuf編碼成位元陣列, 再通過des加密
    /// 封包解碼: des解密, 再通過protobuf解碼
    /// </summary>
    public partial class PListProc : Procmgr
    {
        /// <summary>
        /// plist使用的填充模式
        /// </summary>
        private const PaddingMode padding = PaddingMode.PKCS7;

        public override byte[] Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not PListMsg message)
                throw new InvalidMessageException("encode");

            var protoBytes = message.ToByteArray();
            var encode = DesCBC.Encrypt(padding, key, iv, protoBytes);
            return encode;
        }

        public override object Decode(byte[] input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            var decode = DesCBC.Decrypt(padding, key, iv, input);
            var message = PListMsg.Parser.ParseFrom(decode);
            return message;
        }

        public override void Process(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not PListMsg message)
                throw new InvalidMessageException("process");

            var process = Get(message.MessageID);

            if (process == null)
                throw new UnprocessException(message.MessageID);

            process(message);
        }

        /// <summary>
        /// 設定/取得密鑰
        /// </summary>
        public byte[] Key
        {
            set { key = value; }
            get { return key; }
        }

        /// <summary>
        /// 設定/取得密鑰
        /// </summary>
        public string KeyStr
        {
            set { Key = Encoding.UTF8.GetBytes(value); }
            get { return Encoding.UTF8.GetString(Key); }
        }

        /// <summary>
        /// 設定/取得初始向量
        /// </summary>
        public byte[] IV
        {
            set { iv = value; }
            get { return iv; }
        }

        /// <summary>
        /// 設定/取得初始向量
        /// </summary>
        public string IVStr
        {
            set { IV = Encoding.UTF8.GetBytes(value); }
            get { return Encoding.UTF8.GetString(IV); }
        }

        /// <summary>
        /// 密鑰
        /// </summary>
        private byte[] key = null;

        /// <summary>
        /// 初始向量
        /// </summary>
        private byte[] iv = null;
    }

    public partial class PListProc
    {
        /// <summary>
        /// plist訊息序列化
        /// </summary>
        /// <param name="sender">plist傳送器</param>
        /// <returns>訊息物件</returns>
        public static PListMsg Marshal(PListSender sender)
        {
            if (sender == null)
                throw new ArgumentNullException("sender");

            return sender.Marshal();
        }

        /// <summary>
        /// plist訊息序列化
        /// </summary>
        /// <param name="messageID">訊息編號</param>
        /// <param name="message">訊息物件</param>
        /// <returns>訊息物件</returns>
        public static PListMsg Marshal(MessageID messageID, IMessage message)
        {
            if (message == null)
                throw new ArgumentNullException("message");

            return new PListMsg { MessageID = messageID, Message = Any.Pack(message) };
        }

        /// <summary>
        /// plist訊息反序列化
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

            if (input is not PListMsg data)
                throw new InvalidMessageException("unmarshal");

            messageID = data.MessageID;
            message = data.Message.Unpack<T>();
        }
    }
}

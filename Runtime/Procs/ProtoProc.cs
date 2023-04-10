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
    /// proto處理器, 封包結構使用ProtoMsg, 可以選擇是否啟用base64編碼或是des-cbc加密
    /// 訊息定義: support/proto/mizugo/protomsg.proto
    /// 封包編碼: protobuf編碼成位元陣列, (可選)des-cbc加密, (可選)base64編碼
    /// 封包解碼: (可選)base64解碼, (可選)des-cbc解密, protobuf解碼成訊息結構
    /// </summary>
    public partial class ProtoProc : Procmgr
    {
        /// <summary>
        /// proto使用的填充模式
        /// </summary>
        private const PaddingMode padding = PaddingMode.PKCS7;

        public override byte[] Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not ProtoMsg message)
                throw new InvalidMessageException("encode");

            var output = message.ToByteArray();

            if (desCBC)
                output = DesCBC.Encrypt(padding, desKey, desIV, output);

            if (base64)
                output = Base64.Encode(output);

            return output;
        }

        public override object Decode(byte[] input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (base64)
                input = Base64.Decode(input);

            if (desCBC)
                input = DesCBC.Decrypt(padding, desKey, desIV, input);

            var message = ProtoMsg.Parser.ParseFrom(input);
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
                throw new UnprocessException(message.MessageID);

            process(message);
        }

        public ProtoProc SetBase64(bool enable)
        {
            base64 = enable;
            return this;
        }

        public ProtoProc SetDesCBC(bool enable, string key, string iv)
        {
            desCBC = enable;
            desKey = Encoding.UTF8.GetBytes(key);
            desIV = Encoding.UTF8.GetBytes(iv);
            return this;
        }

        /// <summary>
        /// 是否啟用base64
        /// </summary>
        private bool base64 = false;

        /// <summary>
        /// 是否啟用des-cbc加密
        /// </summary>
        private bool desCBC = false;

        /// <summary>
        /// des密鑰
        /// </summary>
        private byte[] desKey = null;

        /// <summary>
        /// des初始向量
        /// </summary>
        private byte[] desIV = null;
    }

    public partial class ProtoProc
    {
        /// <summary>
        /// proto訊息序列化
        /// </summary>
        /// <param name="messageID">訊息編號</param>
        /// <param name="message">訊息物件</param>
        /// <returns>訊息物件</returns>
        public static ProtoMsg Marshal(MessageID messageID, IMessage message)
        {
            if (message == null)
                throw new ArgumentNullException("message");

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

            if (input is not ProtoMsg data)
                throw new InvalidMessageException("unmarshal");

            messageID = data.MessageID;
            message = data.Message.Unpack<T>();
        }
    }
}

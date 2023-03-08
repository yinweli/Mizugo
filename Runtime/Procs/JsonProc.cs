using Newtonsoft.Json;
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
    /// json處理器, 封包結構使用JsonMsg, 可以選擇是否啟用base64編碼或是des-cbc加密
    /// 訊息定義: support/proto/mizugo/msg-go/msgs-json/jsonmsg.go
    /// 訊息定義: support/proto/mizugo/msg-cs/msgs-json/Jsonmsg.cs
    /// 封包編碼: json編碼成位元陣列, (可選)des-cbc加密, (可選)base64編碼
    /// 封包解碼: (可選)base64解碼, (可選)des-cbc解密, json解碼成訊息結構
    /// </summary>
    public partial class JsonProc : Procmgr
    {
        /// <summary>
        /// json使用的填充模式
        /// </summary>
        private const PaddingMode padding = PaddingMode.PKCS7;

        public override byte[] Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not JsonMsg message)
                throw new InvalidMessageException("encode");

            var json = JsonConvert.SerializeObject(message);
            var output = Encoding.UTF8.GetBytes(json);

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

            var json = Encoding.UTF8.GetString(input);
            var message = JsonConvert.DeserializeObject<JsonMsg>(json);
            return message;
        }

        public override void Process(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not JsonMsg message)
                throw new InvalidMessageException("process");

            var process = Get(message.MessageID);

            if (process == null)
                throw new UnprocessException(message.MessageID);

            process(message);
        }

        public JsonProc SetBase64(bool enable)
        {
            base64 = enable;
            return this;
        }

        public JsonProc SetDesCBC(bool enable, string key, string iv)
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

    public partial class JsonProc
    {
        /// <summary>
        /// json訊息序列化
        /// </summary>
        /// <param name="messageID">訊息編號</param>
        /// <param name="message">訊息物件</param>
        /// <returns>訊息物件</returns>
        public static JsonMsg Marshal(MessageID messageID, object message)
        {
            if (message == null)
                throw new ArgumentNullException("message");

            var json = JsonConvert.SerializeObject(message);
            var jsonBytes = Encoding.UTF8.GetBytes(json);
            return new JsonMsg { MessageID = messageID, Message = jsonBytes };
        }

        /// <summary>
        /// json訊息反序列化
        /// </summary>
        /// <typeparam name="T">訊息類型</typeparam>
        /// <param name="input">輸入物件</param>
        /// <param name="messageID">訊息編號</param>
        /// <param name="message">訊息物件</param>
        public static void Unmarshal<T>(object input, out MessageID messageID, out T message)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not JsonMsg data)
                throw new InvalidMessageException("unmarshal");

            var json = Encoding.UTF8.GetString(data.Message);

            messageID = data.MessageID;
            message = JsonConvert.DeserializeObject<T>(json);
        }
    }
}

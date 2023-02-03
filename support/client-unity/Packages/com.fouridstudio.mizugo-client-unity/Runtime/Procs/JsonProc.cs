using Newtonsoft.Json;
using System;
using System.Text;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    /// <summary>
    /// json處理器, 封包結構使用JsonMsg
    /// 沒有使用加密技術, 所以安全性很低, 僅用於傳送簡單訊息或是傳送密鑰使用
    /// 訊息內容: support/proto/mizugo/msg-cs/msgs-json/Jsonmsg.cs
    /// 封包編碼: json編碼成位元陣列, 再通過base64編碼
    /// 封包解碼: base64解碼, 再通過json解碼成訊息結構
    /// </summary>
    public partial class JsonProc : Procmgr
    {
        public override byte[] Encode(object input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            if (input is not JsonMsg message)
                throw new InvalidMessageException("encode");

            var json = JsonConvert.SerializeObject(message);
            var jsonBytes = Encoding.UTF8.GetBytes(json);
            var encode = Base64.Encode(jsonBytes);
            return encode;
        }

        public override object Decode(byte[] input)
        {
            if (input == null)
                throw new ArgumentNullException("input");

            var decode = Base64.Decode(input);
            var json = Encoding.UTF8.GetString(decode);
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
